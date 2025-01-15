package server

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
	"github.com/ferdiebergado/goexpress"
)

type Server struct {
	server *http.Server
	cfg    *config.HTTPServerConfig
}

func New(cfg *config.HTTPServerConfig, router *goexpress.Router) *Server {
	srv := &http.Server{ // #nosec G112 -- Server timeouts delegated to reverse proxy
		Addr:    fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Handler: router,
	}

	return &Server{
		server: srv,
		cfg:    cfg,
	}
}

// Start the server
func (s *Server) Start() {
	slog.Info("HTTP Server listening", "addr", s.cfg.Addr, slog.Int("port", s.cfg.Port))

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server ListenAndServe error", "error", err)
	}

	slog.Info("Server has stopped listening")
}

// Handle server shutdown on signal
func (s *Server) WaitForShutdown(wg *sync.WaitGroup, idleConnsClosed chan struct{}) {
	defer wg.Done()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	slog.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("HTTP server Shutdown error", "error", err)
	}

	close(idleConnsClosed)
}
