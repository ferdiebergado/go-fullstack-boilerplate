package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/auth"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
	"github.com/ferdiebergado/goexpress"
)

type App struct {
	cfg            *config.Config
	db             *sql.DB
	router         *goexpress.Router
	htmlTemplate   *html.Template
	sessionManager session.Manager
}

func New(cfg *config.Config, conn *sql.DB, router *goexpress.Router, htmlTmpl *html.Template, sessMgr session.Manager) *App {
	return &App{
		cfg:            cfg,
		db:             conn,
		router:         router,
		htmlTemplate:   htmlTmpl,
		sessionManager: sessMgr,
	}
}

func (a *App) registerGlobalMiddlewares() {
	a.router.Use(goexpress.LogRequest)
	a.router.Use(goexpress.StripTrailingSlashes)
	a.router.Use(goexpress.Middleware(auth.SessionMiddleware(a.cfg.Session, a.sessionManager)))
	a.router.Use(goexpress.RecoverFromPanic)
}

func (a *App) AddBaseHandler() *BaseHandler {
	repo := NewRepo(a.db, &a.cfg.DB)
	service := NewService(repo, a.cfg)
	return NewHandler(a.router, service, a.cfg, a.htmlTemplate)
}

func (a *App) AddAuthHandler() *auth.Handler {
	repo := auth.NewAuthRepo(&a.cfg.DB, a.db)
	service := auth.NewAuthService(a.cfg, repo)
	return auth.NewHandler(a.cfg, a.router, service, a.htmlTemplate, a.sessionManager)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.router, a.AddBaseHandler(), a.sessionManager)
	auth.RegisterAuthRoutes(a.router, a.AddAuthHandler(), a.sessionManager)
}
