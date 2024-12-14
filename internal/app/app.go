package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	DB     *sql.DB
	Router *goexpress.Router
	Config *config.Config
}

func New(conn *sql.DB, router *goexpress.Router, config *config.Config) *App {
	return &App{
		DB:     conn,
		Router: router,
		Config: config,
	}
}

func (a *App) registerGlobalMiddlewares() {
	a.Router.Use(goexpress.LogRequest)
	a.Router.Use(goexpress.StripTrailingSlashes)
	a.Router.Use(goexpress.RecoverFromPanic)
}

func (a *App) AddBaseHandler() *BaseHandler {
	repo := NewRepo(a.DB)
	service := NewService(repo)
	htmlTemplate := html.NewTemplate(&a.Config.HTML)
	return NewHandler(a.Router, service, a.Config, htmlTemplate)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerRoutes(a.Router, *a.AddBaseHandler())
}
