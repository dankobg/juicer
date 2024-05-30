package server

import (
	"log/slog"
	"net/http"
	"time"
)

func NewServer(log *slog.Logger, h http.Handler) *http.Server {
	srv := http.Server{
		Addr:              ":1337",
		Handler:           h,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       90 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelDebug),
	}

	return &srv
}
