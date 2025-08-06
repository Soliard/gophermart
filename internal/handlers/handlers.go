package handlers

import "github.com/Soliard/gophermart/internal/services"

type Handlers struct {
	User *userHandler
}

func New(services *services.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(services.User),
	}
}
