package middleware

import (
	"context"
	"net/http"
	"strings"

	"starterpack-golang-cleanarch/internal/utils"
	"starterpack-golang-cleanarch/internal/utils/errors"
	"starterpack-golang-cleanarch/internal/utils/log"
)

// ContextKey defines a custom type for context keys to avoid collisions.
type ContextKey string

const (
	ContextKeyUserID   ContextKey = "userID"
	ContextKeyTenantID ContextKey = "tenantID"
)

// AuthMiddleware validates JWT and populates context with user and tenant ID.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warnf(r.Context(), "Authentication attempt without Authorization header for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, errors.ErrUnauthorized, r)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf(r.Context(), "Invalid Authorization header format for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, errors.NewBadRequest("Invalid Authorization header format", nil), r)
			return
		}

		// --- FIX: Using a VALID UUID for tenantID placeholder ---
		// In a real scenario, you'd parse tokenParts[1] (the actual JWT)
		// to extract claims like user_id, tenant_id, roles, expiration.
		// For this demo, we use hardcoded placeholders.
		userID := "8c84b42b-5f33-4f9e-b2d4-1a2b3c4d5e6f"
		tenantID := "a1b2c3d4-e5f6-4a7b-8c9d-0f1e2d3c4b5a"
		// --------------------------------------------------------

		if userID == "" || tenantID == "" {
			log.Warnf(r.Context(), "AuthMiddleware: UserID or TenantID missing after placeholder token check for path: %s", r.URL.Path)
			utils.HandleHTTPError(w, errors.ErrUnauthorized, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, ContextKeyTenantID, tenantID)
		log.Debugf(ctx, "AuthMiddleware: Authenticated user %s for tenant %s accessing path: %s", userID, tenantID, r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
