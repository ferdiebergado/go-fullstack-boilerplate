package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type Database struct {
	conn   *sql.DB
	config *config.DBConfig
}

func New(cfg config.DBConfig) *Database {
	return &Database{
		config: &cfg,
	}
}

func (d *Database) Connect(ctx context.Context) (*sql.DB, error) {
	slog.Info("Connecting to the database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.Port, d.config.User, d.config.Password, d.config.DB, d.config.SSLMode)

	db, err := sql.Open(d.config.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("database initialization: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, d.config.PingTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	db.SetConnMaxLifetime(d.config.ConnMaxLifetime)
	db.SetMaxIdleConns(d.config.MaxIdleConnections)
	db.SetMaxOpenConns(d.config.MaxOpenConnections)

	d.conn = db

	slog.Info("Connected.")

	return db, nil
}

func (d *Database) Disconnect() {
	slog.Info("Closing database connection...")

	if err := d.conn.Close(); err != nil {
		slog.Error("failed to close the database connection", "error", err.Error())
	}

	slog.Info("Database closed successfully.")
}
