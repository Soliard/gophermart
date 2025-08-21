package models

import (
	"time"
)

var (
	StatusNew        OrderStatus = "NEW"
	StatusRegistered OrderStatus = "REGISTERED"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type OrderStatus string

type Order struct {
	Number     string      `json:"number" db:"number"`
	UserID     string      `json:"-" db:"user_id"`
	Status     OrderStatus `json:"status" db:"status"`
	Accrual    *float64    `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

func NewOrder(number, userID string) *Order {
	return &Order{
		Number:     number,
		UserID:     userID,
		Status:     StatusNew,
		Accrual:    nil,
		UploadedAt: time.Now().UTC(),
	}
}
