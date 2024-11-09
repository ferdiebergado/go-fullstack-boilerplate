package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

func Connect(ctx context.Context, dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = db.PingContext(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping the database: %v\n", err)
		os.Exit(1)
	}

	return db
}
