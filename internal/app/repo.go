package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type repo struct {
	db  *sql.DB
	cfg *config.DBConfig
}

type Repo interface {
	Stats() sql.DBStats
	Ping(context.Context) error
}

func NewRepo(conn *sql.DB, cfg *config.DBConfig) Repo {
	return &repo{
		db:  conn,
		cfg: cfg,
	}
}

func (r *repo) Stats() sql.DBStats {
	return r.db.Stats()
}

func (r *repo) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, r.cfg.PingTimeout)
	defer cancel()

	if err := r.db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
