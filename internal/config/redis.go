package config

import "strconv"

// Redis struct
type Redis struct {
	Addrs    string
	Password string
	Port     string
	Channel  string
	DB       int
	PoolSize int
}

func (r *Redis) GetRedisEnv() *Redis {
	r.Addrs = GetEnv("REDIS_ADDRS")
	r.Password = GetEnv("REDIS_PASSWORD")
	db, err := strconv.Atoi(GetEnv("REDIS_DB"))
	if err != nil {
		r.DB = 0
	} else {
		r.DB = db
	}

	r.Channel = GetEnv("REDIS_CHANNEL")
	ps, err := strconv.Atoi(GetEnv("REDIS_POOL_SIZE"))
	if err != nil {
		r.PoolSize = 1
	} else {
		r.PoolSize = ps
	}
	return r
}
