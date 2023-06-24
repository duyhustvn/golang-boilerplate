package config

// Redis struct
type Redis struct {
	Host    string
	Port    string
	Channel string
}

func (r *Redis) GetRedisEnv() *Redis {
	r.Host = GetEnv("REDIS_HOST")
	r.Port = GetEnv("REDIS_PORT")
	r.Channel = GetEnv("REDIS_CHANNEL")
	return r
}
