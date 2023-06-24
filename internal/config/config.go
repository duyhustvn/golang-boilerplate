package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    Env
	Logger Logger
	Redis  Redis
	Server HTTPS
}

// GetEnv return environment variable of name "name"
func GetEnv(name string) string {
	return os.Getenv(name)
}

// Load loads the config from file
func Load() (*Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &c, nil
}
