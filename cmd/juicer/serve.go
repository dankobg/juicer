package juicer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankobg/juicer/auth/keto"
	"github.com/dankobg/juicer/auth/kratos"
	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/httpserver"
	"github.com/dankobg/juicer/logging"
	"github.com/dankobg/juicer/mailer"
	"github.com/dankobg/juicer/persistence/postgres"
	redisp "github.com/dankobg/juicer/persistence/redis"
	"github.com/dankobg/juicer/redis"
	"github.com/dankobg/juicer/server"
	"github.com/dankobg/juicer/ws"
)

type ServeCommand struct{}

func (sc *ServeCommand) Run() error {
	engine.InitPrecalculatedTables()

	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	logger := logging.New(
		logging.WithConsolePretty(cfg.ENV != "production" && cfg.Logger.Pretty),
		logging.WithLevel(slog.LevelDebug),
	)

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

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	rdbPersistor := redisp.New(rdb)

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)

	ketoClient, err := keto.NewClient()
	if err != nil {
		return err
	}

	pool, err := postgres.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("postgres.NewPool: %w", err)
	}
	defer pool.Close()

	pg := postgres.New(pool)

	persistor := struct {
		*redisp.RedisPersistor
		*postgres.PgPersistor
	}{rdbPersistor, pg}

	hub := ws.NewHub(persistor, rdb, logger)

	apiHandler := server.New(cfg, logger, rdb, kratosClient, ketoClient, smtpClient, hub, persistor)

	if err := apiHandler.FetchCategoryThresholds(context.Background()); err != nil {
		return fmt.Errorf("FetchCategoryThresholds: %w", err)
	}
	if err := apiHandler.FetchProtoMappingsCacheLookups(context.Background()); err != nil {
		return fmt.Errorf("FetchProtoMappingsCacheLookups: %w", err)
	}

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	h := apiHandler.SetupRoutes(cfg.ENV, cfg.UploadDir)

	srv := httpserver.New(
		httpserver.WithHostPort("", cfg.Port),
		httpserver.WithHandler(h),
		httpserver.WithReadTimeout(cfg.Server.ReadTimeout),
		httpserver.WithReadHeaderTimeout(cfg.Server.ReadHeaderTimeout),
		httpserver.WithWriteTimeout(cfg.Server.WriteTimeout),
		httpserver.WithIdleTimeout(cfg.Server.IdleTimeout),
		httpserver.WithErrorSlog(logger, slog.LevelDebug),
		httpserver.WithBaseContext(func(l net.Listener) context.Context { return rootCtx }),
	)

	logger.Info("juicer info", slog.String("env", cfg.ENV), slog.String("website_url", cfg.WebsiteURL), slog.String("logger_level", cfg.Logger.Level), slog.Bool("mailer_enabled", cfg.Email.Enabled))

	srvErr := make(chan error, 1)

	go func() {
		logger.Info("server is listening", slog.String("addr", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", slog.Any("error", err))

			srvErr <- err
		}
	}()

	go apiHandler.PubsubProcess(rootCtx)

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

	logger.Info("server shut down")

	return errors.Join(shutdownErrors...)
}
