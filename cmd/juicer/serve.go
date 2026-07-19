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
	"github.com/dankobg/juicer/bus"
	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/features/chat"
	chatpg "github.com/dankobg/juicer/features/chat/persistence/postgres"
	chatrdb "github.com/dankobg/juicer/features/chat/persistence/redis"
	"github.com/dankobg/juicer/features/game"
	gamepg "github.com/dankobg/juicer/features/game/persistence/postgres"
	gamerdb "github.com/dankobg/juicer/features/game/persistence/redis"
	"github.com/dankobg/juicer/features/idp"
	idppg "github.com/dankobg/juicer/features/idp/persistence/postgres"
	"github.com/dankobg/juicer/features/presence"
	presencerdb "github.com/dankobg/juicer/features/presence/persistence/redis"
	"github.com/dankobg/juicer/features/webhooks"
	"github.com/dankobg/juicer/httpserver"
	"github.com/dankobg/juicer/logging"
	"github.com/dankobg/juicer/mailer"
	"github.com/dankobg/juicer/postgres"
	"github.com/dankobg/juicer/redis"
	"github.com/dankobg/juicer/server"
	"github.com/dankobg/juicer/ws"
	"github.com/google/uuid"
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

	pgPst := postgres.NewPgPersistor(pool)

	hub := ws.NewHub(rdb, logger)

	gamePst := gamepg.NewPgGamePersistor(pgPst)
	gvPst := gamepg.NewPgGameVariantPersistor(pgPst)
	gtcPst := gamepg.NewPgGameTimeCategoryPersistor(pgPst)
	gtkPst := gamepg.NewPgGameTimeKindPersistor(pgPst)
	agPst := gamerdb.NewRedisActiveGamePersistor(rdb)
	grPst := gamepg.NewPgGameResultPersistor(pgPst)
	grsPst := gamepg.NewPgGameResultStatusPersistor(pgPst)
	gsPst := gamepg.NewPgGameStatePersistor(pgPst)
	userPst := idppg.NewPgUserPersistor(pgPst)
	ratingPst := gamepg.NewPgRatingPersistor(pgPst)
	poolPst := gamerdb.NewRedisPoolPersistor(rdb)
	presencePst := presencerdb.NewRedisPresencePersistor(rdb)

	bus := bus.NewBus(rdb)

	redisChatPst := chatrdb.NewRedisChatPersistor(rdb)
	sqlChatPst := chatpg.NewPostgresChatPersistor(pgPst)
	presenceSvc := presence.NewPresenceService(bus, presencePst, logger)

	// chatSvc := chat.NewChatService(redisChatPst, sqlChatPst, sqlChatPst, sqlChatPst, logger)
	chatSvc := chat.NewChatService(redisChatPst, redisChatPst, redisChatPst, redisChatPst, logger)
	_ = sqlChatPst
	idpr := idp.NewIdentityProvider(kratosClient, ketoClient, cfg.KratosAPIKey, cfg.KetoAPIKey, userPst, gtcPst, ratingPst, logger)

	pst := game.Persistor{
		Game:             gamePst,
		GameVariant:      gvPst,
		GameTimeCategory: gtcPst,
		GameTimeKind:     gtkPst,
		ActiveGame:       agPst,
		GameResult:       grPst,
		GameResultStatus: grsPst,
		GameState:        gsPst,
		Rating:           ratingPst,
		Pool:             poolPst,
	}

	gameSvc := game.NewGameService(presenceSvc, chatSvc, bus, usrRdr{idpr}, pst, logger)
	webhooksSvc := webhooks.NewWebhooks(idpr, logger)

	apiHandler := server.New(cfg, logger, rdb, kratosClient, ketoClient, smtpClient, hub, gameSvc, chatSvc, idpr, webhooksSvc)

	if err := gameSvc.FetchCategoryThresholds(context.Background()); err != nil {
		return fmt.Errorf("FetchCategoryThresholds: %w", err)
	}

	if err := gameSvc.FetchProtoMappingsCacheLookups(context.Background()); err != nil {
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

	go gameSvc.PubsubProcess(rootCtx)
	go gameSvc.StartMatchmaking(rootCtx)

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

	hub.Stop()
	bus.Close()

	return errors.Join(shutdownErrors...)
}

type usrRdr struct {
	*idp.IdentityProvider
}

func (ur usrRdr) GetUserInfo(ctx context.Context, userID uuid.UUID) (game.UserInfo, error) {
	userInfo, err := ur.IdentityProvider.GetUserInfo(ctx, userID)
	if err != nil {
		return game.UserInfo{}, err
	}

	return game.UserInfo(userInfo), nil
}

func (ur usrRdr) GetUsername(ctx context.Context, userID uuid.UUID) (string, error) {
	return ur.IdentityProvider.GetUsername(ctx, userID)
}
