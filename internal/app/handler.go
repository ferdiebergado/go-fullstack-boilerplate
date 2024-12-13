package app

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type BaseHandler struct {
	Router  *goexpress.Router
	Service Service
}

func NewHandler(router *goexpress.Router, service Service) *BaseHandler {
	return &BaseHandler{
		Router:  router,
		Service: service,
	}
}

func (h *BaseHandler) registerGlobalMiddlewares() {
	h.Router.Use(goexpress.LogRequest)
	h.Router.Use(goexpress.StripTrailingSlashes)
	h.Router.Use(goexpress.RecoverFromPanic)
}

func (h *BaseHandler) RegisterRoutes() {
	h.registerGlobalMiddlewares()
	h.Router.Get("/", h.HandleNotFound)
	h.Router.Get("/dbstats", h.HandleDBStats)
}

func (h *BaseHandler) HandleDBStats(w http.ResponseWriter, _ *http.Request) {
	stats := h.Service.Stats()
	html.Render(w, stats, "pages/dbstats.html")
}

func (h *BaseHandler) HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	html.Render(w, nil, "pages/404.html")
}
