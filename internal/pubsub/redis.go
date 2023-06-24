package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) Publish(ctx context.Context, channel string, payload string) error {
	if err := r.Client.Publish(ctx, channel, payload).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Subscribe(ctx context.Context, channel string) interface{} {
	// there is no error because go-redis automatically reconencts on error
	subscribe := r.Client.Subscribe(ctx, channel)
	fmt.Println("Subscribe to channel: ", channel)

	return subscribe
}

func NewRedisPubSub(rdb *redis.Client) *Redis {
	return &Redis{Client: rdb}
}
