package server

import (
	"expvar"
	"io/fs"
	"net/http"
	"net/http/pprof"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

func SetupRoutes(e *echo.Echo, h *ApiHandler, publicFiles fs.FS) {
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimitWithConfig(bodyLimitConfig))
	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))
	e.Use(h.AttachSessionData)
	e.Use(requestsCount)

	// expvar routes
	e.GET("/debug/vars", echo.WrapHandler(expvar.Handler()))

	// pprof routes
	e.GET("/debug/pprof/", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/allocs", echo.WrapHandler(pprof.Handler("allocs")))
	e.GET("/debug/pprof/block", echo.WrapHandler(pprof.Handler("block")))
	e.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	e.GET("/debug/pprof/goroutine", echo.WrapHandler(pprof.Handler("goroutine")))
	e.GET("/debug/pprof/heap", echo.WrapHandler(pprof.Handler("heap")))
	e.GET("/debug/pprof/mutex", echo.WrapHandler(pprof.Handler("mutex")))
	e.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	e.POST("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	e.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	e.GET("/debug/pprof/threadcreate", echo.WrapHandler(pprof.Handler("threadcreate")))
	e.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))

	// static files
	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", http.FileServer(http.FS(publicFiles)))))

	// websocket
	e.GET("/ws", h.serverWs)

	v1Group := e.Group("/api/v1")

	v1Group.POST("/webhooks/kratos/registration_after_password", h.registrationAfterPassword)
	v1Group.POST("/webhooks/kratos/registration_after_oidc", h.registrationAfterOidc)

	apiSrv := api.NewStrictHandler(h, make([]strictecho.StrictEchoMiddlewareFunc, 0))
	api.RegisterHandlers(v1Group, apiSrv)
}
