package postgr

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}

func (r *BalanceRepository) GetUserBalance(ctx context.Context, userID string) (*models.Balance, error) {
	balance := models.Balance{}
	balance.UserID = userID
	query := `
		SELECT 
			o.total_accrual - w.total_withdrawn as current,
			w.total_withdrawn as withdrawn
		FROM 
			(SELECT COALESCE(SUM(accrual), 0) as total_accrual FROM orders WHERE user_id = $1) o,
			(SELECT COALESCE(SUM(sum), 0) as total_withdrawn FROM withdrawals WHERE user_id = $1) w;
	`
	err := r.db.GetContext(ctx, &balance, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return &balance, nil
}
