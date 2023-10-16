package kafkaclient

import (
	"boilerplate/internal/config"
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

func NewKafkaConnection(ctx context.Context, cfg *config.Config) (*kafka.Conn, error) {
	brokers := strings.Split(cfg.Kafka.Brokers, ",")
	return kafka.DialContext(ctx, "tcp", brokers[0])
}
