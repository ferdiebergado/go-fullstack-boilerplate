package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/errtypes"
	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/http/response"
	"github.com/ferdiebergado/go-fullstack-boilerplate/web"
)

const suffix = ".html"

// Retrieve the template func maps
func funcMap() template.FuncMap {
	return template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s) // #nosec G203 -- No user input
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s) // #nosec G203 -- No user input
		},
		"url": func(s string) template.URL {
			return template.URL(s) // #nosec G203 -- No user input
		},
		"js": func(s string) template.JS {
			return template.JS(s) // #nosec G203 -- No user input
		},
		"jsstr": func(s string) template.JSStr {
			return template.JSStr(s) // #nosec G203 -- No user input
		},
		"css": func(s string) template.CSS {
			return template.CSS(s) // #nosec G203 -- No user input
		},
	}
}

// Parse all partial templates into the layout template
func parsePartials(layoutTmpl *template.Template, partialsDir string) {
	err := fs.WalkDir(web.TemplatesFS, partialsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, suffix) {
			_, parseErr := layoutTmpl.ParseFS(web.TemplatesFS, path)
			if parseErr != nil {
				return fmt.Errorf("parse partials: %w", parseErr)
			}
			slog.Debug("parsed partial", "path", path)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("load partials templates: %w", err))
	}

	slog.Debug("layout", "name", layoutTmpl.Name())
}

// Parse main templates from pagesDir
func parsePages(layoutTmpl *template.Template, templatePagesDir string) templateMap {
	tmplMap := make(templateMap)
	err := fs.WalkDir(web.TemplatesFS, templatePagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, suffix) {
			name := strings.TrimPrefix(path, templatePagesDir+"/")
			tmplMap[name] = template.Must(template.Must(layoutTmpl.Clone()).ParseFS(web.TemplatesFS, path))
			slog.Debug("parsed page", "path", path, "name", name)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("load pages templates: %w", err))
	}

	return tmplMap
}

type templateMap map[string]*template.Template

type Template struct {
	templates templateMap
}

func NewTemplate(cfg *config.HTMLTemplateConfig) *Template {
	if cfg.TemplateDir == "" || cfg.LayoutFile == "" || cfg.PagesDir == "" {
		panic("invalid template configuration: template directory, layout file, and pages directory are required")
	}

	templateDir := cfg.TemplateDir + "/"
	layoutFile := templateDir + cfg.LayoutFile
	partialsDir := templateDir + cfg.PartialsDir
	pagesDir := templateDir + cfg.PagesDir

	layoutTmpl := template.Must(template.New("layout").Funcs(funcMap()).ParseFS(web.TemplatesFS, layoutFile))

	parsePartials(layoutTmpl, partialsDir)

	return &Template{templates: parsePages(layoutTmpl, pagesDir)}
}

func (t *Template) Render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := t.templates[name]

	if !ok {
		err := &TemplateNotFoundError{Template: name}
		response.RenderError(w, nil, errtypes.ServerError(err))
		return
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		execErr := fmt.Errorf("execute template: %w", err)
		response.RenderError(w, nil, errtypes.ServerError(execErr))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := buf.WriteTo(w)

	if err != nil {
		writeErr := fmt.Errorf("write response: %w", err)
		response.RenderError(w, nil, errtypes.ServerError(writeErr))
		return
	}
}
