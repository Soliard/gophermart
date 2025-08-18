package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/models"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, req *dto.LoginRequest) (token string, err error)
}

type RegistrationServiceInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error)
}

type OrderServiceInterface interface {
	UploadOrder(ctx context.Context, userID, orderNumber string) (*models.Order, error)
	ValidateOrderNumber(ctx context.Context, orderNumber string) bool
	GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error)
}

type JWTServiceInterface interface {
	GenerateToken(u *models.User) (string, error)
	GetClaims(token string) (*UserContext, error)
}

type AccrualServiceInterface interface {
	UpdateOrders(ctx context.Context) error
}
