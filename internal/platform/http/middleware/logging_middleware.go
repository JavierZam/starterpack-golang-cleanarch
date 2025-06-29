package middleware

import (
	"net/http"
	"time"

	"starterpack-golang-cleanarch/internal/utils/log"
)

// LoggingMiddleware logs incoming requests and their duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Infof(r.Context(), "Request: Method=%s Path=%s Duration=%s", r.Method, r.URL.Path, duration)
	})
}
