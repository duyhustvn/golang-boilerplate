package server

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	"boilerplate/internal/middleware"
	authcacherepo "boilerplate/internal/modules/auth/repository/cache"
	authsqlrepo "boilerplate/internal/modules/auth/repository/sql"
	authsvc "boilerplate/internal/modules/auth/service"
	authkafka "boilerplate/internal/modules/auth/transport/kafka"
	authrest "boilerplate/internal/modules/auth/transport/rest"
	healthchecksvc "boilerplate/internal/modules/healthcheck/service"
	healthcheckrest "boilerplate/internal/modules/healthcheck/transport/rest"
	kafkaclient "boilerplate/pkg/kafka"
	postgres "boilerplate/pkg/postgresql"
	redisclient "boilerplate/pkg/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strings"
	"sync"

	_ "boilerplate/docs"

	swaggerMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

type Server struct {
	router           *mux.Router
	kafkaConn        *kafka.Conn
	Cfg              *config.Config
	log              logger.Logger
	metricsCollector metrics.IMetricCollector
}

// GetApp returns main app
func GetApp() *Server {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %+v\n", err)
	}

	if err := loadVars(cfg); err != nil {
		log.Fatalf("Error loading var: %+v\n", err)
	}

	log, err := logger.GetLogger(cfg)
	if err != nil {
		log.Fatalf("Error initialize custom logger: %s\n", err)
	}

	kafkaConn, err := kafkaclient.NewKafkaConnection(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Cannot connect to kafka %+v", err)
	}

	return &Server{
		router:    mux.NewRouter(),
		Cfg:       cfg,
		log:       log,
		kafkaConn: kafkaConn,
	}
}

func loadVars(c *config.Config) error {
	c.Env.GetKeys()
	c.Redis.GetRedisEnv()
	c.Logger.GetLoggerEnv()
	c.Server.GetHTTPSEnv()
	c.Kafka.GetKafkaEnv()

	if err := c.Postgres.GetPostgresEnv(); err != nil {
		return err
	}

	if _, err := c.Monitoring.GetMonitoringEnv(); err != nil {
		return err
	}

	return nil
}

// @title Swagger
// @version 1.0
// @description This is a list of API

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func (s *Server) Run() error {
	defer s.kafkaConn.Close()

	postgresCfg := s.Cfg.Postgres
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresCfg.Username, url.QueryEscape(postgresCfg.Password), postgresCfg.Host, postgresCfg.Port, postgresCfg.DBName)
	sqlxDB, err := postgres.NewSqlx(dataSourceName, s.log)
	if err != nil {
		return err
	}
	defer sqlxDB.Close()

	rdb := redisclient.NewUniversalRedisClient(s.Cfg.Redis)
	authCacheRepo := authcacherepo.NewRedisRepo(rdb, s.log)
	authSqlRepo := authsqlrepo.NewSqlRepo(sqlxDB, s.log)

	authSvc := authsvc.NewAuthSvc(authCacheRepo, authSqlRepo, s.log)
	s.metricsCollector = metrics.NewCollector(s.Cfg, s.log)

	authReaderMessageProcess := authkafka.NewAuthMessageProcessor(s.log, s.Cfg, authSvc, s.metricsCollector)
	brokers := strings.Split(s.Cfg.Kafka.Brokers, ",")
	cg := kafkaclient.NewConsumerGroup(brokers, "SASL_PLAIN", s.Cfg.Kafka.GroupID, s.log, *s.Cfg)

	go cg.ConsumeTopic(context.Background(), []string{s.Cfg.Kafka.SignupUserTopic}, s.Cfg.Kafka.PoolSize, authReaderMessageProcess.ProcessMessage)

	// profiling go programs
	// https://go.dev/blog/pprof
	s.router.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)

	// API route
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	// Health check
	healthCheckSvc, _ := healthchecksvc.NewHealthCheckSvc(s.log)
	healthCheckHandlers := healthcheckrest.NewHealthCheckHandlers(apiRouter, s.log, s.Cfg, healthCheckSvc, s.metricsCollector)
	healthCheckHandlers.RegisterRouter()

	authHandlers := authrest.NewAuthHandlers(apiRouter, s.log, s.Cfg, authSvc, s.metricsCollector)
	authHandlers.RegisterRouter()

	// Swagger: http://localhost:8001/docs
	s.router.Handle("/swagger.yaml", http.FileServer(http.Dir("./docs")))
	opts := swaggerMiddleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := swaggerMiddleware.SwaggerUI(opts, nil)
	s.router.Handle("/docs", sh)

	mw := middleware.NewMiddleware(s.log)
	s.router.Use(mw.LoggingMiddleware)

	runHTTP := func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Printf("Listening on port: %s ...", s.Cfg.Server.Port)

		if err := http.ListenAndServe(fmt.Sprintf(":%s", s.Cfg.Server.Port), s.router); err != nil {
			log.Fatal("ListenAndServe error: ", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go runHTTP(wg)
	wg.Wait()

	return nil
}
