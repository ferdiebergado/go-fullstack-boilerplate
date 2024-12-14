package app

import (
	"github.com/ferdiebergado/goexpress"
)

func registerRoutes(router *goexpress.Router, handler BaseHandler) {
	router.Get("/", handler.HandleNotFound)
	router.Get("/dbstats", handler.HandleDBStats)
}
