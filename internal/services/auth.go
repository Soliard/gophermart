package services

import (
	"context"
	"errors"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	user repository.UserRepositoryInterface
	jwt  JWTServiceInterface
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtService JWTServiceInterface) *authService {
	return &authService{
		user: userRepo,
		jwt:  jwtService,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	now := time.Now().UTC()
	u, err := s.user.GetByLogin(ctx, req.Login)
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

	s.user.UpdateLoginTime(ctx, u.ID, now)

	return tokenString, nil
}
