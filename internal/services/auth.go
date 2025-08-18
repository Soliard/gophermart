package services

import (
	"context"
	"errors"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginner interface {
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateLoginTime(ctx context.Context, userID string, time time.Time) error
}

type authService struct {
	loginner UserLoginner
	jwt      JWTServiceInterface
}

func NewAuthService(userRepo UserLoginner, jwtService JWTServiceInterface) *authService {
	return &authService{
		loginner: userRepo,
		jwt:      jwtService,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	now := time.Now().UTC()
	u, err := s.loginner.GetByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, errs.UserNotFound) {
			return "", errs.WrongLoginOrPassword
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errs.WrongLoginOrPassword
	}

	tokenString, err := s.jwt.GenerateToken(u)
	if err != nil {
		return "", err
	}

	u.LastLoginAt = &now

	s.loginner.UpdateLoginTime(ctx, u.ID, now)

	return tokenString, nil
}
