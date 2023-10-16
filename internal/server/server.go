package server

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	authkafka "boilerplate/internal/modules/auth/delivery/kafka"
	authrest "boilerplate/internal/modules/auth/delivery/rest"
	authrepo "boilerplate/internal/modules/auth/repository"
	authsvc "boilerplate/internal/modules/auth/service"
	kafkaclient "boilerplate/pkg/kafka"
	redisclient "boilerplate/pkg/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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
	env, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %+v\n", err)
	}

	if err := loadVars(env); err != nil {
		log.Fatalf("Error loading var: %+v\n", err)
	}

	log, err := logger.GetLogger(env)
	if err != nil {
		log.Fatalf("Error initialize custom logger: %s\n", err)
	}

	kafkaConn, err := kafkaclient.NewKafkaConnection(context.Background(), env)
	if err != nil {
		log.Fatalf("Cannot connect to kafka %+v", err)
	}
	defer kafkaConn.Close()

	return &Server{
		router:    mux.NewRouter(),
		Cfg:       env,
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
	if _, err := c.Monitoring.GetMonitoringEnv(); err != nil {
		return err
	}

	return nil
}

// Run the https server
func (s *Server) Run() {
	rdb := redisclient.NewUniversalRedisClient(s.Cfg.Redis)
	authCacheRepo := authrepo.NewRedisRepo(rdb, s.log)

	authSvc := authsvc.NewAuthSvc(authCacheRepo, s.log)
	s.metricsCollector = metrics.NewCollector(s.Cfg, s.log)

	authReaderMessageProcess := authkafka.NewAuthMessageProcessor(s.log, s.Cfg, authSvc, s.metricsCollector)
	brokers := strings.Split(s.Cfg.Kafka.Brokers, ",")
	cg := kafkaclient.NewConsumerGroup(brokers, s.Cfg.Kafka.GroupID, s.log)

	go cg.ConsumeTopic(context.Background(), []string{s.Cfg.Kafka.SignupUserTopic}, s.Cfg.Kafka.PoolSize, authReaderMessageProcess.ProcessMessage)

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	authHandlers := authrest.NewAuthHandlers(apiRouter, s.log, s.Cfg, authSvc, s.metricsCollector)
	authHandlers.RegisterRouter()

	// Healthz
	apiRouter.HandleFunc("/healthz", func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is running")
	}).Methods(http.MethodGet)

	runHTTP := func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println((fmt.Sprintf("Listening on port: %s ...", s.Cfg.Server.Port)))

		if err := http.ListenAndServe(fmt.Sprintf(":%s", s.Cfg.Server.Port), s.router); err != nil {
			log.Fatal("ListenAndServe error: ", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go runHTTP(wg)
	wg.Wait()
}
