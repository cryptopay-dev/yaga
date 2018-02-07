package doc

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

const (
	defaultDocURL = "doc"
)

type params struct {
	title           string
	docURL          string
	swaggerFilePath string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newHTMLRenderer() (*Template, error) {
	t, err := template.New(defaultDocURL).Parse(htmlTemplate)
	return &Template{
		templates: t,
	}, err
}

func (p *params) apiDoc(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, defaultDocURL, struct {
		Title string
		YAML  string
	}{
		Title: p.title,
		YAML:  swaggerURL(p.docURL),
	})
}

// swaggerYAML - GET /doc/swagger.yaml (set url in config)
func (p *params) swaggerYAML(ctx echo.Context) error {
	return ctx.File(p.swaggerFilePath)
}

// AddDocumentation allows you to add documentation for the url for the swagger file
func AddDocumentation(e *echo.Echo, url, title, swaggerFilePath string) {
	if len(url) == 0 {
		url = defaultDocURL
	}
	renderer, err := newHTMLRenderer()
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.Renderer = renderer
	p := &params{
		title:           title,
		docURL:          url,
		swaggerFilePath: swaggerFilePath,
	}
	e.GET("/"+url, p.apiDoc)
	e.GET("/"+swaggerURL(url), p.swaggerYAML)
}

func swaggerURL(url string) string {
	return url + "/swagger.yaml"
}
