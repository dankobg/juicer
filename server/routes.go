package server

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(mux *echo.Echo, h *ApiHandler, publicFiles fs.FS) {
	mux.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFiles)))))

	mux.GET("/ws", h.ws)

	api := mux.Group("/api/v1")

	api.GET("/health/alive", h.health)
	api.GET("/health/ready", h.ready)
}
