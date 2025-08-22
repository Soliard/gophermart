package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/models"
)

type BalanceProvider interface {
	GetUserBalance(ctx context.Context, userID string) (*models.Balance, error)
}

type balanceService struct {
	repo BalanceProvider
}

func NewBalanceService(repo BalanceProvider) *balanceService {
	return &balanceService{
		repo: repo,
	}
}

func (s *balanceService) GetBalance(ctx context.Context, userID string) (*models.Balance, error) {
	return s.repo.GetUserBalance(ctx, userID)
}
