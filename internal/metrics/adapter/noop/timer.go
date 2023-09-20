package noop

import "boilerplate/internal/metrics/adapter"

type Timer struct {
}

func (s *Timer) Observe(n int64, labels adapter.Labels) {
}
