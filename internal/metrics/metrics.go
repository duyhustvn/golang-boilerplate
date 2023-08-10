package metrics

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics/adapter"
	"boilerplate/internal/metrics/adapter/noop"
	"boilerplate/internal/metrics/adapter/statsd"
	"time"
)

// IMetricCollector
type IMetricCollector interface {
	RegisterCounter(adapter.CollectorOptions) adapter.Counter
	RegisterGauge(adapter.CollectorOptions) adapter.Gauge
	RegisterTimer(adapter.CollectorOptions) adapter.Timer
	Shutdown()
}

func NewCollector(cfg *config.Config, log logger.Logger) IMetricCollector {
	monitoringConfig := &cfg.Monitoring
	isStatsd := monitoringConfig != nil && monitoringConfig.Statsd != nil

	if isStatsd {
		flushDuration := 100 * time.Millisecond
		if monitoringConfig.Statsd.FlushPeriod > 0 {
			flushDuration = time.Duration(monitoringConfig.Statsd.FlushPeriod) * time.Microsecond
		}

		return statsd.NewCollector(monitoringConfig.Statsd.Addr, monitoringConfig.Statsd.Prefix, log, flushDuration)
	}

	log.Info("config_skipping_empty_metrics_provider")
	return noop.NewCollector()
}
