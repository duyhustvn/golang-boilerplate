package metrics

import (
	"boilerplate/internal/config"
	"fmt"

	"github.com/smira/go-statsd"
)

type Metrics struct {
	Client *statsd.Client
	Cfg    *config.Config
	// stats key
	RegisterKafkaMessage string
	SuccessKafkaMessage  string
	ErrorKafkaMessage    string
	RegisterUserDelay    string
}

func NewMetrics(client *statsd.Client, cfg *config.Config) *Metrics {
	return &Metrics{
		Client:               client,
		Cfg:                  cfg,
		RegisterKafkaMessage: fmt.Sprintf("%s_register_kafka_message", cfg.Env.ServiceName),
		SuccessKafkaMessage:  fmt.Sprintf("%s_success_processed_kafka_message_total", cfg.Env.ServiceName),
		ErrorKafkaMessage:    fmt.Sprintf("%s_error_processed_kafka_message_total", cfg.Env.ServiceName),
		RegisterUserDelay:    fmt.Sprintf("%s_register_delay", cfg.Env.ServiceName),
	}
}

func (m *Metrics) StringTag(name, value string) statsd.Tag {
	return statsd.StringTag(name, value)
}

func (m *Metrics) IntTag(name string, value int) statsd.Tag {
	return statsd.IntTag(name, value)
}
