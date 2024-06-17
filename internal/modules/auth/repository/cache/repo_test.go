package authcacherepo_test

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	authcacherepo "boilerplate/internal/modules/auth/repository/cache"
	redisclient "boilerplate/pkg/redis"
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewRedisRepo(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	rdb := redisclient.NewUniversalRedisClient(config.Redis{
		Addrs: endpoint,
	})

	cfg := config.Config{
		Env: config.Env{
			Environment: "test",
		},
	}
	log, err := logger.GetLogger(&cfg)
	if err != nil {
		t.Fatalf("failed to init logger %+v", err)
	}

	repo := authcacherepo.NewRedisRepo(rdb, log)

	_ = repo

	t.Run("Save New User", func(t *testing.T) {
		tests := []struct {
			username string
			password string
		}{
			{"admin", "changeme"},
			{"admin", "update1"},
			{"user1", "password1"},
		}

		for _, test := range tests {
			err = repo.SaveNewUser(ctx, test.username, test.password)
			if err != nil {
				t.Fatalf("should get no error but got %+v", err)
			}

			result, err := rdb.HGet(ctx, authcacherepo.PrefixKey, test.username).Result()
			if err != nil {
				t.Fatalf("should get no error but got %+v", err)
			}
			if result != test.password {
				t.Fatalf("should get %s got %s", test.password, result)
			}
		}
	})
}
