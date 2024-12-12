package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	glog "github.com/ferdiebergado/gopherkit/log"
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

	// Initialize the application.
	application := app.New(conn, router.NewRouter(), glog.CreateLogger())

	// Start the server
	if err = server.Start(signalCtx, application.Router, cfg.Server); err != nil {
		return err
	}

	// Wait for a Signal from the OS
	<-signalCtx.Done()

	return nil
}

func main() {
	config := config.LoadConfig()
	ctx := context.Background()

	if err := run(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
