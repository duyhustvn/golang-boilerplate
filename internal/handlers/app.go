package handlers

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/pubsub"
	"boilerplate/internal/redis"
	"context"
	"fmt"
	"log"
)

type App struct {
	Pubsub pubsub.Pubsub
	Cfg    *config.Config
	log    logger.Logger
}

// GetApp returns main app
func GetApp() *App {
	env, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s\n", err)
	}

	loadVars(env)

	log, err := logger.GetLogger(env)
	if err != nil {
		log.Fatalf("Error initialize custom logger: %s\n", err)
	}

	rdb, err := redis.Open(env)
	if err != nil {
		log.Fatalf("Error open redis connection: %s\n", err)
	}

	if err := rdb.Publish(context.Background(), env.Redis.Channel, "hello"); err != nil {
		fmt.Printf("publish to channel %s failed \n", env.Redis.Channel)
		fmt.Println(err)
	} else {
		fmt.Println("Publish to channel successfully")
	}

	redisPubSub := pubsub.NewRedisPubSub(rdb)

	return &App{
		Pubsub: redisPubSub,
		Cfg:    env,
		log:    log,
	}
}

func loadVars(c *config.Config) {
	c.Redis.GetRedisEnv()
	c.Logger.GetLoggerEnv()
	c.Server.GetHTTPSEnv()
}
