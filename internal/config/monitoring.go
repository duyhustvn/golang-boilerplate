package config

import (
	"strconv"
)

type StatsdConfig struct {
	// Addr host:port
	Addr string
	// Prefix for statsd metrics
	Prefix string
	// FlushPeriod in ms
	FlushPeriod int
}

// MonitoringConfig for profiler
type MonitoringConfig struct {
	Statsd *StatsdConfig
}

func (m *MonitoringConfig) GetMonitoringEnv() (*MonitoringConfig, error) {
	m.Statsd = &StatsdConfig{}
	m.Statsd.Addr = GetEnv("STATSD_ADDR")
	m.Statsd.Prefix = GetEnv("STATSD_PREFIX")

	flushPeriod, err := strconv.Atoi(GetEnv("STATSD_FLUSH_PERIOD"))
	if err != nil {
		return nil, err
	}

	m.Statsd.FlushPeriod = flushPeriod
	return m, nil
}
