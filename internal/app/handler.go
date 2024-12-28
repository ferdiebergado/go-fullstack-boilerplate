package app

import (
	"context"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
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
	h.htmlTemplate.Render(w, nil, "dashboard.html")
}

func (h *BaseHandler) HandleDBStats(w http.ResponseWriter, _ *http.Request) {
	h.htmlTemplate.Render(w, nil, "dbstats.html")
}

type HealthResponse struct {
	Status     string          `json:"status"`
	Components ComponentHealth `json:"components"`
}

func (h *BaseHandler) performHealthCheck(ctx context.Context) (*HealthResponse, *errtypes.HTTPError) {
	dbHealth, err := h.service.DBStats(ctx)

	if err != nil {
		healthResponse := &HealthResponse{
			Status:     "unhealthy",
			Components: ComponentHealth{DB: dbHealth},
		}

		healthErr := &errtypes.HTTPError{
			AppError: &errtypes.AppError{Description: err.Error(), Err: err, Severity: errtypes.High},
			Code:     http.StatusServiceUnavailable,
		}

		return healthResponse, healthErr
	}

	cpuHealth := h.service.CPUStats()
	ramHealth := h.service.MemStats()
	healthResponse := &HealthResponse{
		Status: "healthy",
		Components: ComponentHealth{DB: dbHealth, App: &Health{
			CPU: cpuHealth,
			RAM: ramHealth}},
	}

	return healthResponse, nil
}

func (h *BaseHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	health, err := h.performHealthCheck(r.Context())

	if err != nil {
		response.RenderError(w, err, health)
	}

	response.RenderJSON(w, http.StatusOK, health)
}

func (h *BaseHandler) HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.htmlTemplate.Render(w, nil, "404.html")
}
