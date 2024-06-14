package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dankobg/juicer/random"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const juicerClientIDCookieName = "juicer_client_id"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func (h *ApiHandler) serverWs(c echo.Context) error {
	sess := GetSession(c.Request().Context())
	var clientID string
	authStatus := clientAnonymous

	if sess != nil && sess.Active != nil && *sess.Active {
		clientID = sess.Identity.Id
		authStatus = clientAuth
	} else {
		cookie, err := c.Cookie(juicerClientIDCookieName)
		if err != nil {
			clientID = random.AlphaNumeric(32)
		} else {
			clientID = cookie.Value
		}
	}

	h.Log.Debug("ws client connected", slog.String("client_id", clientID), slog.String("auth_status", authStatus.String()))

	respHeader := http.Header{}
	if authStatus.Anonymous() {
		respHeader.Set("Set-Cookie", newClientIDCookie(clientID).String())
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), respHeader)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan *Message, 256),
		Hub:  h.Hub,
		Log:  h.Log,
	}

	h.Hub.ClientConnected <- client

	go client.writePump()
	go client.readPump()

	return nil
}
