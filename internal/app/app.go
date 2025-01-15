package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	config       *config.Config
	db           *sql.DB
	router       *goexpress.Router
	htmlTemplate *html.Template
}

func New(config *config.Config, conn *sql.DB, router *goexpress.Router) *App {
	return &App{
		config:       config,
		db:           conn,
		router:       router,
		htmlTemplate: html.NewTemplate(&config.HTML),
	}
}

func (a *App) registerGlobalMiddlewares() {
	a.router.Use(goexpress.LogRequest)
	a.router.Use(goexpress.StripTrailingSlashes)
	a.router.Use(goexpress.RecoverFromPanic)
}

func (a *App) AddBaseHandler() *BaseHandler {
	repo := NewRepo(a.db, &a.config.DB)
	service := NewService(repo, a.config)
	return NewHandler(a.router, service, a.config, a.htmlTemplate)
}

func (a *App) AddAuthHandler() *user.Handler {
	repo := user.NewAuthRepo(&a.config.DB, a.db)
	service := user.NewAuthService(a.config, repo)
	return user.NewHandler(a.config, a.router, service, a.htmlTemplate)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.router, a.AddBaseHandler())
	user.RegisterAuthRoutes(a.router, a.AddAuthHandler())
}
