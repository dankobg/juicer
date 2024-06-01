package server

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(mux *echo.Echo, h *ApiHandler, publicFiles fs.FS) {
	mux.Use(middleware.Recover())
	mux.Use(middleware.BodyLimitWithConfig(bodyLimitConfig))
	mux.Use(middleware.CORSWithConfig(corsConfig))
	mux.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))
	mux.Use(h.AttachSessionData)

	mux.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFiles)))))

	mux.GET("/ws", h.serverWs)

	api := mux.Group("/api/v1")

	api.GET("/health/alive", h.health)
	api.GET("/health/ready", h.ready)
}
