package app

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/web/html"
	"github.com/ferdiebergado/goexpress"
)

func MountBaseRoutes(router *goexpress.Router) *goexpress.Router {
	// Add routes here, see https://github.com/ferdiebergado/go-express for the documentation.

	// global middlewares
	router.Use(goexpress.LogRequest)
	router.Use(goexpress.StripTrailingSlashes)
	router.Use(goexpress.RecoverFromPanic)

	// home page
	router.Get("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		html.Render(w, nil, "pages/home.html")
	})

	// 404 page
	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		html.Render(w, nil, "pages/404.html")
	})

	return router
}
