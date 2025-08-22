package middlewares

import (
	"net/http"
	"slices"

	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/services"
)

func Authorization(allowedRoles ...models.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			userCtx, err := services.GetUserFromContext(ctx)
			if err != nil {
				log.Error("Failed to get user context from context", logger.F.Error(err))
				http.Error(w, "Failed to authorize user", http.StatusInternalServerError)
				return
			}

			if !hasAnyAllowedRole(userCtx.Roles, allowedRoles) {
				log.Warn("Attempt to access resource without required role",
					logger.F.String("userID", userCtx.ID),
					logger.F.Any("userRoles", userCtx.Roles),
					logger.F.Any("requiredRoles", allowedRoles))
				http.Error(w, "User does not have access to the resource", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hasAnyAllowedRole(userRoles models.Roles, required models.Roles) bool {
	for _, requiredRole := range required {
		if slices.Contains(userRoles, requiredRole) {
			return true
		}
	}
	return false
}
