package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-e-commerce-app-go/backend/internal/config"
	"ai-e-commerce-app-go/backend/internal/database"
	"ai-e-commerce-app-go/backend/internal/routes"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	router := routes.New(routes.Dependencies{
		Config: cfg,
		DB:     db,
		Logger: logger,
	})

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting api server", "port", cfg.HTTPPort, "env", cfg.AppEnv)
		serverErr <- router.Run(":" + cfg.HTTPPort)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			logger.Error("server stopped with error", "error", err)
			os.Exit(1)
		}
	}

	time.Sleep(250 * time.Millisecond)
	logger.Info("api server stopped")
}
