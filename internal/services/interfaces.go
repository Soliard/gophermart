package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/models"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *dto.LoginRequest) (token string, err error)
}

type OrderServiceInterface interface {
	UploadOrder(ctx context.Context, userID, orderNumber string) (*models.Order, error)
	ValidateOrderNumber(ctx context.Context, orderNumber string) bool
}

type JWTServiceInterface interface {
	GenerateToken(u *models.User) (string, error)
	GetClaims(token string) (*UserContext, error)
}
