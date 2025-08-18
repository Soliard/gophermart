package models

import (
	"time"
)

var (
	StatusRegistered OrderStatus = "REGISTERED" //заказ зарегистрирован, но вознаграждение не рассчитано
	StatusProcessing OrderStatus = "PROCESSING" //расчёт начисления в процессе
	StatusInvalid    OrderStatus = "INVALID"    //заказ не принят к расчёту, и вознаграждение не будет начислено
	StatusProcessed  OrderStatus = "PROCESSED"  //расчёт начисления окончен
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
		Status:     StatusRegistered,
		Accrual:    nil,
		UploadedAt: time.Now().UTC(),
	}
}
