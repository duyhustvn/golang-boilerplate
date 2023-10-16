package healthcheckrest

import "net/http"

func (handler *healthcheckHandlers) RegisterRouter() {
	s := handler.router.PathPrefix("/").Subrouter()
	s.HandleFunc("/healthz", handler.HealthCheckHandler()).Methods(http.MethodGet)
}
