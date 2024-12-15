package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/env"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Loads an environment file when in development
func loadEnvFile() error {
	const envVar = "APP_ENV"
	const envFile = ".env"
	const dev = "development"

	if environment := env.Get(envVar, dev); environment == dev {
		if err := env.Load(envFile); err != nil {
			return err
		}
	}

	return nil
}

// Run the application
func run(ctx context.Context) error {
	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load .env file when in development mode
	if err := loadEnvFile(); err != nil {
		return err
	}

	// Load config
	cfg := config.Load()

	// Connect to the database.
	database := db.New(cfg.DB)
	conn, err := database.Connect(signalCtx)
	if err != nil {
		return err
	}

	defer database.Disconnect()

	// Create the router
	router := goexpress.New()

	// Create the application
	application := app.New(cfg, conn, router)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(cfg.Server, router)
	if err := httpServer.Start(signalCtx); err != nil {
		return err
	}

	return nil
}

func main() {
	logging.SetLogger()

	slog.Info("Running application...")

	if err := run(context.Background()); err != nil {
		slog.Error("Fatal error occurred.", "reason", err, "severity", "FATAL")
		os.Exit(1)
	}

	slog.Info("Done.")
}
