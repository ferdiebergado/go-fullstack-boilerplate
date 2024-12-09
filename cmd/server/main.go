package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	genv "github.com/ferdiebergado/gopherkit/env"
	glog "github.com/ferdiebergado/gopherkit/log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	serverShutdownTimeout = 10
	serverReadTimeout     = 10
	serverWriteTimeout    = 10
	serverIdleTimeout     = 60
)

func startServer(ctx context.Context, router *router.Router, port string) error {
	// Configure HTTP server
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       serverReadTimeout * time.Second,
		ReadHeaderTimeout: serverReadTimeout * time.Second,
		WriteTimeout:      serverWriteTimeout * time.Second,
		IdleTimeout:       serverIdleTimeout * time.Second,
	}

	// Wait for shutdown signal
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()

		log.Println("Shutdown signal received.")

		// Shutdown the server
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		log.Println("Done.")
	}()

	// Start the server
	log.Printf("HTTP Server listening on %s... (Press Ctrl-C to exit)", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("HTTP server ListenAndServe: %v", err)
		return err
	}

	wg.Wait()

	return nil
}

func run(ctx context.Context, dsn string, port string) error {
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to the database.
	conn, err := db.Connect(ctx, dsn)

	if err != nil {
		return err
	}

	defer func() {
		// Close the database connection
		log.Println("Closing database connection...")

		if err = conn.Close(); err != nil {
			log.Printf("conn close: %v", err)
		}

		log.Println("Done.")
	}()

	// Initialize the application.
	application := app.NewApp(conn, router.NewRouter(), glog.CreateLogger())

	// Start the server
	if err = startServer(signalCtx, application.Router, port); err != nil {
		return err
	}

	<-signalCtx.Done()

	return nil
}

func main() {
	port := genv.Must("PORT")
	dsn := genv.Must("DATABASE_URL")

	if err := run(context.Background(), dsn, port); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
