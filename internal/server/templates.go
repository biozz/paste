package server

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/biozz/paste/web"
	"github.com/labstack/echo/v4"
)

type TemplatesLoader struct {
	dev       bool
	templates map[string]*template.Template
}

func NewTemplatesLoader(dev bool) *TemplatesLoader {
	t := TemplatesLoader{dev: dev}
	t.Load()
	return &t
}

type CustomRenderer struct {
	dev    bool
	loader *TemplatesLoader
}

func NewRenderer(dev bool, loader *TemplatesLoader) echo.Renderer {
	r := CustomRenderer{dev, loader}
	return &r
}

func (t *TemplatesLoader) Load() {
	var webFS fs.FS
	webFS = web.TemplatesFS
	if t.dev {
		webFS = os.DirFS("web")
	}
	t.templates = map[string]*template.Template{
		"index": template.Must(template.New("").ParseFS(webFS, "templates/index.go.html", "templates/layout.go.html")),
		"view":  template.Must(template.New("").ParseFS(webFS, "templates/view.go.html", "templates/layout.go.html")),
	}

}

var errTemplateNotFound = errors.New("template not found")

func (t *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.dev {
		t.loader.Load()
	}
	tmpl, ok := t.loader.templates[name]
	if !ok {
		return errTemplateNotFound
	}
	layout := tmpl.Lookup("layout.go.html")
	if layout == nil {
		return tmpl.ExecuteTemplate(w, name, data)
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}
