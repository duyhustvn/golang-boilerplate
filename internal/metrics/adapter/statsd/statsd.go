package statsd

import (
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics/adapter"
	"time"

	sd "github.com/smira/go-statsd"
)

// Collector for statsd
type Collector struct {
	client *sd.Client
}

// NewCollector
func NewCollector(addr, prefix string, log logger.Logger, flushPeriod time.Duration) *Collector {
	client := sd.NewClient(addr, sd.MetricPrefix(prefix), sd.FlushInterval(flushPeriod))
	log.Infof("new_statsd_client address: %v, flushPeriod: %v", addr, flushPeriod)
	return &Collector{
		client: client,
	}
}

// RegisterCounter create a new counter for statsd
func (c *Collector) RegisterCounter(opts adapter.CollectorOptions) adapter.Counter {
	return &Counter{client: c.client, name: opts.Name}
}

func (c *Collector) RegisterGauge(opts adapter.CollectorOptions) adapter.Gauge {
	return &Gauge{client: c.client, name: opts.Name}
}

func (c *Collector) RegisterTimer(opts adapter.CollectorOptions) adapter.Timer {
	return &Timer{client: c.client, name: opts.Name}
}

// Shutdown close the current statsd connection
func (c *Collector) Shutdown() {
	_ = c.client.Close()
}
