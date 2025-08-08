package memory

import (
	"context"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
)

type Orders map[string]*models.Order

type OrderRepository struct {
	Orders Orders
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		Orders: Orders{},
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	r.Orders[order.Number] = order
	return nil
}

func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	order, ok := r.Orders[number]
	if ok {
		return order, nil
	}
	return nil, errs.OrderNotFound

}
