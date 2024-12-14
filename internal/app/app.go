package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	config *config.Config
	db     *sql.DB
	router *goexpress.Router
	logger *logging.Logger
}

func New(config *config.Config, conn *sql.DB, router *goexpress.Router, logger *logging.Logger) *App {
	return &App{
		config: config,
		db:     conn,
		router: router,
		logger: logger,
	}
}

func (a *App) registerGlobalMiddlewares() {
	a.router.Use(goexpress.LogRequest)
	a.router.Use(goexpress.StripTrailingSlashes)
	a.router.Use(goexpress.RecoverFromPanic)
}

func (a *App) AddBaseHandler() *BaseHandler {
	repo := NewRepo(a.db)
	service := NewService(repo)
	htmlTemplate := html.NewTemplate(&a.config.HTML, a.logger)
	return NewHandler(a.router, service, a.config, htmlTemplate, a.logger)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.router, a.AddBaseHandler())
}
