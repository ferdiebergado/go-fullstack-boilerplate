package app

import (
	"net/http"

	router "github.com/ferdiebergado/go-express"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/web/html"
)

func AddRoutes(router *router.Router) *router.Router {
	// Add routes here, see https://github.com/ferdiebergado/go-express for the documentation.
	router.Get("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		html.Render(w, nil, "pages/home.html")
	})

	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		html.Render(w, nil, "pages/404.html")
	})

	return router
}
