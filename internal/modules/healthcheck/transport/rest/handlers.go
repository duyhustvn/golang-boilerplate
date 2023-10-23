package healthcheckrest

import (
	"boilerplate/internal/common"
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	"boilerplate/internal/metrics/adapter"
	healthchecksvc "boilerplate/internal/modules/healthcheck/service"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Metrics struct {
	HealcheckApiRequests adapter.Counter
	SuccessHealcheckApi  adapter.Counter
	ErrorHealcheckApi    adapter.Counter
}

var (
	metricsRegistry Metrics
	metricsOne      sync.Once
)

type healthcheckHandlers struct {
	router           *mux.Router
	log              logger.Logger
	cfg              *config.Config
	metricsCollector metrics.IMetricCollector
	healthCheckSvc   *healthchecksvc.HealthCheckSvc
}

func registerMetrics(metricsCollector metrics.IMetricCollector) {
	metricsRegistry.HealcheckApiRequests = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name: "healcheck_api_requests_total",
		Help: "The total number of healcheck api requests",
	})

	metricsRegistry.SuccessHealcheckApi = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name: "success_healcheck_api_total",
		Help: "The total number of success healcheck api",
	})

	metricsRegistry.ErrorHealcheckApi = metricsCollector.RegisterCounter(adapter.CollectorOptions{
		Name: "error_healcheck_api_total",
		Help: "The total number of error healcheck api",
	})
}

func registerMetricsOnce(metricsCollector metrics.IMetricCollector) {
	metricsOne.Do(func() { registerMetrics(metricsCollector) })
}

func NewHealthCheckHandlers(router *mux.Router, log logger.Logger, cfg *config.Config, healthcheckSvc *healthchecksvc.HealthCheckSvc, metricCollector metrics.IMetricCollector) *healthcheckHandlers {
	registerMetricsOnce(metricCollector)
	return &healthcheckHandlers{router: router, log: log, cfg: cfg, healthCheckSvc: healthcheckSvc, metricsCollector: metricCollector}
}

func (handler *healthcheckHandlers) HealthCheckHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsRegistry.HealcheckApiRequests.Add(1, adapter.Labels{"record_type": "healcheck_api_requests_total"})
		err := handler.healthCheckSvc.HealthCheck()
		if err != nil {
			common.ResponseError(w, http.StatusInternalServerError, nil, err.Error())
			return
		}

		res := struct {
			Status string `json:"status"`
		}{
			Status: "running",
		}

		common.ResponseOk(w, http.StatusOK, res)
	}
}
