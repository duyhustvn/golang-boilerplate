package statsd

import (
	"boilerplate/internal/metrics/adapter"

	sd "github.com/smira/go-statsd"
)

type Timer struct {
	client *sd.Client
	name   string
}

func (t *Timer) Observe(n int64, labels adapter.Labels) {
	tags := getTags(labels)
	t.client.Timing(t.name, n, tags...)
}
