package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

// Starts the HTTP server
func Start(ctx context.Context, router *router.Router, cfg config.HTTPServerConfig) error {
	// Configure the server
	srv := &http.Server{
		Addr:              cfg.Addr + ":" + cfg.Port,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
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
		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
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
