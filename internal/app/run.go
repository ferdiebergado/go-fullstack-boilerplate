package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
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
	go func() {
		defer wg.Done() // Signal completion of database shutdown
		<-dbSignalCtx.Done()

		slog.Info("Closing database connection...")

		if err := conn.Close(); err != nil {
			slog.Error("failed to close the database connection", "error", err)
			return
		}

		slog.Info("Database closed successfully.")
	}()

	// Create the router
	router := goexpress.New()

	// Create the application
	application := New(cfg, conn, router)
	application.SetupRouter()

	// Start the httpServer
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Addr, cfg.Server.Port),
		Handler: router,
	}

	// Channel for server shutdown completion
	idleConnsClosed := make(chan struct{})

	// Goroutine to handle server shutdown on signal
	go func() {
		defer wg.Done() // Signal completion of server shutdown
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		slog.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("HTTP server Shutdown error", "error", err)
		}
		close(idleConnsClosed)
	}()

	// Start the server
	go func() {
		slog.Info("HTTP Server listening", "addr", cfg.Server.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server ListenAndServe error", "error", err) // Handle unexpected errors
		}
		slog.Info("Server has stopped listening") // Log after server shutdown
	}()

	// Block until both database and server shutdown are complete
	<-idleConnsClosed // Wait for server to shut down
	wg.Wait()         // Wait for all shutdown tasks (including database)
	slog.Info("All shutdown tasks completed. Exiting.")

	return nil
}
