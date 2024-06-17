package postgres

import (
	"boilerplate/internal/logger"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewSqlx(dataSourceName string, log logger.Logger) (*sqlx.DB, error) {
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
