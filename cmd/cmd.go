package cmd

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/core"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/dankobg/juicer/loggerx"
	"github.com/dankobg/juicer/mailer"
	"github.com/dankobg/juicer/server"
	"github.com/dankobg/juicer/store/postgres"
	"github.com/dankobg/juicer/store/redis"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (tr *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

func Run(publicFiles, templateFiles fs.FS) error {
	engine.InitPrecalculatedTables()

	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	logger := loggerx.New(
		loggerx.WithConsolePretty(cfg.ENV != "production" && cfg.Logger.Pretty),
		loggerx.WithLevel(slog.LevelDebug),
	)

	publicFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return fmt.Errorf("failed to get FS subtree out of embedded public files dir: %w", err)
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
		mailer.WithFromName(cfg.Email.FromName),
		mailer.WithFromAddress(cfg.Email.FromAddress),
		mailer.WithLog(logger),
	)

	rdb, err := redis.Connect(cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	rdbStore := redis.New(rdb)

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)
	ketoClient, err := keto.NewClient()
	if err != nil {
		return err
	}

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	pgStore := postgres.New(db)

	store := struct {
		*redis.RedisStore
		*postgres.PgStore
	}{rdbStore, pgStore}

	hub := core.NewHub(store, rdb, logger)

	apiHandler := server.NewApiHandler(cfg, logger, rdb, kratosClient, ketoClient, smtpClient, hub, store)
	apiHandler.Echo.Renderer = tr

	hub.SetUserInfoFetcher(apiHandler)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	srv := server.NewServer(
		server.WithHostPort("", cfg.Port),
		server.WithHandler(apiHandler),
		server.WithReadTimeout(cfg.Server.ReadTimeout),
		server.WithReadHeaderTimeout(cfg.Server.ReadHeaderTimeout),
		server.WithWriteTimeout(cfg.Server.WriteTimeout),
		server.WithIdleTimeout(cfg.Server.IdleTimeout),
		server.WithErrorSlog(logger, slog.LevelDebug),
		server.WithBaseContext(func(l net.Listener) context.Context { return rootCtx }),
	)

	logger.Info("juicer info", slog.String("env", cfg.ENV), slog.String("website_url", cfg.WebsiteURL), slog.String("logger_level", cfg.Logger.Level), slog.Bool("mailer_enabled", cfg.Email.Enabled))

	server.SetupRoutes(apiHandler.Echo, apiHandler, publicFS)

	srvErr := make(chan error, 1)

	go func() {
		logger.Info("server is listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", slog.Any("error", err))
			srvErr <- err
		}
	}()

	go func() {
		if err := apiHandler.Hub.Run(rootCtx); err != nil {
			logger.Error("failed to run hub", slog.Any("error", err))
		}
	}()

	select {
	case <-rootCtx.Done():
		logger.Info("received interruption signal")
	case err := <-srvErr:
		logger.Error("received server err", slog.Any("error", err))
		stop()
	}

	logger.Info("starting shutdown", slog.Duration("graceful_timeout", cfg.Server.GracefulTimeout))
	shutdownCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.Server.GracefulTimeout, fmt.Errorf("graceful shutdown timeout"))
	defer cancel()

	var shutdownErrors []error

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
		shutdownErrors = append(shutdownErrors, err)
	}

	if err := apiHandler.Hub.Stop(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
		shutdownErrors = append(shutdownErrors, err)
	}

	logger.Info("server shut down")
	return errors.Join(shutdownErrors...)
}
