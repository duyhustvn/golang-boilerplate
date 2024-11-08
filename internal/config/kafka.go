package config

import "strconv"

type Kafka struct {
	Brokers         string
	GroupID         string
	ClientUser      string
	ClientPassword  string
	PoolSize        int
	SignupUserTopic string
}

func (k *Kafka) GetKafkaEnv() *Kafka {
	k.Brokers = GetEnv("KAFKA_BROKERS")
	k.ClientUser = GetEnv("KAFKA_CLIENT_USER")
	k.ClientPassword = GetEnv("KAFKA_CLIENT_PASSWORD")
	k.GroupID = GetEnv("KAFKA_GROUP_ID")
	ps, err := strconv.Atoi(GetEnv("KAFKA_POOL_SIZE"))
	if err != nil {
		k.PoolSize = 30
	} else {
		k.PoolSize = ps
	}

	k.SignupUserTopic = GetEnv("KAFKA_SIGNUP_USER_TOPIC")
	return k
}
