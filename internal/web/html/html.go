package html

import (
	"bytes"
	"embed"
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

	funcMap := template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"url": func(s string) template.URL {
			return template.URL(s)
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"jsstr": func(s string) template.JSStr {
			return template.JSStr(s)
		},
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
	}

	templates, err := template.New("template").Funcs(funcMap).ParseFS(templatesFS, targetTemplates...)

	if err != nil {
		response.ServerError(w, "Parse templates", err)
		return
	}

	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, layoutFile, data); err != nil {
		response.ServerError(w, "Execute template", err)
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		response.ServerError(w, "Write to buffer", err)
		return
	}
}
