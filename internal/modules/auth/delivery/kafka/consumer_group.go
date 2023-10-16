package authkafka

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	"boilerplate/internal/metrics/adapter"
	authsvc "boilerplate/internal/modules/auth/service"
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Metrics struct {
	kafkaRegisterUserDelay    adapter.Timer
	registerUserKafkaMessages adapter.Counter
	successKafkaMessages      adapter.Counter
	errorKafkaMessages        adapter.Counter
}

var (
	metricsRegistry Metrics
	metricsOnce     sync.Once
)

type authMessageProcessor struct {
	log              logger.Logger
	cfg              *config.Config
	authSvc          *authsvc.AuthSvc
	metricsCollector metrics.IMetricCollector
}

func NewAuthMessageProcessor(log logger.Logger, cfg *config.Config, authSvc *authsvc.AuthSvc, metricsCollector metrics.IMetricCollector) *authMessageProcessor {
	registerMetricsOnce(metricsCollector)

	return &authMessageProcessor{log: log, cfg: cfg, authSvc: authSvc, metricsCollector: metricsCollector}
}

func (s *authMessageProcessor) ProcessMessage(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Errorf("[ProcessMessage] FetchMessage workerID: %d, err: %+v", workerID, err)
			continue
		}

		switch msg.Topic {
		case s.cfg.Kafka.SignupUserTopic:
			metricsRegistry.registerUserKafkaMessages.Inc(map[string]string{})
			s.processUserRegister(ctx, r, msg)
			metricsRegistry.kafkaRegisterUserDelay.Observe(int64(time.Since(msg.Time)), map[string]string{})
		}
	}
}

func registerMetricsOnce(metricsCollector metrics.IMetricCollector) {
	metricsOnce.Do(func() { registerMetrics(metricsCollector) })
}

func registerMetrics(metricsCollector metrics.IMetricCollector) {
	metricsRegistry.kafkaRegisterUserDelay = metricsCollector.RegisterTimer(adapter.CollectorOptions{
		Name:   "kafka_register_user_delay",
		Help:   "The ms spent on register user",
		Labels: []string{},
	})

	metricsRegistry.registerUserKafkaMessages = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name:   "register_user_kafka_messages_total",
		Help:   "The number of register user kafka messages",
		Labels: []string{},
	})

	metricsRegistry.successKafkaMessages = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name: "success_kafka_processed_messages",
		Help: "The total number of success kafka processed message",
	})

	metricsRegistry.errorKafkaMessages = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name: "error_kafka_processed_messages",
		Help: "The total number of error kafka processed message",
	})
}
