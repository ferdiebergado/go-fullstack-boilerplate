package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		slog.Error("Fatal error occurred.", "reason", err)
		os.Exit(1)
	}
}
