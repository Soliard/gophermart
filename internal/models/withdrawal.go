package models

import (
	"time"

	"github.com/google/uuid"
)

type Withdrawal struct {
	ID          string    `json:"-" db:"id"`
	UserID      string    `json:"-" db:"user_id"`
	OrderNumber string    `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

func NewWithdrawal(userID, orderNumber string, sum float64) *Withdrawal {
	return &Withdrawal{
		ID:          uuid.New().String(),
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now().UTC(),
	}
}
