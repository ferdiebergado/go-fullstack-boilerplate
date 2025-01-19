package app

import (
	"context"
	"errors"
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
	conn, err := db.Connect(ctx, cfg.DB)
	if err != nil {
		return err
	}

	// WaitGroup to wait for all shutdown tasks to complete
	var wg sync.WaitGroup
	wg.Add(3)

	// Register OS Signal Listener
	dbSignalCtx, dbCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer dbCancel()

	// Goroutine to handle database connection closure on signal
	go db.WaitDisconnect(dbSignalCtx, &wg, conn)

	// Create the application
	idleConnsClosed := make(chan struct{})
	sessionManager := session.NewInMemorySession(cfg.Server.SessionDuration)
	session, ok := sessionManager.(*session.InMemorySession)

	if !ok {
		return errors.New("sessionManager is not a session.InMemorySession")
	}

	session.StartCleanup(&wg, idleConnsClosed, cfg.Session.CleanUpInterval)
	htmlTemplate := html.NewTemplate(&cfg.HTML)
	router := goexpress.New()
	application := New(cfg, conn, router, htmlTemplate, sessionManager)
	application.SetupRouter()

	// Start the httpServer
	httpServer := server.New(&cfg.Server, router)

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
