package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *goexpress.Router
	Logger *logging.Logger
}

func New(config *config.Config, conn *sql.DB, router *goexpress.Router, logger *logging.Logger) *App {
	return &App{
		Config: config,
		DB:     conn,
		Router: router,
		Logger: logger,
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
	htmlTemplate := html.NewTemplate(&a.Config.HTML, a.Logger)
	return NewHandler(a.Router, service, a.Config, htmlTemplate, a.Logger)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.Router, a.AddBaseHandler())
}
