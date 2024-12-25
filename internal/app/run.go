package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
	"github.com/ferdiebergado/goexpress"
)

// Run the application
func Run(ctx context.Context) error {
	// Initialize the logger
	logging.Init()
	slog.Info("Running application...")

	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
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
	application := New(cfg, conn, router)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(cfg.Server, router)
	if err := httpServer.Start(signalCtx); err != nil {
		return err
	}

	slog.Info("Done.")

	return nil
}
