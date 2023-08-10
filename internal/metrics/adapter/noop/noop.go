package noop

import "boilerplate/internal/metrics/adapter"

// Collector for noop
type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

// RegisterCounter returns noop Counter
func (c *Collector) RegisterCounter(opts adapter.CollectorOptions) adapter.Counter {
	return &Counter{}
}

func (c *Collector) RegisterGauge(opts adapter.CollectorOptions) adapter.Gauge {
	return &Gauge{}
}

func (c *Collector) RegisterTimer(opts adapter.CollectorOptions) adapter.Timer {
	return &Timer{}
}

// Shutdown close the current statsd connection
func (c *Collector) Shutdown() {}
