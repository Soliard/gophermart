package services

import (
	"context"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
)

type Withdrawer interface {
	Create(ctx context.Context, w *models.Withdrawal) error
	WithdrawalExists(ctx context.Context, userID string, orderNumber string) (bool, error)
	GetWithdrawals(ctx context.Context, userID string) ([]*models.Withdrawal, error)
}

type withdrawalService struct {
	repo    Withdrawer
	balance BalanceServiceInterface
	orders  OrderServiceInterface
}

func NewWithdrawalService(repo Withdrawer, balance BalanceServiceInterface, orders OrderServiceInterface) *withdrawalService {
	return &withdrawalService{
		repo:    repo,
		balance: balance,
		orders:  orders,
	}
}

func (s *withdrawalService) ProcessWithdraw(ctx context.Context, userID, orderNumber string, sum float64) error {
	isValid := s.orders.ValidateOrderNumber(ctx, orderNumber)
	if !isValid {
		return errs.OrderIsNotValid
	}

	exists, err := s.repo.WithdrawalExists(ctx, userID, orderNumber)
	if err != nil {
		return err
	}
	if exists {
		return errs.WithdrawalAlreadyProcessed
	}

	balance, err := s.balance.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Current < sum {
		return errs.BalanceInsufficient
	}

	withdrawal := models.NewWithdrawal(userID, orderNumber, sum)
	err = s.repo.Create(ctx, withdrawal)
	if err != nil {
		return err
	}

	return nil
}

func (s *withdrawalService) GetWithdrawals(ctx context.Context, userID string) ([]*models.Withdrawal, error) {
	return s.repo.GetWithdrawals(ctx, userID)
}
