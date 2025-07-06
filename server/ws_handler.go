package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dankobg/juicer/core"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	orykratos "github.com/ory/client-go"
)

const (
	wsReadBufferSize  = 1024
	wsWriteBufferSize = 1024
)

const juicerClientIDCookieName = "juicer_client_id"

func newClientIDCookie(clientID string) *http.Cookie {
	expires := time.Now().AddDate(10, 0, 0)
	maxAge := int(time.Until(expires).Seconds())

	return &http.Cookie{
		Name:     juicerClientIDCookieName,
		Value:    clientID,
		Path:     "/",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func getClientData(c echo.Context) (core.ClientAuthState, *orykratos.Session, string) {
	sess := GetSession(c.Request().Context())
	var playerID string
	authState := core.ClientGuest

	if sess != nil && sess.Active != nil && *sess.Active {
		playerID = sess.Identity.Id
		authState = core.ClientAuth
	} else {
		cookie, err := c.Cookie(juicerClientIDCookieName)
		if err != nil {
			playerID = uuid.NewString()
		} else {
			playerID = cookie.Value
		}
	}

	return authState, sess, playerID
}

func setCookieHeader(authState core.ClientAuthState, playerID string) http.Header {
	respHeader := http.Header{}
	if authState == core.ClientGuest {
		respHeader.Set("Set-Cookie", newClientIDCookie(playerID).String())
	}
	return respHeader
}

func (h *ApiHandler) serverWs(c echo.Context) error {
	authState, _, clientIDStr := getClientData(c)
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), setCookieHeader(authState, clientIDStr))
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		h.Log.Error("failed to parse client uuid", slog.String("client_id", clientIDStr))
		return err
	}
	client := core.NewClient(clientID, h.Hub, conn, authState, h.Log, c.Request().URL.Query())
	h.Hub.ClientConnected <- client
	go client.ReadLoop()
	go client.WriteLoop()
	return nil
}
