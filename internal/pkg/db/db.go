package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
)

type Database struct {
	conn   *sql.DB
	config *config.DBConfig
	logger *logging.Logger
}

func New(cfg config.DBConfig, logger *logging.Logger) *Database {
	return &Database{
		config: &cfg,
		logger: logger,
	}
}

func (d *Database) Connect(ctx context.Context) (*sql.DB, error) {
	d.logger.Info("Connecting to the database...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.Port, d.config.User, d.config.Password, d.config.DB, d.config.SSLMode)

	db, err := sql.Open(d.config.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("initialize the database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, d.config.PingTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("connect to the database: %w", err)
	}

	db.SetConnMaxLifetime(d.config.ConnMaxLifetime)
	db.SetMaxIdleConns(d.config.MaxIdleConnections)
	db.SetMaxOpenConns(d.config.MaxOpenConnections)

	d.conn = db

	d.logger.Info("Connected.")

	return db, nil
}

func (d *Database) Disconnect() {
	d.logger.Info("Closing database connection...")

	if err := d.conn.Close(); err != nil {
		d.logger.Error("failed to close the database connection", slog.String("error", err.Error()))
	}

	d.logger.Info("Done.")
}
