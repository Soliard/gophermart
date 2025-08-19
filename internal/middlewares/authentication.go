package middlewares

import (
	"net/http"
	"strings"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

const bearerPrefix = "Bearer "
const authHeader = "Authorization"

func Authentication(jwtService services.JWTServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			log := logger.FromContext(ctx)
			tokenString, err := extractTokenFromHeaders(req)
			if err != nil {
				http.Error(res, "Token not found", http.StatusUnauthorized)
				return
			}

			userInfo, err := jwtService.GetClaims(tokenString)
			if err != nil {
				log.Warn("Failed to get claims from JWT", logger.F.Error(err))
				http.Error(res, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx = services.ContextWithUser(ctx, userInfo)

			next.ServeHTTP(res, req.WithContext(ctx))

		})
	}
}

func extractTokenFromHeaders(req *http.Request) (string, error) {

	token := req.Header.Get(authHeader)
	if token == "" {
		return "", errs.ErrTokenNotFound
	}

	if strings.HasPrefix(token, bearerPrefix) {
		return token[len(bearerPrefix):], nil
	}

	return token, nil
}
