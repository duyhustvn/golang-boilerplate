package kafkaclient

import (
	"context"
	"sync"

	"boilerplate/internal/config"
	"boilerplate/internal/logger"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// MessageProcessor processor methods must implement kafka.Worker func method interface
type MessageProcessor interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
}

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

type ConsumerGroup interface {
	ConsumeTopic(ctx context.Context, cancel context.CancelFunc, groupID, topic string, poolSize int, worker Worker)
	GetNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type consumerGroup struct {
	brokers       []string
	groupID       string
	authMechanism string
	log           logger.Logger
	cfg           config.Config
}

// NewConsumerGroup kafka consumer group constructor
func NewConsumerGroup(
    brokers []string, 
    authMechanism string, 
    groupID string, 
    log logger.Logger, 
    cfg config.Config,
) *consumerGroup {
	return &consumerGroup{brokers: brokers, groupID: groupID, authMechanism: authMechanism, log: log, cfg: cfg}
}

// GetNewKafkaReader create new kafka reader
func (c *consumerGroup) GetNewKafkaReader(
    kafkaURL []string, 
    groupTopics []string, 
    groupID string,
) *kafka.Reader {
	c.log.Infof("Listen to topic: %+v", groupTopics)
	dialer := kafka.Dialer{
		Timeout: dialTimeout,
	}

	if c.authMechanism == "SASL_PLAIN" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: c.cfg.Kafka.ClientUser,
			Password: c.cfg.Kafka.ClientPassword,
		}
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		GroupTopics:            groupTopics,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		MaxAttempts:            maxAttempts,
		MaxWait:                maxWait,
		Dialer:                 &dialer,
	})
}

// GetNewKafkaWriter create new kafka producer
func (c *consumerGroup) GetNewKafkaWriter() *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		// Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}

	return w
}

// ConsumeTopic start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker) {
	r := c.GetNewKafkaReader(c.brokers, groupTopics, c.groupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.log.Warnf("consumerGroup.r.Close: %v", err)
		}
	}()

	c.log.Infof("Starting consumer groupID: %s, topic: %+v, pool size: %v", c.groupID, groupTopics, poolSize)

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}
