package app

import "database/sql"

type repo struct {
	db *sql.DB
}

type Repo interface {
	Stats() sql.DBStats
}

func NewRepo(conn *sql.DB) Repo {
	return &repo{
		db: conn,
	}
}

func (r *repo) Stats() sql.DBStats {
	return r.db.Stats()
}
