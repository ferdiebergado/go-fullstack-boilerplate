package html

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
)

//go:embed templates/*
var templatesFS embed.FS

type Template struct {
	templates map[string]*template.Template
}

func NewTemplate(cfg *config.HTMLTemplateConfig) *Template {
	templateDir := cfg.TemplateDir + "/"
	layoutFile := cfg.LayoutFile
	pagesDir := cfg.PagesDir + "/"
	templatePagesDir := templateDir + pagesDir

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

	layoutTmpl := template.Must(template.New("layout").Funcs(funcMap).ParseFS(templatesFS, templateDir+layoutFile))

	parseTemplate := func(htmlFile string) *template.Template {
		return template.Must(template.Must(layoutTmpl.Clone()).ParseFS(templatesFS, htmlFile))
	}

	return &Template{
		templates: map[string]*template.Template{
			"dbstats": parseTemplate(templatePagesDir + "dbstats.html"),
			"404":     parseTemplate(templatePagesDir + "404.html"),
		},
	}
}

func (t *Template) Render(w http.ResponseWriter, data any, name string) {
	tmpl, ok := t.templates[name]

	if !ok {
		response.RenderError(w, response.TemplateNotFoundError(name))
		return
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		response.RenderError(w, response.ServerError(fmt.Errorf("execute template: %v", err)))
		return
	}

	_, err := buf.WriteTo(w)

	if err != nil {
		response.RenderError(w, response.ServerError(err))
		return
	}
}
