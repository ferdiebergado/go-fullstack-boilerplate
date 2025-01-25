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
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
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
	database, err := db.Connect(cfg.DB)
	if err != nil {
		return err
	}

	// Register OS Signal Listener
	dbSignalCtx, dbCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer dbCancel()

	// WaitGroup to wait for all shutdown tasks to complete
	var wg sync.WaitGroup

	// Goroutine to handle database connection closure on signal
	wg.Add(1)
	go db.WaitDisconnect(dbSignalCtx, &wg, database)

	// Create the application
	sessionManager := session.NewDatabaseSession(cfg.Session, database)
	htmlTemplate := html.NewTemplate(&cfg.HTML)
	router := goexpress.New()
	application := New(cfg, database, router, htmlTemplate, sessionManager)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(&cfg.Server, router)

	// Goroutine to handle server shutdown on signal
	wg.Add(1)
	go httpServer.WaitForShutdown(&wg)

	// Start the server
	wg.Add(1)
	go httpServer.Start()

	// Block until both database and server shutdown are complete
	wg.Wait() // Wait for all shutdown tasks (including database)
	slog.Info("All shutdown tasks completed. Exiting.")

	return nil
}
