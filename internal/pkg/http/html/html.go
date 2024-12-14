package html

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
)

//go:embed templates/*
var templatesFS embed.FS

type Template struct {
	TemplateDir string
	LayoutFile  string
}

func NewTemplate(cfg *config.HTMLTemplateConfig) *Template {
	return &Template{
		TemplateDir: cfg.TemplateDir,
		LayoutFile:  cfg.LayoutFile,
	}
}

func (t *Template) Render(w http.ResponseWriter, data any, templateFiles ...string) {
	templateDir := t.TemplateDir
	layoutFile := t.LayoutFile

	layoutTemplate := filepath.Join(templateDir, layoutFile)

	targetTemplates := []string{layoutTemplate}

	for _, file := range templateFiles {
		targetTemplate := filepath.Join(templateDir, file)
		targetTemplates = append(targetTemplates, targetTemplate)
	}

	templates, err := template.New("template").Funcs(getFuncMap()).ParseFS(templatesFS, targetTemplates...)

	if err != nil {
		response.RenderError(w, response.ServerError(err))
		return
	}

	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, layoutFile, data); err != nil {
		response.RenderError(w, response.ServerError(err))
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		response.RenderError(w, response.ServerError(err))
		return
	}
}
