package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-express/middleware"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	serverShutdownTimeout = 10
)

func run(ctx context.Context, _ []string, getenv func(string) string, _ io.Reader, _, stderr io.Writer) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	db.Connect(ctx, getenv("DATABASE_URL"))

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	readTimeout, err := strconv.Atoi(getenv("SERVER_READ_TIMEOUT"))

	if err != nil {
		readTimeout = 10
	}

	writeTimeout, err := strconv.Atoi(getenv("SERVER_WRITE_TIMEOUT"))

	if err != nil {
		writeTimeout = 10
	}

	idleTimeout, err := strconv.Atoi(getenv("SERVER_IDLE_TIMEOUT"))

	if err != nil {
		idleTimeout = 60
	}

	router := router.NewRouter()
	router.Use(middleware.RequestLogger)
	router.Use(middleware.PanicRecovery)

	app.AddRoutes(router)

	httpServer := &http.Server{
		Addr:         ":" + getenv("PORT"),
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeout * int(time.Second)),
		WriteTimeout: time.Duration(writeTimeout * int(time.Second)),
		IdleTimeout:  time.Duration(idleTimeout * int(time.Second)),
	}

	go func() {
		fmt.Printf("HTTP Server listening on %s...\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, serverShutdownTimeout*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
