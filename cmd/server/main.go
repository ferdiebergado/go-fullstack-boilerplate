package main

import (
	"context"
	"fmt"
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

func loadEnvFile() error {
	const envVar = "ENV"
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
func run() error {
	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create the logger
	logger := logging.New()

	// Load .env file when in development mode
	if err := loadEnvFile(); err != nil {
		return err
	}

	// Load config
	cfg := config.Load()

	// Connect to the database.
	database := db.New(cfg.DB, logger)
	conn, err := database.Connect(signalCtx)
	if err != nil {
		return err
	}

	defer database.Disconnect()

	// Create the router
	router := goexpress.New()

	// Create the application
	application := app.New(cfg, conn, router, logger)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(cfg.Server, router, logger)
	if err := httpServer.Start(signalCtx); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error occurred.\n%v\n", err)
		os.Exit(1)
	}
}
