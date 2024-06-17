package middleware

import (
	"boilerplate/internal/logger"
	"net/http"
	"time"
)

type middleWare struct {
	log logger.Logger
}

func NewMiddleware(log logger.Logger) *middleWare {
	return &middleWare{log: log}
}

func (m *middleWare) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)

		m.log.Infof("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
