package authsvc

import (
	"boilerplate/internal/logger"
	authrepo "boilerplate/internal/modules/auth/repository"
	"context"
)

type AuthSvc struct {
	cacheRepo authrepo.ICacheRepository
	log       logger.Logger
}

func NewAuthSvc(cacheRepo authrepo.ICacheRepository, log logger.Logger) *AuthSvc {
	return &AuthSvc{cacheRepo: cacheRepo, log: log}
}

func (svc *AuthSvc) Register(ctx context.Context, username string, password string) error {
	if err := svc.cacheRepo.SaveNewUser(ctx, username, password); err != nil {
		svc.log.Errorf("[Register] error %+v", err)
		return err
	}

	return nil
}
