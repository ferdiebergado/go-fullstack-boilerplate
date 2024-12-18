package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/goexpress"
)

type Server struct {
	*http.Server
	config config.HTTPServerConfig
}

var ErrServerStart = errors.New("failed to start the server")

func New(cfg config.HTTPServerConfig, router *goexpress.Router) *Server {
	srv := &http.Server{
		Addr:              cfg.Addr + ":" + cfg.Port,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return &Server{
		Server: srv,
		config: cfg,
	}
}

// Starts the HTTP Server
func (s *Server) Start(ctx context.Context) error {
	slog.Info("Starting http server...")

	go func() {
		<-ctx.Done()
		slog.Info("Received shutdown signal.")
		s.Shutdown()
	}()

	slog.Info("HTTP Server listening", "addr", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w %v", ErrServerStart, err)
	}

	return nil
}

// Shuts down the HTTP Server
func (s *Server) Shutdown() {
	slog.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.Server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed shutting down the server", slog.String("error", err.Error()))
	}

	slog.Info("HTTP Server shut down successfully.")
}
