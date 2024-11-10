package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-express/middleware"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	serverShutdownTimeout = 10
	serverReadTimeout     = 10
	serverWriteTimeout    = 10
	serverIdleTimeout     = 60
)

func run(ctx context.Context, _ []string, getenv func(string) string, _ io.Reader, _, stderr io.Writer) error {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Connect to the database
	db := Connect(ctx, getenv("DATABASE_URL"))
	defer db.Close()

	// Handle OS interrupt signals
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Parse server timeouts
	readTimeout := parseTimeout(getenv("SERVER_READ_TIMEOUT"), serverReadTimeout)
	writeTimeout := parseTimeout(getenv("SERVER_WRITE_TIMEOUT"), serverWriteTimeout)
	idleTimeout := parseTimeout(getenv("SERVER_IDLE_TIMEOUT"), serverIdleTimeout)

	// Initialize router and add middlewares
	router := setupRouter()

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:         "localhost:" + getenv("PORT"),
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeout * int(time.Second)),
		WriteTimeout: time.Duration(writeTimeout * int(time.Second)),
		IdleTimeout:  time.Duration(idleTimeout * int(time.Second)),
	}

	// Start the server in a goroutine
	go startServer(httpServer, stderr)

	// Wait for shutdown
	waitForShutdown(ctx, httpServer, stderr)

	return nil
}

// Helper to parse timeout values with a default
func parseTimeout(value string, defaultValue int) int {
	timeout, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return timeout
}

// Setup router with middlewares
func setupRouter() *router.Router {
	r := router.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.PanicRecovery)
	app.AddRoutes(r)
	return r
}

// Start the HTTP server and handle errors
func startServer(httpServer *http.Server, stderr io.Writer) {
	fmt.Printf("HTTP Server listening on %s...\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
	}
}

// Wait for context cancellation and gracefully shut down the server
func waitForShutdown(ctx context.Context, httpServer *http.Server, stderr io.Writer) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Connect(ctx context.Context, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal("Unable to connect to database", err)
	}

	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("Unable to ping the database", err)
	}

	return db
}
