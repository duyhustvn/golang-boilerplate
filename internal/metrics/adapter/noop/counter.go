package noop

import "boilerplate/internal/metrics/adapter"

// Counter for noop
type Counter struct {
}

func (c *Counter) Inc(label adapter.Labels) {
}

func (c *Counter) Add(n int64, label adapter.Labels) {
}
