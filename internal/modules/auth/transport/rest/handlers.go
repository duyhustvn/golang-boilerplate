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

// @Summary Login API
// @Description Login
// @Accept json
// @Produce json
// @Param user body authmodel.User true "User object"
// @Success 200 {object} common.RestResponse{data=authmodel.LoginResponse}
// @Router /api/auth/login [post]
// @Tags Auth
func (handler *authHandlers) Login(w http.ResponseWriter, r *http.Request) {
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
