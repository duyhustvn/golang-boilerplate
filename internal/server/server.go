package server

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	authkafka "boilerplate/internal/modules/auth/delivery/kafka"
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
)

type Server struct {
	Cfg *config.Config
	log logger.Logger
}

// GetApp returns main app
func GetApp() *Server {
	env, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s\n", err)
	}

	loadVars(env)

	log, err := logger.GetLogger(env)
	if err != nil {
		log.Fatalf("Error initialize custom logger: %s\n", err)
	}

	return &Server{
		Cfg: env,
		log: log,
	}
}

func loadVars(c *config.Config) {
	c.Redis.GetRedisEnv()
	c.Logger.GetLoggerEnv()
	c.Server.GetHTTPSEnv()
	c.Kafka.GetKafkaEnv()
	c.Env.GetKeys()
}

// Run the https server
func (s *Server) Run() {

	rdb := redisclient.NewUniversalRedisClient(s.Cfg.Redis)
	authCacheRepo := authrepo.NewRedisRepo(rdb, s.log)
	authSvc := authsvc.NewAuthSvc(authCacheRepo, s.log)

	readerMessageProcess := authkafka.NewReaderMessageProcessor(s.log, s.Cfg, authSvc)
	brokers := strings.Split(s.Cfg.Kafka.Brokers, ",")
	cg := kafkaclient.NewConsumerGroup(brokers, s.Cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(context.Background(), []string{s.Cfg.Kafka.SignupUserTopic}, s.Cfg.Kafka.PoolSize, readerMessageProcess.ProcessMessage)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", s.Cfg.Server.Port),
	}

	runHTTP := func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println((fmt.Sprintf("Listening on port: %s ...", s.Cfg.Server.Port)))

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe error: ", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go runHTTP(wg)
	wg.Wait()
}
