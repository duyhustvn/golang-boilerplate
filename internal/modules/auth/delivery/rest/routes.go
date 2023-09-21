package authrest

import "net/http"

func (handler *authHandlers) RegisterRouter() {
	s := handler.router.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/login", handler.Login()).Methods(http.MethodPost)
}
