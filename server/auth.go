package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	orykratos "github.com/ory/client-go"
)

type contextKey string

const (
	oryKratosSessionCtxKey contextKey = "ory_kratos_session_key"
)

const (
	oryKratosCsrfCookiePrefix  = "csrf_token"
	oryKratosSessionCookieName = "ory_kratos_session"

	prefixBearer = "Bearer"
)

type authHeadersResult struct {
	csrf         *http.Cookie
	session      *http.Cookie
	authHeader   string
	cookieHeader string
}

func ExtractAuthHeadersFromRequest(r *http.Request) *authHeadersResult {
	result := &authHeadersResult{
		cookieHeader: r.Header.Get(echo.HeaderCookie),
		authHeader:   r.Header.Get(echo.HeaderAuthorization),
	}

	for _, c := range r.Cookies() {
		if c != nil {
			if ok := strings.HasPrefix(c.Name, oryKratosCsrfCookiePrefix); ok {
				result.csrf = c
			}
		}
	}

	sessionCookie, _ := r.Cookie(oryKratosSessionCookieName)
	if sessionCookie != nil {
		result.session = sessionCookie
	}

	return result
}

func WithSession(ctx context.Context, sess *orykratos.Session) context.Context {
	return context.WithValue(ctx, oryKratosSessionCtxKey, sess)
}

func GetSession(ctx context.Context) *orykratos.Session {
	sess, ok := ctx.Value(oryKratosSessionCtxKey).(*orykratos.Session)
	if !ok {
		return nil
	}
	return sess
}

func (h *ApiHandler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := GetSession(c.Request().Context())

		if sess == nil {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "session is required"}
		}

		if sess.Active != nil && !*sess.Active {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "session is invalid or has expired already"}
		}

		return next(c)
	}
}

func (h *ApiHandler) RequireAnonymous(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := GetSession(c.Request().Context())

		if sess != nil && sess.Active != nil && *sess.Active {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "must have no session"}
		}

		return next(c)
	}
}

func (h *ApiHandler) AttachSessionData(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		info := ExtractAuthHeadersFromRequest(c.Request())

		var hasAuthHeader bool
		var hasCookieHeader bool

		if info.authHeader != "" && strings.HasPrefix(info.authHeader, prefixBearer) {
			hasAuthHeader = true
		}

		if info.cookieHeader != "" {
			hasCookieHeader = true
		}

		if !hasAuthHeader && !hasCookieHeader {
			return next(c)
		}

		req := h.Kratos.Public.FrontendAPI.ToSession(ctx).Cookie(info.session.String())
		session, _, err := req.Execute()
		if err != nil {
			return next(c)
		}

		if session != nil && session.Active != nil && !*session.Active {
			return next(c)
		}

		if session != nil {
			ctxWithSession := WithSession(ctx, session)
			c.SetRequest(c.Request().WithContext(ctxWithSession))
		}

		return next(c)
	}
}
