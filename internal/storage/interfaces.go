package storage

import (
	"context"

	"github.com/Soliard/gophermart/internal/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UserExists(ctx context.Context, login string) (bool, error)
}
