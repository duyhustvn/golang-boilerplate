package authrepo

import "context"

type ICacheRepository interface {
	SaveNewUser(ctx context.Context, username string, password string) error
}
