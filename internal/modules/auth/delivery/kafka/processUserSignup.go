package authkafka

import (
	authmodel "boilerplate/internal/modules/auth/models"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) processUserRegister(ctx context.Context, r *kafka.Reader, msg kafka.Message) {
	s.log.Infof("[processUserRegister] %s", string(msg.Value))

	var m authmodel.User
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		s.log.Errorf("[processUserRegister] unmarshal message error %+v", err)
		s.commitErrMessage(ctx, r, msg)
		return
	}

	if err := s.authSvc.Register(ctx, m.Username, m.Password); err != nil {
		s.log.Errorf("[processUserRegister] Register error %+v", err)
		s.commitErrMessage(ctx, r, msg)
		return
	}

	s.log.Info("[processUserRegister] successful")
	s.commitMessage(ctx, r, msg)
}
