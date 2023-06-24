package redis

import (
	"boilerplate/internal/config"
	"context"

	"github.com/redis/go-redis/v9"
)

func Open(c *config.Config) (client *redis.Client, err error) {
	addr := c.Redis.Host + ":" + c.Redis.Port

	rdb := redis.NewClient(&redis.Options{Addr: addr})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
