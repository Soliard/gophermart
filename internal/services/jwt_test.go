package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateToken_GetClaims(t *testing.T) {
	secret := "test-secret-key-12345"
	expires := 1 * time.Hour
	service := NewJWTService(secret, expires)

	testUser := &models.User{
		ID:    "user123",
		Login: "testuser",
		Roles: []models.Role{models.RoleUser},
	}

	t.Run("GenerateToken and GetClaims success", func(t *testing.T) {
		tokenString, err := service.GenerateToken(testUser)
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		claims, err := service.GetClaims(tokenString)
		require.NoError(t, err)
		require.Equal(t, testUser.ID, claims.ID)
		require.EqualValues(t, testUser.Roles, claims.Roles)
	})

	t.Run("GetClaims with expired token", func(t *testing.T) {
		shortExpireService := NewJWTService(secret, -1*time.Hour)

		tokenString, err := shortExpireService.GenerateToken(testUser)
		require.NoError(t, err)

		claims, err := service.GetClaims(tokenString)
		require.Error(t, err)
		require.True(t, errors.Is(err, errs.ErrTokenExpired))
		require.Nil(t, claims)
	})

	t.Run("GetClaims with invalid signature", func(t *testing.T) {
		tokenString, err := service.GenerateToken(testUser)
		require.NoError(t, err)

		wrongSecretService := NewJWTService("wrong-secret-key", expires)

		claims, err := wrongSecretService.GetClaims(tokenString)
		require.Error(t, err)
		require.True(t, errors.Is(err, errs.ErrTokenInvalid))
		require.Nil(t, claims)
	})

	t.Run("GetClaims with malformed token", func(t *testing.T) {
		claims, err := service.GetClaims("not.a.valid.token")
		require.Error(t, err)
		require.Nil(t, claims)
	})

	t.Run("GetClaims with wrong signing method", func(t *testing.T) {
		wrongMethodToken := jwt.NewWithClaims(jwt.SigningMethodNone, &claims{
			User: UserContext{ID: testUser.ID, Roles: testUser.Roles},
		})

		tokenString, err := wrongMethodToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		claims, err := service.GetClaims(tokenString)
		require.Error(t, err)
		require.True(t, errors.Is(err, errs.ErrTokenInvalid))
		require.Nil(t, claims)
	})
}

func TestJWTService_ContextUser(t *testing.T) {
	t.Run("ContextWithUser and GetUserFromContext success", func(t *testing.T) {
		userCtx := &UserContext{
			ID:    "user123",
			Roles: []models.Role{models.RoleUser, models.RoleAdmin},
		}

		ctx := context.Background()
		ctxWithUser := ContextWithUser(ctx, userCtx)

		retrievedUser, err := GetUserFromContext(ctxWithUser)
		require.NoError(t, err)
		require.Equal(t, userCtx.ID, retrievedUser.ID)
		require.Equal(t, userCtx.Roles, retrievedUser.Roles)
	})

	t.Run("GetUserFromContext with empty context", func(t *testing.T) {
		ctx := context.Background()

		user, err := GetUserFromContext(ctx)
		require.Error(t, err)
		require.True(t, errors.Is(err, errs.ErrEmptyContextUser))
		require.Nil(t, user)
	})

	t.Run("GetUserFromContext with wrong value type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxKeyUserContext, "not-a-user-context")

		user, err := GetUserFromContext(ctx)
		require.Error(t, err)
		require.Nil(t, user)
	})
}
