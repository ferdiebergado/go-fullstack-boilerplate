package app

import (
	"database/sql"

	"github.com/ferdiebergado/goexpress"
)

type App struct {
	DB     *sql.DB
	Router *goexpress.Router
}

func New(conn *sql.DB, router *goexpress.Router) *App {
	return &App{
		DB:     conn,
		Router: router,
	}
}

func (a *App) AddHandlers() {
	repo := NewRepo(a.DB)
	service := NewService(repo)
	handler := NewHandler(a.Router, service)
	handler.RegisterRoutes()
}
