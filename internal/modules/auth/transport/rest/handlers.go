package authrest

import (
	"boilerplate/internal/common"
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	authmodel "boilerplate/internal/modules/auth/models"
	authsvc "boilerplate/internal/modules/auth/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type authHandlers struct {
	router           *mux.Router
	log              logger.Logger
	cfg              *config.Config
	authSvc          *authsvc.AuthSvc
	metricsCollector metrics.IMetricCollector
}

func NewAuthHandlers(router *mux.Router, log logger.Logger, cfg *config.Config, authSvc *authsvc.AuthSvc, metricCollector metrics.IMetricCollector) *authHandlers {
	return &authHandlers{router: router, log: log, cfg: cfg, authSvc: authSvc, metricsCollector: metricCollector}
}

func (handler *authHandlers) Login() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginUser authmodel.User
		if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
			handler.log.Errorf("[LoginHandler] Failed to parse body %+v", err)
			common.ResponseError(w, http.StatusBadRequest, nil, "invalid body")
			return
		}

		response := authmodel.LoginResponse{
			AccessToken: "update logic",
		}

		common.ResponseOk(w, http.StatusCreated, response)
	}
}
