package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	"github.com/ferdiebergado/goexpress"
	"github.com/ferdiebergado/gopherkit/env"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Run the application
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

	// Create the application
	application := app.New(conn, goexpress.New(), cfg)
	application.SetupRouter()

	// Start the server
	if err = server.Start(signalCtx, application.Router, cfg.Server); err != nil {
		return err
	}

	// Wait for a Signal from the OS
	<-signalCtx.Done()

	return nil
}

func main() {
	const envFile = ".env"
	const dev = "development"

	if environment := env.Get("ENV", dev); environment == dev {
		if err := env.Load(envFile); err != nil {
			log.Fatalf("failed to load %s file: %v", envFile, err)
		}
	}

	if err := run(context.Background(), config.Load()); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
