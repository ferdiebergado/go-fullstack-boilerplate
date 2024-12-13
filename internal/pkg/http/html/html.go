package html

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
)

//go:embed templates/*
var templatesFS embed.FS

func Render(w http.ResponseWriter, data any, templateFiles ...string) {
	const (
		templateDir      = "templates"
		layoutFile       = "layout.html"
		partialTemplates = "partials/*.html"
	)

	layoutTemplate := filepath.Join(templateDir, layoutFile)

	targetTemplates := []string{layoutTemplate}

	for _, file := range templateFiles {
		targetTemplate := filepath.Join(templateDir, file)
		targetTemplates = append(targetTemplates, targetTemplate)
	}

	templates, err := template.New("template").Funcs(getFuncMap()).ParseFS(templatesFS, targetTemplates...)

	if err != nil {
		response.RenderServerError(w, fmt.Errorf("parse template: %v", err))
		return
	}

	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, layoutFile, data); err != nil {
		response.RenderServerError(w, fmt.Errorf("execute template: %v", err))
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		response.RenderServerError(w, fmt.Errorf("write to buffer: %v", err))
		return
	}
}
