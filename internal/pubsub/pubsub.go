package pubsub

import "context"

type Pubsub interface {
	Publish(ctx context.Context, channel string, payload string) error
	Subscribe(ctx context.Context, channel string) interface{}
}
