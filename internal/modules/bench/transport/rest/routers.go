package benchrest

import "net/http"

func (handler *benchHandlers) RegisterRouter() {
	s := handler.router.PathPrefix("/bench").Subrouter()
	s.HandleFunc("/timeout", handler.TimeoutAPI).Methods(http.MethodGet)
}
