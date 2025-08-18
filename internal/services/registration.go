package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistrator interface {
	Create(ctx context.Context, user *models.User) error
	UserExists(ctx context.Context, login string) (bool, error)
}

type registrationService struct {
	registrator UserRegistrator
}

func NewRegistrationService(userRepo UserRegistrator) *registrationService {
	return &registrationService{
		registrator: userRepo,
	}
}

func (s *registrationService) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	if req.Login == "" || req.Password == "" {
		return nil, errs.EmptyLoginOrPassword
	}

	exists, err := s.registrator.UserExists(ctx, req.Login)
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

	err = s.registrator.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
