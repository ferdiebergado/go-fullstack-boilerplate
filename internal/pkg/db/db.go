package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

func Connect(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	slog.Info("Connecting to the database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)

	slog.Info("Connected to the database", "database", cfg.DB, "user", cfg.User)
	slog.Debug("Database config", slog.Group("Connection", slog.Duration("conn_max_life_time", cfg.ConnMaxLifetime), slog.Int("max_idle_connections", cfg.MaxIdleConnections), slog.Int("max_open_connections", cfg.MaxOpenConnections)), slog.Duration("ping_timeout", cfg.PingTimeout), "ssl_mode", cfg.SSLMode)

	return db, nil
}
