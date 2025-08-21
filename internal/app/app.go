package app

import (
	"context"
	"time"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/handlers"
	"github.com/Soliard/gophermart/internal/middlewares"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/services"
	"github.com/Soliard/gophermart/internal/storage/postgr"
	"github.com/Soliard/gophermart/internal/workers"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Config   *config.Config
	Handlers *handlers.Handlers
	Services *services.Services
	db       *sqlx.DB
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := postgr.NewConnection(ctx, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	repoUser := postgr.NewUserRepository(db)
	repoOrder := postgr.NewOrderRepository(db)
	repoWithdrawal := postgr.NewWithdrawalRepository(db)
	repoBalance := postgr.NewBalanceRepository(db)

	services := services.New(repoUser, repoOrder, repoWithdrawal, repoBalance, cfg)
	handlers := handlers.New(services)

	accrualUpdater := workers.NewAccrualUpdater(services.Accrual, time.Duration(time.Second*10))
	go accrualUpdater.Start(ctx)

	return &App{
		Config:   cfg,
		Handlers: handlers,
		Services: services,
		db:       db,
	}, nil
}

func (a *App) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.Logging)

	r.Post("/api/user/register", a.Handlers.User.Register)
	r.Post("/api/user/login", a.Handlers.User.Login)

	// Защещенные роуты
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Authentication(a.Services.JWT))

		// Роуты юзеров
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authorization(models.RoleUser))
			r.Post("/api/user/orders", a.Handlers.Order.UploadOrder)
			r.Get("/api/user/orders", a.Handlers.Order.GetUserOrders)
			r.Get("/api/user/balance", a.Handlers.Balance.GetBalance)
			r.Post("/api/user/balance/withdraw", a.Handlers.Withdrawal.ProcessWithdrawal)
			r.Get("/api/user/withdrawals", a.Handlers.Withdrawal.GetWithdrawals)
		})

	})

	return r
}

func (a *App) Close() error {
	return a.db.Close()
}
