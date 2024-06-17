package authrepo

import "context"

type AuthCacheRepo interface {
	SaveNewUser(ctx context.Context, username string, password string) error
}
