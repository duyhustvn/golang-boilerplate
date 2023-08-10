package authrepo

import (
	"boilerplate/internal/logger"
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

const prefixKey = "user"

type RedisRepo struct {
	prefixKey string
	client    redis.UniversalClient
	log       logger.Logger
}

func NewRedisRepo(rdb redis.UniversalClient, log logger.Logger) *RedisRepo {
	return &RedisRepo{prefixKey: prefixKey, client: rdb, log: log}
}

func (rr *RedisRepo) SaveNewUser(ctx context.Context, username string, password string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisRepo.SaveNewUser")
	defer span.Finish()

	if err := rr.client.HSet(ctx, rr.prefixKey, username, password).Err(); err != nil {
		rr.log.Errorf("[SaveNewUser] Save user error: %+v", err)
		return err
	}

	return nil
}
