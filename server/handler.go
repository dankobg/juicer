package server

import (
	"log/slog"
	"net/http"

	"github.com/dankobg/juicer/keto"
	"github.com/dankobg/juicer/kratos"
	"github.com/labstack/echo/v4"
)

type ApiHandler struct {
	Log    *slog.Logger
	Kratos *kratos.Client
	Keto   *keto.Client
	Hub    *Hub
	Echo   *echo.Echo
}

func NewApiHandler(log *slog.Logger, kratos *kratos.Client, keto *keto.Client, hub *Hub) *ApiHandler {
	e := echo.New()
	e.HideBanner = true

	return &ApiHandler{
		Log:    log,
		Kratos: kratos,
		Keto:   keto,
		Echo:   echo.New(),
		Hub:    hub,
	}
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Echo.ServeHTTP(w, r)
}
