package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *ApiHandler) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"alive": true,
	})
}

func (h *ApiHandler) ready(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"ready": true,
	})
}
