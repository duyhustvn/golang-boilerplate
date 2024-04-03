package config

import "errors"

type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func (p Postgres) Validator() error {
	if p.Host == "" || p.Port == "" || p.Password == "" || p.DBName == "" {
		return errors.New("all postgresql environment must be provided")
	}
	return nil
}

func (p *Postgres) GetPostgresEnv() error {
	p.Host = GetEnv("POSTGRES_HOST")
	p.Port = GetEnv("POSTGRES_PORT")
	p.Username = GetEnv("POSTGRES_USERNAME")
	p.Password = GetEnv("POSTGRES_PASSWORD")
	p.DBName = GetEnv("POSTGRES_DBNAME")

	if err := p.Validator(); err != nil {
		return err
	}

	return nil
}
