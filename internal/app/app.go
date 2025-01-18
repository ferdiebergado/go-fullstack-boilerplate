package app

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/app/user"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/middleware"
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

func New(cfg *config.Config, conn *sql.DB, router *goexpress.Router) *App {
	return &App{
		cfg:            cfg,
		db:             conn,
		router:         router,
		htmlTemplate:   html.NewTemplate(&cfg.HTML),
		sessionManager: session.NewMemorySessionStore(cfg.Server.SessionDuration),
	}
}

func (a *App) registerGlobalMiddlewares() {
	a.router.Use(goexpress.LogRequest)
	a.router.Use(goexpress.StripTrailingSlashes)
	a.router.Use(goexpress.Middleware(middleware.SessionMiddleware(a.cfg.Server, a.sessionManager)))
	a.router.Use(goexpress.RecoverFromPanic)
}

func (a *App) AddBaseHandler() *BaseHandler {
	repo := NewRepo(a.db, &a.cfg.DB)
	service := NewService(repo, a.cfg)
	return NewHandler(a.router, service, a.cfg, a.htmlTemplate)
}

func (a *App) AddAuthHandler() *user.Handler {
	repo := user.NewAuthRepo(&a.cfg.DB, a.db)
	service := user.NewAuthService(a.cfg, repo)
	return user.NewHandler(a.cfg, a.router, service, a.htmlTemplate, a.sessionManager)
}

func (a *App) SetupRouter() {
	a.registerGlobalMiddlewares()
	registerBaseRoutes(a.router, a.AddBaseHandler())
	user.RegisterAuthRoutes(a.router, a.AddAuthHandler())
}
