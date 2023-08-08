package redisclient

import (
	"boilerplate/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries      = 5
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 3 * time.Second
	minIdleConns    = 20
	poolTimeout     = 6 * time.Second
	idleTimeout     = 12 * time.Second
)

func NewUniversalRedisClient(cfg config.Redis) redis.UniversalClient {
	universalClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           []string{cfg.Addrs},
		Password:        cfg.Password, // no password set
		DB:              cfg.DB,       // use default DB
		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    minIdleConns,
		PoolTimeout:     poolTimeout,
		ConnMaxIdleTime: idleTimeout,
	})

	return universalClient
}
