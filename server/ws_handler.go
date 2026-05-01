package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/dankobg/juicer/ws"
	"github.com/google/uuid"
	orykratos "github.com/ory/client-go"
)

const juicerClientIDCookieName = "juicer_user_id"

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

func getClientAuthData(r *http.Request) (string, ws.ClientAuthState, *orykratos.Session) {
	sess := GetSession(r.Context())

	var userID string

	authState := ws.ClientGuest

	if sess != nil && sess.Active != nil && *sess.Active {
		userID = sess.Identity.Id
		authState = ws.ClientAuth
	} else {
		cookie, err := r.Cookie(juicerClientIDCookieName)
		if err != nil {
			userID = uuid.NewString()
		} else {
			userID = cookie.Value
		}
	}

	return userID, authState, sess
}

func setClientIDCookie(w http.ResponseWriter, authState ws.ClientAuthState, playerID string) {
	if authState == ws.ClientGuest {
		w.Header().Set("Set-Cookie", newClientIDCookie(playerID).String())
	}
}

func (a *ApiHandler) serverWs(w http.ResponseWriter, r *http.Request) {
	clientID, authState, _ := getClientAuthData(r)

	setClientIDCookie(w, authState, clientID)

	var (
		onPingReceived func(ctx context.Context, payload []byte) bool
		onPongReceived func(ctx context.Context, payload []byte)
	)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: a.Cfg.ENV == "development",
		OriginPatterns:     a.Cfg.Cors.AllowOrigins,
		OnPingReceived: func(ctx context.Context, payload []byte) bool {
			if onPingReceived != nil {
				return onPingReceived(ctx, payload)
			}

			return false
		},
		OnPongReceived: func(ctx context.Context, payload []byte) {
			if onPongReceived != nil {
				onPongReceived(ctx, payload)
			}
		},
	})
	if err != nil {
		a.Log.Error("websocket.Accept", slog.String("user_id", clientID), slog.Any("error", err))
		return
	}

	defer func() {
		_ = conn.CloseNow()
	}()

	userID, err := uuid.Parse(clientID)
	if err != nil {
		a.Log.Error("failed to parse client uuid", slog.String("user_id", clientID))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())

	client := ws.NewClient(
		userID,
		a.Hub,
		conn,
		authState,
		a.Log,
		r.URL.Query(),
		cancel,
		func(fn func(ctx context.Context, payload []byte) bool) {
			onPingReceived = fn
		},
		func(fn func(ctx context.Context, payload []byte)) {
			onPongReceived = fn
		},
	)

	initialChannels, err := a.Hub.InitializeChannels(r.Context(), client)
	if err != nil {
		a.Log.Error("failed to initialize channels", slog.String("user_id", clientID))
		return
	}

	channels := make([]ws.Channel, len(initialChannels))
	for i, channel := range initialChannels {
		channels[i] = ws.Channel(channel)
	}

	client.JoinChannels(channels)

	a.Hub.ClientConnected <- client

	go client.ReadLoop(ctx)

	client.WriteLoop(ctx)
}
