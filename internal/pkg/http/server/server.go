package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	router "github.com/ferdiebergado/go-express"
)

const (
	serverShutdownTimeout = 10
	serverReadTimeout     = 10
	serverWriteTimeout    = 10
	serverIdleTimeout     = 60
)

// Starts the HTTP server
func Start(ctx context.Context, router *router.Router, port string) error {
	// Configure the server
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
