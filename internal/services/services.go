package services

import (
	"time"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/repository"
)

type Services struct {
	Auth    AuthServiceInterface
	Reg     RegistrationServiceInterface
	JWT     JWTServiceInterface
	Order   OrderServiceInterface
	Accrual AccrualServiceInterface
}

func New(s repository.Storage, c *config.Config) *Services {
	services := &Services{}
	services.JWT = NewJWTService(c.TokenSecret, time.Duration(c.TokenExpMinutes)*time.Minute)
	services.Auth = NewAuthService(s.UserRepository(), services.JWT)
	services.Reg = NewRegistrationService(s.UserRepository())
	services.Order = NewOrderService(s.OrderRepository())
	services.Accrual = NewAccrualService(s.OrderRepository(), c.AccrualAddress)
	return services
}
