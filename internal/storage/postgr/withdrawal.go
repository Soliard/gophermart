package postgr

import (
	"context"

	"github.com/Soliard/gophermart/internal/models"
	"github.com/jmoiron/sqlx"
)

type WithdrawalRepository struct {
	db *sqlx.DB
}

func NewWithdrawalRepository(db *sqlx.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) Create(ctx context.Context, w *models.Withdrawal) error {
	query := `
		INSERT INTO withdrawals (id, user_id, order_number, sum, processed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, w.ID, w.UserID, w.OrderNumber, w.Sum, w.ProcessedAt)
	return err
}

func (r *WithdrawalRepository) WithdrawalExists(ctx context.Context, userID string, orderNumber string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM withdrawals WHERE user_id = $1 and order_number = $2)`
	err := r.db.GetContext(ctx, &exists, query, userID, orderNumber)
	return exists, err
}

func (r *WithdrawalRepository) GetWithdrawals(ctx context.Context, userID string) ([]*models.Withdrawal, error) {
	var withdrawals []*models.Withdrawal
	query := `SELECT * FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC`
	err := r.db.SelectContext(ctx, &withdrawals, query, userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
