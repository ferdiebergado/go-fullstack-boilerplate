package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

func Connect(cfg config.DBConfig) (*sql.DB, error) {
	slog.Info("Connecting to the database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.SSLMode)

	database, err := sql.Open(cfg.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()

	if err = database.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	database.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	database.SetMaxIdleConns(cfg.MaxIdleConnections)
	database.SetMaxOpenConns(cfg.MaxOpenConnections)

	slog.Info("Connected to the database", "database", cfg.DB, "user", cfg.User)
	slog.Debug("Database config", slog.Group("Connection", slog.Duration("conn_max_life_time", cfg.ConnMaxLifetime), slog.Int("max_idle_connections", cfg.MaxIdleConnections), slog.Int("max_open_connections", cfg.MaxOpenConnections)), slog.Duration("ping_timeout", cfg.PingTimeout), "ssl_mode", cfg.SSLMode)

	return database, nil
}

func WaitDisconnect(ctx context.Context, wg *sync.WaitGroup, database *sql.DB) {
	defer wg.Done()
	<-ctx.Done()

	slog.Info("Closing database connection...")

	if err := database.Close(); err != nil {
		slog.Error("failed to close the database connection", "error", err)
		return
	}

	slog.Info("Database closed successfully.")
}
