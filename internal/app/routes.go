package app

import (
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/auth"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/session"
	"github.com/ferdiebergado/goexpress"
)

func registerBaseRoutes(router *goexpress.Router, handler *BaseHandler, sessMgr session.Manager) {
	router.Get("/dashboard", handler.HandleDashboard, goexpress.Middleware(auth.RequireUserMiddleware(sessMgr)))
	router.Get("/dbstats", handler.HandleDBStats)
	router.Get("/api/health", handler.HandleHealthCheck)
}
