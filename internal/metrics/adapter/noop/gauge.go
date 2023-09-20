package noop

import "boilerplate/internal/metrics/adapter"

type Gauge struct {
}

func (g *Gauge) Add(n int64, labels adapter.Labels) {
}

func (g *Gauge) Sub(n int64, labels adapter.Labels) {
}

func (g *Gauge) Inc(labels adapter.Labels) {
}

func (g *Gauge) Set(n int64, labels adapter.Labels) {
}
