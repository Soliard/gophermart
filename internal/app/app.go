package app

import (
	"context"

	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/handlers"
	"github.com/Soliard/gophermart/internal/middlewares"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/Soliard/gophermart/internal/services"
	"github.com/Soliard/gophermart/internal/storage"
	"github.com/go-chi/chi"
)

type App struct {
	Config   *config.Config
	Handlers *handlers.Handlers
	Services *services.Services
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	store, err := storage.New(ctx, cfg)
	if err != nil {
		return nil, err
	}
	services := services.New(store, cfg)
	handlers := handlers.New(services)
	return &App{
		Config:   cfg,
		Handlers: handlers,
		Services: services,
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
		})

	})

	return r
}
