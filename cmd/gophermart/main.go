package main

import (
	"context"
	stdlog "log"
	"net/http"

	"github.com/Soliard/gophermart/internal/app"
	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		stdlog.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Log.Info("Starting...")

	config, err := config.New()
	if err != nil {
		logger.Log.Fatal("Failed to create config", logger.F.Error(err))
	}

	logger.Log.Info("Used config", logger.F.Any("config", config))

	ctx := context.Background()
	app, err := app.New(ctx, config)
	if err != nil {
		logger.Log.Fatal("Failed to create app", logger.F.Error(err))
	}

	router := app.Router()
	server := &http.Server{
		Addr:    app.Config.ServerHost,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Server failed", logger.F.Error(err))
	}

}
