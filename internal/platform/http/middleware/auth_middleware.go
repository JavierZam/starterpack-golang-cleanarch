package middleware

import (
	"context"
	"net/http"
	"strings"

	"starterpack-golang-cleanarch/internal/utils"
	globalErrors "starterpack-golang-cleanarch/internal/utils/errors"
	"starterpack-golang-cleanarch/internal/utils/log"
)

type ContextKey string

const (
	ContextKeyUserID   ContextKey = "userID"
	ContextKeyTenantID ContextKey = "tenantID"
	ContextKeyUserRole ContextKey = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warnf(r.Context(), "Auth: Missing Authorization header for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, globalErrors.ErrUnauthorized, r)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf(r.Context(), "Auth: Invalid Authorization header format for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, globalErrors.NewBadRequest("Invalid Authorization header format", nil), r)
			return
		}

		tokenString := tokenParts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.Warnf(r.Context(), "Auth: Invalid token for path: %s, error: %v", r.URL.Path, err)
			utils.HandleHTTPError(w, globalErrors.NewBadRequest("Invalid or expired token", nil), r)
			return
		}

		userID := claims.UserID
		tenantID := claims.TenantID
		userRole := claims.Role

		if userID == "" || tenantID == "" || userRole == "" {
			log.Warnf(r.Context(), "Auth: Claims missing (UserID/TenantID/Role) in token for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, globalErrors.ErrUnauthorized, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, ContextKeyTenantID, tenantID)
		ctx = context.WithValue(ctx, ContextKeyUserRole, userRole)

		log.Debugf(ctx, "Auth: Authenticated user %s (Role: %s) for tenant %s accessing path: %s", userID, userRole, tenantID, r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
