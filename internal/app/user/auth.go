package user

import (
	"database/sql"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type Auth struct {
	config *config.Config
	db     *sql.DB
	router *goexpress.Router
}

func NewAuthHandler(config *config.Config, conn *sql.DB, router *goexpress.Router) *Auth {
	return &Auth{
		config: config,
		db:     conn,
		router: router,
	}
}

func (a *Auth) AddAuthHandler() *Handler {
	repo := NewAuthRepo(&a.config.DB, a.db)
	service := NewAuthService(a.config, repo)
	htmlTemplate := html.NewTemplate(&a.config.HTML)
	return NewHandler(a.config, a.router, service, htmlTemplate)
}
