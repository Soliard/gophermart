package services

import (
	"context"
	"errors"
	"time"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

const ctxKeyUserContext ctxKey = "user"

var signingMethod jwt.SigningMethod = jwt.SigningMethodHS256

type ctxKey string

type jwtService struct {
	secret       string
	expiresAfter time.Duration
}

type UserContext struct {
	ID    string        `json:"id"`
	Roles []models.Role `json:"roles"`
}

type claims struct {
	jwt.RegisteredClaims
	User UserContext `json:"user"`
}

func NewJWTService(secret string, expires time.Duration) *jwtService {
	return &jwtService{
		secret:       secret,
		expiresAfter: expires,
	}
}

func (s *jwtService) GenerateToken(u *models.User) (string, error) {
	token := jwt.NewWithClaims(signingMethod, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresAfter)),
		},
		User: UserContext{ID: u.ID, Roles: u.Roles},
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *jwtService) GetClaims(tokenString string) (*UserContext, error) {
	claims := &claims{}
	keyfunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != signingMethod {
			return nil, errs.ErrTokenInvalid
		}
		return []byte(s.secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrTokenExpired
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errs.ErrTokenInvalid
	}

	return &claims.User, nil
}

func ContextWithUser(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, ctxKeyUserContext, u)
}

func GetUserFromContext(ctx context.Context) (*UserContext, error) {
	user, ok := ctx.Value(ctxKeyUserContext).(*UserContext)
	if !ok {
		return nil, errs.ErrEmptyContextUser
	}
	return user, nil
}
