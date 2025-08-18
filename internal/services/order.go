package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/repository"
	"github.com/phedde/luhn-algorithm"
)

type orderService struct {
	order repository.OrderRepositoryInterface
}

func NewOrderService(orderRepository repository.OrderRepositoryInterface) *orderService {
	return &orderService{
		order: orderRepository,
	}
}

func (s *orderService) UploadOrder(ctx context.Context, userID, orderNumber string) (*models.Order, error) {
	order, err := s.order.GetByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, errs.OrderNotFound) {
		return nil, err
	}

	if order != nil {
		if order.UserID != userID {
			return nil, errs.OrderAlreadyUploadedByOtherUser
		} else {
			return nil, errs.OrderAlreadyUploadedByThisUser
		}
	}

	newOrder := models.NewOrder(orderNumber, userID)
	err = s.order.Create(ctx, newOrder)
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
	return s.order.GetUserOrders(ctx, userID)
}
