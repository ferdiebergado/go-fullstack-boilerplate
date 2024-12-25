package app

import (
	"github.com/ferdiebergado/goexpress"
)

func registerBaseRoutes(router *goexpress.Router, handler *BaseHandler) {
	router.Get("/dashboard", handler.HandleDashboard)
	router.Get("/dbstats", handler.HandleDBStats)
	router.Get("/api/health", handler.HandleHealthCheck)
	router.Get("/", handler.HandleNotFound)
}
