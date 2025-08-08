package models

import (
	"time"
)

var (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type OrderStatus string

type Order struct {
	Number     string      `json:"number"`
	UserID     string      `json:"-"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func NewOrder(number, userID string) *Order {
	return &Order{
		Number:     number,
		UserID:     userID,
		Status:     StatusNew,
		Accrual:    nil,
		UploadedAt: time.Now(),
	}
}
