package server

import (
	"log/slog"
	"net/http"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/config"
	"github.com/dankobg/juicer/core"
	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/dankobg/juicer/mailer"
	"github.com/dankobg/juicer/store"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var _ api.StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	Cfg      *config.Config
	Log      *slog.Logger
	Kratos   *kratos.Client
	Keto     *keto.Client
	Hub      *core.Hub
	Echo     *echo.Echo
	Rdb      *redis.Client
	store    store.Store
	Mailer   mailer.Mailer
	upgrader websocket.Upgrader
}

func NewApiHandler(cfg *config.Config, log *slog.Logger, rdb *redis.Client, kratos *kratos.Client, keto *keto.Client, mailer mailer.Mailer, hub *core.Hub, st store.Store) *ApiHandler {
	e := echo.New()
	e.HideBanner = true

	upgrader := websocket.Upgrader{
		ReadBufferSize:  wsReadBufferSize,
		WriteBufferSize: wsWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			if cfg.ENV == "production" {
				origin := r.Header.Get("Origin")
				return origin == cfg.Host
			}
			return true
		},
	}

	return &ApiHandler{
		Cfg:      cfg,
		Log:      log,
		Kratos:   kratos,
		Keto:     keto,
		Echo:     e,
		Hub:      hub,
		Rdb:      rdb,
		store:    st,
		Mailer:   mailer,
		upgrader: upgrader,
	}
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Echo.ServeHTTP(w, r)
}
