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

type WithdrawRepository interface {
	Withdrawer
}

type BalanceRepository interface {
	BalanceProvider
}

type Services struct {
	Auth       AuthServiceInterface
	Reg        RegistrationServiceInterface
	JWT        JWTServiceInterface
	Order      OrderServiceInterface
	Accrual    AccrualServiceInterface
	Withdrawal WithdrawalServiceInterface
	Balance    BalanceServiceInterface
}

func New(
	users UserRepository, orders OrderRepository,
	withdrawals WithdrawRepository, balance BalanceRepository,
	c *config.Config) *Services {

	services := &Services{}
	services.JWT = NewJWTService(c.TokenSecret, time.Duration(c.TokenExpMinutes)*time.Minute)
	services.Balance = NewBalanceService(balance)
	services.Auth = NewAuthService(users, services.JWT)
	services.Reg = NewRegistrationService(users)
	services.Order = NewOrderService(orders)
	services.Accrual = NewAccrualService(orders, c.AccrualAddress)
	services.Withdrawal = NewWithdrawalService(withdrawals, services.Balance, services.Order)

	return services
}
