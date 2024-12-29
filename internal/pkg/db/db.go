package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type Database struct {
	conn   *sql.DB
	config *config.DBConfig
}

var (
	ErrDatabaseInit    = errors.New("failed to initialize database")
	ErrDatabaseConnect = errors.New("failed to connect to the database")
)

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
		return nil, fmt.Errorf("%w %v", ErrDatabaseInit, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, d.config.PingTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("%w %v", ErrDatabaseConnect, err)
	}

	db.SetConnMaxLifetime(d.config.ConnMaxLifetime)
	db.SetMaxIdleConns(d.config.MaxIdleConnections)
	db.SetMaxOpenConns(d.config.MaxOpenConnections)

	d.conn = db

	slog.Info("Connected to the database", "database", d.config.DB, "user", d.config.User)
	slog.Debug("Database Config", slog.Group("Connection", slog.Duration("conn_max_life_time", d.config.ConnMaxLifetime), slog.Int("max_idle_connections", d.config.MaxIdleConnections), slog.Int("max_open_connections", d.config.MaxOpenConnections)), slog.Duration("ping_timeout", d.config.PingTimeout), "ssl_mode", d.config.SSLMode)

	return db, nil
}

func (d *Database) Disconnect() {
	slog.Info("Closing database connection...")

	if err := d.conn.Close(); err != nil {
		slog.Error("failed to close the database connection", "error", err)
		return
	}

	slog.Info("Database closed successfully.")
}
