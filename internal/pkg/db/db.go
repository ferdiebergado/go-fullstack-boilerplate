package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type DeleteMode int

const (
	driver             = "pgx"
	connMaxLifetime    = 0
	maxIdleConnections = 50
	maxOpenConnections = 50
	pingTimeout        = 1

	SoftDelete DeleteMode = iota
	HardDelete
)

var ErrRowClose = errors.New("failed to close the rows result set")
var ErrRowScan = errors.New("error occurred while scanning the row into the destination variables")
var ErrRowIteration = errors.New("error encountered during row iteration, possibly due to a database or connection issue")
var ErrModelNotFound = errors.New("model not found")

func Connect(dsn string) (*sql.DB, error) {
	log.Println("Connecting to the database...")

	db, err := sql.Open(driver, dsn)

	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), pingTimeout*time.Second)
	defer cancel()

	err = db.PingContext(pingCtx)

	if err != nil {
		log.Printf("ping database: %v", err)
		return nil, err
	}

	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)

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
