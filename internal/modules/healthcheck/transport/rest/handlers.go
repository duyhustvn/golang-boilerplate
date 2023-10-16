package healthcheckrest

import (
	"boilerplate/internal/common"
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	healthchecksvc "boilerplate/internal/modules/healthcheck/service"
	"net/http"

	"github.com/gorilla/mux"
)

type healthcheckHandlers struct {
	router           *mux.Router
	log              logger.Logger
	cfg              *config.Config
	metricsCollector metrics.IMetricCollector
	healthCheckSvc   *healthchecksvc.HealthCheckSvc
}

func NewHealthCheckHandlers(router *mux.Router, log logger.Logger, cfg *config.Config, healthcheckSvc *healthchecksvc.HealthCheckSvc, metricCollector metrics.IMetricCollector) *healthcheckHandlers {
	return &healthcheckHandlers{router: router, log: log, cfg: cfg, healthCheckSvc: healthcheckSvc, metricsCollector: metricCollector}
}

func (handler *healthcheckHandlers) HealthCheckHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
