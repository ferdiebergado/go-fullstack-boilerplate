package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/server"
	genv "github.com/ferdiebergado/gopherkit/env"
	glog "github.com/ferdiebergado/gopherkit/log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Runs the application
func run(ctx context.Context, dsn string, port string) error {
	// Register OS Signal Listener
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to the database.
	conn, err := db.Connect(ctx, dsn)

	if err != nil {
		return err
	}

	// Close the database connection after running the application
	defer func() {
		log.Println("Closing database connection...")

		if err = conn.Close(); err != nil {
			log.Printf("conn close: %v", err)
		}

		log.Println("Done.")
	}()

	// Initialize the application.
	application := app.NewApp(conn, router.NewRouter(), glog.CreateLogger())

	// Start the server
	if err = server.Start(signalCtx, application.Router, port); err != nil {
		return err
	}

	// Wait for a Signal from the OS
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
