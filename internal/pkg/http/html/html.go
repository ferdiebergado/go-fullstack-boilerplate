package html

import (
	"bytes"
	"errors"
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

var (
	ErrTemplateParse    = errors.New("failed to parse template")
	ErrTemplateNotFound = errors.New("could not find template")
	ErrTemplateExec     = errors.New("failed to execute template")
	ErrTemplateWrite    = errors.New("failed to write template to the response")
)

const suffix = ".html"

// Retrieve the template func maps
func getFuncMap() template.FuncMap {
	return template.FuncMap{
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
				return fmt.Errorf("failed to parse partial: %w", parseErr)
			}
			slog.Debug("parsed partial", "path", path)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("failed to load partial templates: %w", err))
	}

	slog.Debug("layout", "name", layoutTmpl.Name(), "defined_templates", layoutTmpl.DefinedTemplates())
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
			slog.Debug("parsed page", "path", path, "name", name, "define_templates", tmplMap[name].DefinedTemplates())
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("failed to load templates: %w", err))
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

	layoutTmpl := template.Must(template.New("layout").Funcs(getFuncMap()).ParseFS(web.TemplatesFS, layoutFile))

	parsePartials(layoutTmpl, partialsDir)

	return &Template{templates: parsePages(layoutTmpl, pagesDir)}
}

func (t *Template) Render(w http.ResponseWriter, data any, name string) {
	tmpl, ok := t.templates[name]

	if !ok {
		response.RenderError(w, errtypes.ServerError(fmt.Errorf("%w: %s", ErrTemplateNotFound, name)))
		return
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		response.RenderError(w, errtypes.ServerError(fmt.Errorf("%w %v", ErrTemplateExec, err)))
		return
	}

	_, err := buf.WriteTo(w)

	if err != nil {
		response.RenderError(w, errtypes.ServerError(fmt.Errorf("%w %v", ErrTemplateWrite, err)))
		return
	}
}
