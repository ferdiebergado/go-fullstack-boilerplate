package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
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

	// Load config
	cfg := config.Load()

	// Connect to the database.
	conn, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		return err
	}

	// WaitGroup to wait for all shutdown tasks to complete
	var wg sync.WaitGroup
	wg.Add(2) // Add 2 for database and server shutdown

	// Register OS Signal Listener
	dbSignalCtx, dbCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer dbCancel()

	// Goroutine to handle database connection closure on signal
	go db.WaitDisconnect(dbSignalCtx, &wg, conn)

	// Create the router
	router := goexpress.New()

	// Create the application
	application := New(cfg, conn, router)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(&cfg.Server, router)

	// Channel for server shutdown completion
	idleConnsClosed := make(chan struct{})

	// Goroutine to handle server shutdown on signal
	go httpServer.WaitForShutdown(&wg, idleConnsClosed)

	// Start the server
	go httpServer.Start()

	// Block until both database and server shutdown are complete
	<-idleConnsClosed // Wait for server to shut down
	wg.Wait()         // Wait for all shutdown tasks (including database)
	slog.Info("All shutdown tasks completed. Exiting.")

	return nil
}
