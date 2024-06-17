package authsqlrepo

import (
	"boilerplate/internal/logger"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type sqlRepo struct {
	client *sqlx.DB

	log logger.Logger
}

func NewSqlRepo(client *sqlx.DB, log logger.Logger) *sqlRepo {
	return &sqlRepo{client: client, log: log}
}

func (repo *sqlRepo) SaveNewUser(ctx context.Context, username string, hashPassword string) (int64 /*inserted*/, error) {
	query := `INSERT INTO users(username, password) VALUES (:username, :password)`
	result, err := repo.client.NamedExecContext(ctx, query, map[string]interface{}{"username": username, "password": hashPassword})
	if err != nil {
		return 0, fmt.Errorf("[SaveNewUser] %+v", err)
	}

	inserted, _ := result.RowsAffected()
	return inserted, nil
}
