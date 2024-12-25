package app

import (
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/html"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/goexpress"
	gkitResponse "github.com/ferdiebergado/gopherkit/http/response"
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

type Health struct {
	CPU *CPUHealth `json:"cpu,omitempty"`
	RAM *RAMHealth `json:"ram,omitempty"`
}

type ComponentHealth struct {
	DB  *DBHealth `json:"db,omitempty"`
	App *Health   `json:"app,omitempty"`
}

type HealthResponse struct {
	Status     string          `json:"status"`
	Components ComponentHealth `json:"components"`
}

func (h *BaseHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	dbHealth, err := h.service.DBStats(r.Context())

	if err != nil {
		healthResponse := &HealthResponse{
			Status:     "unhealthy",
			Components: ComponentHealth{DB: dbHealth},
		}

		healthErr := &errtypes.HTTPError{
			AppError: &errtypes.AppError{Description: err.Error(), Err: err, Severity: errtypes.High},
			Code:     http.StatusServiceUnavailable,
		}

		response.RenderError(w, healthErr, healthResponse)
	}

	cpuHealth := h.service.CPUStats()
	ramHealth := h.service.MemStats()
	healthResponse := &HealthResponse{
		Status: "healthy",
		Components: ComponentHealth{DB: dbHealth, App: &Health{
			CPU: cpuHealth,
			RAM: ramHealth}},
	}

	if err := gkitResponse.JSON(w, http.StatusOK, healthResponse); err != nil {
		jsonErr := &errtypes.HTTPError{
			AppError: &errtypes.AppError{Description: "failed to encode to json", Err: err, Severity: errtypes.High},
			Code:     http.StatusInternalServerError,
		}

		response.RenderError[any](w, jsonErr, nil)
	}
}

func (h *BaseHandler) HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.htmlTemplate.Render(w, nil, "404.html")
}
