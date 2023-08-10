package authkafka

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	authsvc "boilerplate/internal/modules/auth/service"
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

type readerMessageProcessor struct {
	log     logger.Logger
	cfg     *config.Config
	authSvc *authsvc.AuthSvc
	metrics *metrics.Metrics
}

func NewReaderMessageProcessor(log logger.Logger, cfg *config.Config, authSvc *authsvc.AuthSvc, metrics *metrics.Metrics) *readerMessageProcessor {
	return &readerMessageProcessor{log: log, cfg: cfg, authSvc: authSvc, metrics: metrics}
}

func (s *readerMessageProcessor) ProcessMessage(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
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
			s.processUserRegister(ctx, r, msg)
		}
	}
}
