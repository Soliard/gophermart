package handlers

import "github.com/Soliard/gophermart/internal/services"

type Handlers struct {
	User       *userHandler
	Order      *orderHandler
	Balance    *balanceHandler
	Withdrawal *withdrawalHandler
}

func New(services *services.Services) *Handlers {
	return &Handlers{
		User:       NewUserHandler(services.Reg, services.Auth),
		Order:      NewOrderHandler(services.Order),
		Balance:    NewBalanceHandler(services.Balance),
		Withdrawal: NewWithdrawalHandler(services.Withdrawal),
	}
}
