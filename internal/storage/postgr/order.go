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
	query := `
		INSERT INTO orders (number, user_id, status, accrual, uploaded_at)
		VALUES ($1, $2, $3, $4, $5)
	`
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
