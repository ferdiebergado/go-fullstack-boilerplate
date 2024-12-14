package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
	"github.com/ferdiebergado/goexpress"
)

type Server struct {
	*http.Server
	config config.HTTPServerConfig
	logger *logging.Logger
}

func New(cfg config.HTTPServerConfig, router *goexpress.Router, logger *logging.Logger) *Server {
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
		logger: logger,
	}
}

// Starts the HTTP Server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting http server...")

	go func() {
		<-ctx.Done()
		s.logger.Info("Received shutdown signal.")
		s.Shutdown()
	}()

	s.logger.Info("HTTP Server listening", slog.String("addr", s.Server.Addr))
	if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Shuts down the HTTP Server
func (s *Server) Shutdown() {
	s.logger.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.Server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed shutting down the server", slog.String("error", err.Error()))
	}

	s.logger.Info("Done.")
}
