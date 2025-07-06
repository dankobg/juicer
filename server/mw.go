package server

import (
	"cmp"
	"context"
	"errors"
	"expvar"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	bodyLimitConfig = middleware.BodyLimitConfig{
		Skipper: middleware.DefaultSkipper,
		Limit:   "5M",
	}

	corsConfig = middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3974", "https://juicer-dev.xyz"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Cache-Control", "Content-Language", "Content-Type", "Content-Range", "Expires", "Last-Modified", "Pragma", "Authorization"},
		MaxAge:           96400,
		AllowCredentials: true,
	}

	rateLimiterConfig = middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 10, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			id := c.RealIP()
			return id, nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return &echo.HTTPError{
				Code:     middleware.ErrExtractorError.Code,
				Message:  middleware.ErrExtractorError.Message,
				Internal: err,
			}
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return &echo.HTTPError{
				Code:     middleware.ErrRateLimitExceeded.Code,
				Message:  middleware.ErrRateLimitExceeded.Message,
				Internal: err,
			}
		},
	}

	reqLoggerConfig = func(ctx context.Context, log *slog.Logger) middleware.RequestLoggerConfig {
		return middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			HandleError: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				log.LogAttrs(ctx, slog.LevelDebug, "req",
					slog.String("id", v.RequestID),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("host", v.Host),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.String("user_agent", v.UserAgent),
					slog.Int("status", v.Status),
					slog.String("referer", v.Referer),
					slog.String("err", cmp.Or(v.Error, errors.New("")).Error()),
					slog.Duration("latency", v.Latency),
					slog.String("bytes_in", v.ContentLength),
					slog.Int64("bytes_out", v.ResponseSize),
				)

				return nil
			},
		}
	}
)

var numRequestsVar = expvar.NewInt("num_requests")

func requestsCount(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		numRequestsVar.Add(1)
		return next(c)
	}
}
