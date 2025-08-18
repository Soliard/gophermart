package handlers

import "github.com/Soliard/gophermart/internal/services"

type Handlers struct {
	User  *userHandler
	Order *orderHandler
}

func New(services *services.Services) *Handlers {
	return &Handlers{
		User:  NewUserHandler(services.Reg, services.Auth),
		Order: NewOrderHandler(services.Order),
	}
}
