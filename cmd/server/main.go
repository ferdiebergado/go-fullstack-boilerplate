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
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Run the application
func run(ctx context.Context) error {
	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		slog.Error("Fatal error occurred.", "reason", err, "severity", "Fatal")
		os.Exit(1)
	}

	slog.Info("Done.")
}
