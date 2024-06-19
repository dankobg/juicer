package server

import (
	"log/slog"
	"net/http"

	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/dankobg/juicer/mailer"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type ApiHandler struct {
	Log    *slog.Logger
	Kratos *kratos.Client
	Keto   *keto.Client
	Hub    *hub
	Echo   *echo.Echo
	Rdb    *redis.Client
	Mailer mailer.Mailer
}

func NewApiHandler(log *slog.Logger, rdb *redis.Client, kratos *kratos.Client, keto *keto.Client, mailer mailer.Mailer, hub *hub) *ApiHandler {
	e := echo.New()
	e.HideBanner = true

	return &ApiHandler{
		Log:    log,
		Kratos: kratos,
		Keto:   keto,
		Echo:   e,
		Hub:    hub,
		Rdb:    rdb,
		Mailer: mailer,
	}
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Echo.ServeHTTP(w, r)
}
