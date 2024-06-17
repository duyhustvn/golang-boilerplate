package authrepo

import "context"

type AuthCacheRepo interface {
	SaveNewUser(ctx context.Context, username string, password string) error
}

type AuthSqlRepo interface {
	SaveNewUser(ctx context.Context, username string, hashPassword string) (int64 /*inserted*/, error)
}
