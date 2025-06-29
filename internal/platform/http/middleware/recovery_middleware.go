package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug" // Untuk mendapatkan stack trace

	"starterpack-golang-cleanarch/internal/utils"
	"starterpack-golang-cleanarch/internal/utils/errors"
	"starterpack-golang-cleanarch/internal/utils/log"
)

// RecoveryMiddleware recovers from panics and returns a 500 error.
// This middleware should typically be placed early in the middleware chain.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				// Log the panic with stack trace for debugging purposes
				log.Errorf(r.Context(), "Panic recovered: %v\nStack Trace:\n%s", rcv, debug.Stack())
				// Return a consistent internal server error response to the client
				utils.HandleHTTPError(w, errors.NewInternalServerError(fmt.Errorf("%v", rcv), "An unexpected server error occurred"), r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
