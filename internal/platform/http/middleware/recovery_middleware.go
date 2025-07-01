package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"starterpack-golang-cleanarch/internal/utils"
	globalErrors "starterpack-golang-cleanarch/internal/utils/errors"
	"starterpack-golang-cleanarch/internal/utils/log"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				log.Errorf(r.Context(), "Panic recovered: %v\nStack Trace:\n%s", rcv, debug.Stack())
				utils.HandleHTTPError(w, globalErrors.NewInternalServerError(fmt.Errorf("%v", rcv), "An unexpected server error occurred"), r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
