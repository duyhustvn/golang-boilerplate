package benchrest

import (
	"boilerplate/internal/common"
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	benchmodel "boilerplate/internal/modules/bench/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type benchHandlers struct {
	router           *mux.Router
	log              logger.Logger
	cfg              *config.Config
	metricsCollector metrics.IMetricCollector
}

func NewHandlers(
    router *mux.Router, 
    log logger.Logger, 
    cfg *config.Config, 
    metricCollector metrics.IMetricCollector,
) *benchHandlers {
	return &benchHandlers{
        router: router, 
        log: log, 
        cfg: cfg, 
        metricsCollector: metricCollector,
    }
}

// @Summary Login API
// @Description Login
// @Accept json
// @Produce json
// @Param timeout query string true "Set timeout for this API"
// @Success 200 {object} common.RestResponse{data=benchmodeln.TimeoutResponse}
// @Router /api/timeout [get]
// @Tags Auth
func (handler *benchHandlers) TimeoutAPI(w http.ResponseWriter, r *http.Request) {
    timeoutStr := r.URL.Query().Get("timeout")
    if timeoutStr == "" {
        common.ResponseError(w, http.StatusBadRequest, nil, "Timeout is required") 
        return
    }

    timeout, err := strconv.Atoi(timeoutStr)
    if err != nil {
        common.ResponseError(w, http.StatusBadRequest, nil, "Timeout must be an integer") 
        return
    }

	response := benchmodel.TimeoutResponse{
        Message: "ok",
    }

    time.Sleep(time.Duration(timeout) * time.Second)

	common.ResponseOk(w, http.StatusOK, response)
}
