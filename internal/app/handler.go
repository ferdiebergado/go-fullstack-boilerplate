package app

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type BaseHandler struct {
	router       *goexpress.Router
	service      Service
	config       *config.Config
	htmlTemplate *html.Template
}

func NewHandler(router *goexpress.Router, service Service, cfg *config.Config, htmlTemplate *html.Template) *BaseHandler {
	return &BaseHandler{
		router:       router,
		service:      service,
		config:       cfg,
		htmlTemplate: htmlTemplate,
	}
}

func (h *BaseHandler) HandleDBStats(w http.ResponseWriter, _ *http.Request) {
	stats := h.service.Stats()
	h.htmlTemplate.Render(w, stats, "pages/dbstats.html")
}

func (h *BaseHandler) HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.htmlTemplate.Render(w, nil, "pages/404.html")
}
