package app

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/goexpress"
)

type BaseHandler struct {
	Router       *goexpress.Router
	Service      Service
	Config       *config.Config
	HTMLTemplate *html.Template
}

func NewHandler(router *goexpress.Router, service Service, cfg *config.Config, htmlTemplate *html.Template) *BaseHandler {
	return &BaseHandler{
		Router:       router,
		Service:      service,
		Config:       cfg,
		HTMLTemplate: htmlTemplate,
	}
}

func (h *BaseHandler) HandleDBStats(w http.ResponseWriter, _ *http.Request) {
	stats := h.Service.Stats()
	h.HTMLTemplate.Render(w, stats, "pages/dbstats.html")
}

func (h *BaseHandler) HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.HTMLTemplate.Render(w, nil, "pages/404.html")
}

func (h *BaseHandler) HandleTest(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"test": "this is a test",
	}
	h.HTMLTemplate.Render(w, data, "pages/home.html")
}
