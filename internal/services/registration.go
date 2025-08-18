package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type registrationService struct {
	repo repository.UserRepositoryInterface
}

func NewRegistrationService(userRepo repository.UserRepositoryInterface) *registrationService {
	return &registrationService{
		repo: userRepo,
	}
}

func (s *registrationService) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	if req.Login == "" || req.Password == "" {
		return nil, errs.EmptyLoginOrPassword
	}

	exists, err := s.repo.UserExists(ctx, req.Login)
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

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
