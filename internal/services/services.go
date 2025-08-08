package services

import (
	"time"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/storage"
)

type Services struct {
	User  UserServiceInterface
	JWT   JWTServiceInterface
	Order OrderServiceInterface
}

func New(s *storage.Storage, c *config.Config) *Services {
	services := &Services{}
	services.JWT = NewJWTService(c.TokenSecret, time.Duration(c.TokenExpMinutes)*time.Minute)
	services.User = NewUserService(s.UserRepository, services.JWT)
	services.Order = NewOrderService(s.OrderRepository)
	return services
}
