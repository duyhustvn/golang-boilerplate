package authsvc

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"context"
	"errors"
	"testing"
)

type cacheRepoMock struct {
	FakeSaveNewUser func(ctx context.Context, username string, password string) error
}

func NewCacheRepoMock() *cacheRepoMock {
	return &cacheRepoMock{}
}

func (m *cacheRepoMock) SaveNewUser(ctx context.Context, username string, password string) error {
	return m.FakeSaveNewUser(ctx, username, password)
}

func TestRegister(t *testing.T) {
	cfg := config.Config{
		Env: config.Env{
			Environment: "test",
		},
	}
	log, err := logger.GetLogger(&cfg)
	if err != nil {
		t.Fatalf("failed to init logger %+v", err)
	}

	t.Run("should run success", func(t *testing.T) {
		cacheRepo := cacheRepoMock{
			FakeSaveNewUser: func(ctx context.Context, username string, password string) error {
				return nil
			},
		}

		authSvc := NewAuthSvc(&cacheRepo, log)

		if err := authSvc.Register(context.Background(), "username", "password"); err != nil {
			t.Errorf("should return no error but got = %+v", err)
		}
	})

	t.Run("should return error", func(t *testing.T) {
		cacheRepo := cacheRepoMock{
			FakeSaveNewUser: func(ctx context.Context, username string, password string) error {
				return errors.New("User already exists")
			},
		}

		authSvc := NewAuthSvc(&cacheRepo, log)

		err := authSvc.Register(context.Background(), "username", "password")
		if err == nil {
			t.Errorf("should got error 'User already exists' but got no error")
		}
		if err.Error() != "User already exists" {
			t.Errorf("should got error 'User already exists' but got error = '%s'", err.Error())
		}
	})

}
