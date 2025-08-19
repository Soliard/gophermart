package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/phedde/luhn-algorithm"
)

type OrderCreator interface {
	Create(ctx context.Context, order *models.Order) error
	GetByNumber(ctx context.Context, number string) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error)
}

type orderService struct {
	creator OrderCreator
}

func NewOrderService(orderRepository OrderCreator) *orderService {
	return &orderService{
		creator: orderRepository,
	}
}

func (s *orderService) UploadOrder(ctx context.Context, userID, orderNumber string) (*models.Order, error) {
	order, err := s.creator.GetByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, errs.ErrOrderNotFound) {
		return nil, err
	}

	if order != nil {
		if order.UserID != userID {
			return nil, errs.ErrOrderAlreadyUploadedByOtherUser
		} else {
			return nil, errs.ErrOrderAlreadyUploadedByThisUser
		}
	}

	newOrder := models.NewOrder(orderNumber, userID)
	err = s.creator.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return newOrder, nil

}

func (s *orderService) ValidateOrderNumber(ctx context.Context, orderNumber string) bool {
	num, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return false
	}
	return luhn.IsValid(num)
}

func (s *orderService) GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	return s.creator.GetUserOrders(ctx, userID)
}
