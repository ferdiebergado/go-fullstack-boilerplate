package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	config *config.Config
	db     *sql.DB
	router *goexpress.Router
}

func New(config *config.Config, conn *sql.DB, router *goexpress.Router) *App {
	return &App{
		config: config,
		db:     conn,
		router: router,
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
	htmlTemplate := html.NewTemplate(&a.config.HTML)
	return NewHandler(a.router, service, a.config, htmlTemplate)
}

func (a *App) AddAuthHandler() *user.Handler {
	repo := user.NewAuthRepo(&a.config.DB, a.db)
	service := user.NewAuthService(a.config, repo)
	htmlTemplate := html.NewTemplate(&a.config.HTML)
	return user.NewHandler(a.config, a.router, service, htmlTemplate)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.router, a.AddBaseHandler())
	user.RegisterAuthRoutes(a.router, a.AddAuthHandler())
}
