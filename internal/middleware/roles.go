package middleware

import (
	"net/http"

	"platform/internal/dto"
)

// RequireRoles checks that current user role is in the allowed set.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := CurrentUser(r.Context())
			if claims == nil {
				dto.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				dto.WriteError(w, http.StatusForbidden, "Forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
