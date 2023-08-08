package config

import "strconv"

type Kafka struct {
	Brokers         string
	GroupID         string
	PoolSize        int
	SignupUserTopic string
}

func (k *Kafka) GetKafkaEnv() *Kafka {
	k.Brokers = GetEnv("KAFKA_BROKERS")
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
