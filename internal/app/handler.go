package app

import (
	"context"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
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

func (h *BaseHandler) HandleDashboard(w http.ResponseWriter, _ *http.Request) {
	h.htmlTemplate.Render(w, "dashboard.html", nil)
}

type HealthResponse struct {
	Status     string          `json:"status"`
	Components ComponentHealth `json:"components"`
}

func (h *BaseHandler) performHealthCheck(_ context.Context) *HealthResponse {
	cpuHealth := h.service.CPUStats()
	ramHealth := h.service.MemStats()
	return &HealthResponse{
		Status: "healthy",
		Components: ComponentHealth{App: &Health{
			CPU: cpuHealth,
			RAM: ramHealth}},
	}
}

func (h *BaseHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	health := h.performHealthCheck(r.Context())

	response.RenderJSON(w, http.StatusOK, health)
}
