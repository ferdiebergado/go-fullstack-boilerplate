package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

func Connect(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	db, err := sql.Open(cfg.Driver, cfg.DSN)

	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	if err = db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)

	log.Println("Connected.")

	return db, nil
}

func Disconnect(db *sql.DB) {
	log.Println("Closing database connection...")

	if err := db.Close(); err != nil {
		log.Printf("conn close: %v", err)
	}

	log.Println("Done.")
}
