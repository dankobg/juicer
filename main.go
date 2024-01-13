package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

//go:embed public/*
var embeddedPublic embed.FS

//go:embed templates/*
var embeddedTemplates embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func main() {
	publicFS, err := fs.Sub(embeddedPublic, "public")
	if err != nil {
		log.Fatalf("failed to get FS subtree out of embedded public files")
	}

	tr := &TemplateRenderer{
		templates: template.Must(template.ParseFS(embeddedTemplates, "templates/*.tmpl")),
	}

	e := echo.New()
	e.Renderer = tr

	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFS)))))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]any{
			"path": "/",
			"data": "wtf shit assa",
		})
	})

	e.GET("/about", func(c echo.Context) error {
		tmplData := map[string]any{"Data": "kurva mac"}
		return c.Render(http.StatusOK, "about", tmplData)
	})

	e.GET("/contact", func(c echo.Context) error {
		tmplData := map[string]any{"Data": "Contact data"}
		return c.Render(http.StatusOK, "contact", tmplData)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
