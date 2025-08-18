package postgr

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func newOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `INSERT INTO orders (number, user_id, status, accrual, uploaded_at)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, order.Number, order.UserID, order.Status, order.Accrual, order.UploadedAt)
	return err
}

func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	order := &models.Order{}
	query := `SELECT * FROM orders WHERE number = $1`
	err := r.db.GetContext(ctx, order, query, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.OrderNotFound
		}
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	var orders []*models.Order
	query := `SELECT *
			  FROM orders 
			  WHERE user_id = $1 
			  ORDER BY uploaded_at DESC`
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrdersToAccrualUpdate(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order
	query := `SELECT * 
			  FROM orders 
			  WHERE status not in ($1, $2)`
	err := r.db.SelectContext(ctx, &orders, query, models.StatusInvalid, models.StatusProcessed)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatusAndAccural(
	ctx context.Context,
	numberOrder string,
	status models.OrderStatus,
	accrual *float64) error {

	query := `UPDATE orders 
			  SET status = $1, accrual = $2
			  WHERE number = $3`
	_, err := r.db.ExecContext(ctx, query, status, accrual, numberOrder)
	if err != nil {
		return err
	}
	return nil
}
