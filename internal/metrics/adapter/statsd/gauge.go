package statsd

import (
	"boilerplate/internal/metrics/adapter"

	sd "github.com/smira/go-statsd"
)

type Gauge struct {
	client *sd.Client
	name   string
}

func (g *Gauge) Add(n int64, labels adapter.Labels) {
	tags := getTags(labels)
	g.client.GaugeDelta(g.name, n, tags...)
}

func (g *Gauge) Sub(n int64, labels adapter.Labels) {
	tags := getTags(labels)
	g.client.GaugeDelta(g.name, -n, tags...)
}

func (g *Gauge) Inc(labels adapter.Labels) {
	tags := getTags(labels)
	g.client.GaugeDelta(g.name, 1, tags...)
}

func (g *Gauge) Set(n int64, labels adapter.Labels) {
	tags := getTags(labels)
	g.client.Gauge(g.name, n, tags...)
}
