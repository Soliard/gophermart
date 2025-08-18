package repository

import (
	"context"
	"time"

	"github.com/Soliard/gophermart/internal/models"
)

type Storage interface {
	UserRepository() UserRepositoryInterface
	OrderRepository() OrderRepositoryInterface
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UserExists(ctx context.Context, login string) (bool, error)
	UpdateLoginTime(ctx context.Context, userID string, time time.Time) error
}

type OrderRepositoryInterface interface {
	Create(ctx context.Context, order *models.Order) error
	GetByNumber(ctx context.Context, number string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error)
	GetOrdersToAccrualUpdate(ctx context.Context) ([]*models.Order, error)
	UpdateStatusAndAccural(
		ctx context.Context,
		numberOrder string,
		status models.OrderStatus,
		accrual *float64) error
}
