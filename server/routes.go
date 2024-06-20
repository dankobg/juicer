package server

import (
	"expvar"
	"io/fs"
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(mux *echo.Echo, h *ApiHandler, publicFiles fs.FS) {
	mux.Use(middleware.Recover())
	mux.Use(middleware.RequestID())
	mux.Use(middleware.BodyLimitWithConfig(bodyLimitConfig))
	mux.Use(middleware.CORSWithConfig(corsConfig))
	mux.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))
	mux.Use(h.AttachSessionData)
	mux.Use(requestsCount)

	// expvar routes
	mux.GET("/debug/vars", echo.WrapHandler(expvar.Handler()))

	// pprof routes
	mux.GET("/debug/pprof/", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	mux.GET("/debug/pprof/allocs", echo.WrapHandler(pprof.Handler("allocs")))
	mux.GET("/debug/pprof/block", echo.WrapHandler(pprof.Handler("block")))
	mux.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	mux.GET("/debug/pprof/goroutine", echo.WrapHandler(pprof.Handler("goroutine")))
	mux.GET("/debug/pprof/heap", echo.WrapHandler(pprof.Handler("heap")))
	mux.GET("/debug/pprof/mutex", echo.WrapHandler(pprof.Handler("mutex")))
	mux.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	mux.POST("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	mux.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	mux.GET("/debug/pprof/threadcreate", echo.WrapHandler(pprof.Handler("threadcreate")))
	mux.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))

	// static files
	mux.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFiles)))))

	// websocket
	mux.GET("/ws", h.serverWs)

	// api v1
	api := mux.Group("/api/v1")

	api.GET("/health/alive", h.health)
	api.GET("/health/ready", h.ready)
}
