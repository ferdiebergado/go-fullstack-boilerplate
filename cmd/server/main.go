package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	"github.com/ferdiebergado/goexpress"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Runs the application
func run(ctx context.Context, cfg *config.Config) error {
	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to the database.
	conn, err := db.Connect(ctx, cfg.DB)

	if err != nil {
		return err
	}

	// Close the database connection after running the application
	defer db.Disconnect(conn)

	// Mount the router and register the routes
	router := goexpress.New()
	app.MountBaseRoutes(router)

	// Start the server
	if err = server.Start(signalCtx, router, cfg.Server); err != nil {
		return err
	}

	// Wait for a Signal from the OS
	<-signalCtx.Done()

	return nil
}

func main() {
	if err := run(context.Background(), config.Load()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
