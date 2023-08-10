package authsvc

import (
	"boilerplate/internal/logger"
	authrepo "boilerplate/internal/modules/auth/repository"
	"context"

	"github.com/opentracing/opentracing-go"
)

type AuthSvc struct {
	cacheRepo authrepo.ICacheRepository
	log       logger.Logger
}

func NewAuthSvc(cacheRepo authrepo.ICacheRepository, log logger.Logger) *AuthSvc {
	return &AuthSvc{cacheRepo: cacheRepo, log: log}
}

func (svc *AuthSvc) Register(ctx context.Context, username string, password string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthSvc.Register")
	defer span.Finish()

	if err := svc.cacheRepo.SaveNewUser(ctx, username, password); err != nil {
		svc.log.Errorf("[Register] error %+v", err)
	}

	return nil
}
