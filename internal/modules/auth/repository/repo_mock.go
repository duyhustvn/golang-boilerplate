package authrepo

import "context"

type CacheRepoMock struct {
	FakeSaveNewUser func(ctx context.Context, username string, password string) error
}

func (m *CacheRepoMock) SaveNewUser(ctx context.Context, username string, password string) error {
	return m.FakeSaveNewUser(ctx, username, password)
}
