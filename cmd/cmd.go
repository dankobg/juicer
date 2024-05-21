package cmd

// package main

// import (
// 	"embed"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"io/fs"
// 	"log"
// 	"net/http"

// 	"github.com/dankobg/juicer/config"
// 	"github.com/labstack/echo/v4"
// )

// //go:embed public/*
// var embeddedPublic embed.FS

// //go:embed templates/*
// var embeddedTemplates embed.FS

// type TemplateRenderer struct {
// 	templates *template.Template
// }

// func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
// 	return tr.templates.ExecuteTemplate(w, name, data)
// }

// func main() {
// 	cfg, _, err := config.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	publicFS, err := fs.Sub(embeddedPublic, "public")
// 	if err != nil {
// 		log.Fatalf("failed to get FS subtree out of embedded public files")
// 	}

// 	tr := &TemplateRenderer{
// 		templates: template.Must(template.ParseFS(embeddedTemplates, "templates/*.tmpl")),
// 	}

// 	e := echo.New()
// 	e.Renderer = tr

// 	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFS)))))

// 	e.GET("/api/v1/", func(c echo.Context) error {
// 		return c.JSON(http.StatusOK, map[string]any{
// 			"path": "/api/v1/",
// 			"data": "wtf rofl",
// 		})
// 	})

// 	e.GET("/api/v1/health/alive", func(c echo.Context) error {
// 		return c.JSON(http.StatusOK, map[string]any{
// 			"health": "alive",
// 		})
// 	})

// 	fmt.Printf("%+v\n", cfg.App.Port)
// 	fmt.Printf("%+v\n", cfg)

// 	e.Logger.Fatal(e.Start(":1337"))
// }

func Run() error {
	return nil
}
