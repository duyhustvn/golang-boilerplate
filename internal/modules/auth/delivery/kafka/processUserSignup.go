package authkafka

import (
	authmodel "boilerplate/internal/modules/auth/models"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) processUserRegister(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	start := time.Now()
	s.log.Infof("[processUserRegister] %s", string(msg.Value))
	s.metrics.Client.Incr(s.metrics.RegisterKafkaMessage, 1, s.metrics.StringTag("stat", s.metrics.RegisterKafkaMessage))
	var m authmodel.User
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		s.log.Errorf("[processUserRegister] unmarshal message error %+v", err)
		s.commitErrMessage(ctx, r, msg)
		s.metrics.Client.PrecisionTiming(s.metrics.RegisterUserDelay, time.Since(start), s.metrics.StringTag("stat", s.metrics.RegisterUserDelay))
		return
	}

	if err := s.authSvc.Register(ctx, m.Username, m.Password); err != nil {
		s.log.Errorf("[processUserRegister] Register error %+v", err)
		s.commitErrMessage(ctx, r, msg)
		s.metrics.Client.PrecisionTiming(s.metrics.RegisterUserDelay, time.Since(start), s.metrics.StringTag("stat", s.metrics.RegisterUserDelay))
		return
	}

	s.log.Info("[processUserRegister] successful")
	s.commitMessage(ctx, r, msg)
	s.metrics.Client.PrecisionTiming(s.metrics.RegisterUserDelay, time.Since(start), s.metrics.StringTag("stat", s.metrics.RegisterUserDelay))

}
