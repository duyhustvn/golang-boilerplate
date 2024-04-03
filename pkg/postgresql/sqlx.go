package postgres

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewSqlx(cfg config.Postgres, log logger.Logger) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.DBName)

	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("[NewSqlx] failed to open connection to postgresql %+v", err)
	}

	log.Info("Connected to database successfully")

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[NewSqlx] failed to ping postgresql %+v", err)
	}

	return db, nil
}
