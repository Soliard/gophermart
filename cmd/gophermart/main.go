package main

import (
	"context"
	"errors"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Soliard/gophermart/internal/app"
	"github.com/Soliard/gophermart/internal/config"
	"github.com/Soliard/gophermart/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		stdlog.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Log.Info("Starting server...")

	config, err := config.New()
	if err != nil {
		logger.Log.Fatal("Failed to create config", logger.F.Error(err))
	}

	logger.Log.Info("Used config", logger.F.Any("config", config))

	ctx, ctxStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer ctxStop()

	app, err := app.New(ctx, config)
	if err != nil {
		logger.Log.Fatal("Failed to create app", logger.F.Error(err))
	}
	defer app.Close()

	router := app.Router()
	server := &http.Server{
		Addr:              app.Config.ServerHost,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	errCh := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("shutdown signal recieved")
	case err := <-errCh:
		if err != nil {
			logger.Log.Error("server failed", logger.F.Error(err))
		}
		ctxStop()
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = server.Shutdown(ctxShutdown)
	if err != nil {
		logger.Log.Error("graceful server shutdown failed", logger.F.Error(err))
		err = server.Close()
		if err != nil {
			logger.Log.Error("forceful server shutdown failed", logger.F.Error(err))
		}
	} else {
		logger.Log.Info("server gracefully stopped")
	}
}
