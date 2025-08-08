package storage

import (
	"context"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/storage/memory"
)

type Storage struct {
	UserRepository  UserRepositoryInterface
	OrderRepository OrderRepositoryInterface
}

func New(ctx context.Context, c *config.Config) (*Storage, error) {
	return newMemoryStorage()
}

func newMemoryStorage() (*Storage, error) {
	return &Storage{
		UserRepository:  memory.NewUserRepository(),
		OrderRepository: memory.NewOrderRepository(),
	}, nil
}
