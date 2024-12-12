package app

import (
	"database/sql"
	"log/slog"

	router "github.com/ferdiebergado/go-express"
)

type App struct {
	DB     *sql.DB
	Router *router.Router
	Logger *slog.Logger
}

func New(conn *sql.DB, router *router.Router, logger *slog.Logger) *App {
	MountBaseRoutes(router)

	return &App{
		DB:     conn,
		Router: router,
		Logger: logger,
	}
}
