package statsd

import (
	"boilerplate/internal/metrics/adapter"

	sd "github.com/smira/go-statsd"
)

type Counter struct {
	client *sd.Client
	name   string
}

func (c *Counter) Add(n int64, labels adapter.Labels) {
	tags := getTags(labels)
	c.client.Incr(c.name, n, tags...)
}

func (c *Counter) Inc(labels adapter.Labels) {
	tags := getTags(labels)
	c.client.Incr(c.name, 1, tags...)
}
