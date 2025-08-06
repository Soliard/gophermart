package services

import (
	"context"
	"errors"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserRepository storage.UserRepositoryInterface
	JWTService     JWTServiceInterface
}

func NewUserService(userRepo storage.UserRepositoryInterface, jwtService JWTServiceInterface) *userService {
	return &userService{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}
}

func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	exists, err := s.UserRepository.UserExists(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.LoginAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.NewUser(req.Login, string(hashedPassword))

	err = s.UserRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	u, err := s.UserRepository.GetByLogin(ctx, req.Login)
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

	tokenString, err := s.JWTService.GenerateToken(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
