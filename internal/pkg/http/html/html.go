package html

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/logging"
)

//go:embed templates/*
var templatesFS embed.FS

type Template struct {
	templateDir string
	layoutFile  string
	logger      *logging.Logger
}

func NewTemplate(cfg *config.HTMLTemplateConfig, logger *logging.Logger) *Template {
	return &Template{
		templateDir: cfg.TemplateDir,
		layoutFile:  cfg.LayoutFile,
		logger:      logger,
	}
}

func (t *Template) Render(w http.ResponseWriter, data any, templateFiles ...string) {
	templateDir := t.templateDir
	layoutFile := t.layoutFile

	layoutTemplate := filepath.Join(templateDir, layoutFile)

	targetTemplates := []string{layoutTemplate}

	for _, file := range templateFiles {
		targetTemplate := filepath.Join(templateDir, file)
		targetTemplates = append(targetTemplates, targetTemplate)
	}

	templates, err := template.New("template").Funcs(getFuncMap()).ParseFS(templatesFS, targetTemplates...)

	if err != nil {
		response.RenderError(w, response.ServerError(err), t.logger)
		return
	}

	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, layoutFile, data); err != nil {
		response.RenderError(w, response.ServerError(err), t.logger)
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		response.RenderError(w, response.ServerError(err), t.logger)
		return
	}
}
