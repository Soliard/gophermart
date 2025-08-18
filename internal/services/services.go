package services

import (
	"time"

	"github.com/Soliard/gophermart/internal/config"
)

type UserRepository interface {
	UserLoginner
	UserRegistrator
}

type OrderRepository interface {
	OrderCreator
	AccrualUpdater
}

type Services struct {
	Auth    AuthServiceInterface
	Reg     RegistrationServiceInterface
	JWT     JWTServiceInterface
	Order   OrderServiceInterface
	Accrual AccrualServiceInterface
}

func New(users UserRepository, orders OrderRepository, c *config.Config) *Services {
	services := &Services{}
	services.JWT = NewJWTService(c.TokenSecret, time.Duration(c.TokenExpMinutes)*time.Minute)
	services.Auth = NewAuthService(users, services.JWT)
	services.Reg = NewRegistrationService(users)
	services.Order = NewOrderService(orders)
	services.Accrual = NewAccrualService(orders, c.AccrualAddress)
	return services
}
