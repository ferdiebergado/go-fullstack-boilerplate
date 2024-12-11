package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type DeleteMode int

const (
	SoftDelete DeleteMode = iota
	HardDelete
)

var ErrRowClose = errors.New("failed to close the rows result set")
var ErrRowScan = errors.New("error occurred while scanning the row into the destination variables")
var ErrRowIteration = errors.New("error encountered during row iteration, possibly due to a database or connection issue")
var ErrModelNotFound = errors.New("model not found")

func Connect(ctx context.Context, cfg config.DBConfig) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	db, err := sql.Open(cfg.Driver, cfg.DSN)

	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	err = db.PingContext(pingCtx)

	if err != nil {
		log.Printf("ping database: %v", err)
		return nil, err
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
