package cmd

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/dankobg/juicer/redis"
	"github.com/dankobg/juicer/server"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func Run(publicFiles, templateFiles fs.FS) error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", err)
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))

	publicFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return fmt.Errorf("failed to get FS subtree out of embedded public files dir")
	}

	tr := &TemplateRenderer{
		templates: template.Must(template.ParseFS(templateFiles, "templates/*.tmpl")),
	}

	rdb, err := redis.New()
	if err != nil {
		return fmt.Errorf("failed to initialize redis client")
	}

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)
	ketoClient := keto.NewClient()
	hub := server.NewHub(logger)

	apiHandler := server.NewApiHandler(logger, rdb, kratosClient, ketoClient, hub)
	apiHandler.Echo.Renderer = tr

	srv := server.NewServer(logger, apiHandler)

	server.SetupRoutes(apiHandler.Echo, apiHandler, publicFS)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	go func() {
		logger.Info("server is listening", slog.String("url", "https://localhost:1337"))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", err)
			os.Exit(1)
		}
	}()

	go func() {
		logger.Info("hub is running")
		apiHandler.Hub.Run()
	}()

	<-ctx.Done()
	logger.Info("received interruption signal, starting shutdown process")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", err)
		return err
	}

	// close other resources

	logger.Info("server shut down")
	return nil
}
