package authrest

import (
	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	"boilerplate/internal/metrics"
	authmodel "boilerplate/internal/modules/auth/models"
	authsvc "boilerplate/internal/modules/auth/service"
	"encoding/json"
	"fmt"
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
			http.Error(w, "Failed to parse body", http.StatusBadRequest)
			return
		}

		response := authmodel.LoginResponse{
			AccessToken: "update logic",
		}

		b, err := json.Marshal(response)
		if err != nil {
			handler.log.Errorf("[LoginHandler] Failed to marshal response %+v", err)
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(b))
	}
}
