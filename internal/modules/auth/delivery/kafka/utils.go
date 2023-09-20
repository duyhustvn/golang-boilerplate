package authkafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	metricsRegistry.successKafkaMessages.Inc(map[string]string{})
	if err := r.CommitMessages(ctx, msg); err != nil {
		s.log.Errorf("commitMessage error: %+v", err)
	}
}

func (s *readerMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	metricsRegistry.errorKafkaMessages.Inc(map[string]string{})
	if err := r.CommitMessages(ctx, msg); err != nil {
		s.log.Errorf("commitErrMessage error: %+v", err)
	}
}
