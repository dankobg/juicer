package cmd

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/dankobg/juicer/mailer"
	"github.com/dankobg/juicer/redis"
	"github.com/dankobg/juicer/server"
	"github.com/labstack/echo/v4"
)

var (
	numGoroutinesVar = expvar.NewInt("goroutines")
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
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	publicFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return fmt.Errorf("failed to get FS subtree out of embedded public files dir")
	}

	tr := &TemplateRenderer{
		templates: template.Must(template.ParseFS(templateFiles, "templates/*.tmpl")),
	}

	smtpClient := mailer.NewSmtpClient(
		mailer.WithEnabled(cfg.ENV == "production"),
		mailer.WithDevHost(cfg.Email.DevSMTPHost),
		mailer.WithDevPort(cfg.Email.DevSMTPPort),
		mailer.WithDevUsername(cfg.Email.DevSMTPUsername),
		mailer.WithDevPassword(cfg.Email.DevSMTPPassword),
		mailer.WithHost(cfg.Email.SMTPHost),
		mailer.WithPort(cfg.Email.SMTPPort),
		mailer.WithUsername(cfg.Email.SMTPUsername),
		mailer.WithPassword(cfg.Email.SMTPPassword),
		mailer.WithTLS(true),
		mailer.WithFromName("Juicer"),
		mailer.WithFromAddress("juicer@chess.com"),
		mailer.WithLog(logger),
	)

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to initialize redis client")
	}

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)
	ketoClient := keto.NewClient()

	hub := server.NewHub(logger, rdb)

	apiHandler := server.NewApiHandler(logger, rdb, kratosClient, ketoClient, smtpClient, hub)
	apiHandler.Echo.Renderer = tr

	srv := server.NewServer(
		server.WithHostPort("", cfg.Port),
		server.WithHandler(apiHandler),
		server.WithReadTimeout(cfg.Server.ReadTimeout),
		server.WithReadHeaderTimeout(cfg.Server.ReadHeaderTimeout),
		server.WithWriteTimeout(cfg.Server.WriteTimeout),
		server.WithIdleTimeout(cfg.Server.IdleTimeout),
		server.WithErrorSlog(logger, slog.LevelDebug),
	)

	logger.Info("juicer info", slog.String("env", cfg.ENV), slog.String("website_url", cfg.WebsiteURL))

	server.SetupRoutes(apiHandler.Echo, apiHandler, publicFS)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	go func() {
		logger.Info("server is listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	go func() {
		logger.Info("hub is running")
		apiHandler.Hub.Run(context.TODO())
	}()

	go func() {
		numGoroutinesVar.Set(int64(runtime.NumGoroutine()))
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				numGoroutinesVar.Set(int64(runtime.NumGoroutine()))
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	logger.Info("received interruption signal, starting shutdown process")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
		return err
	}

	// close other resources

	logger.Info("server shut down")
	return nil
}
