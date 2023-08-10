package authkafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	s.metrics.Client.Incr(s.metrics.SuccessKafkaMessage, 1, s.metrics.StringTag("stat", s.metrics.SuccessKafkaMessage))
	if err := r.CommitMessages(ctx, msg); err != nil {
		s.log.Errorf("commitMessage error: %+v", err)
	}
}

func (s *readerMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	s.metrics.Client.Incr(s.metrics.SuccessKafkaMessage, 1, s.metrics.StringTag("stat", s.metrics.ErrorKafkaMessage))
	if err := r.CommitMessages(ctx, msg); err != nil {
		s.log.Errorf("commitErrMessage error: %+v", err)
	}
}
