package ws

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type ClientAuthState uint8

const (
	ClientGuest ClientAuthState = iota
	ClientAuth
)

func (cas ClientAuthState) String() string {
	switch cas {
	case ClientAuth:
		return "auth"
	case ClientGuest:
		return "guest"
	default:
		return "unknown"
	}
}

const (
	pingPeriod       = 5 * time.Second
	pongWait         = pingPeriod + time.Second*5
	outMsgBufferSize = 64
)

type client struct {
	id        uuid.UUID
	connID    uuid.UUID
	conn      *websocket.Conn
	hub       *Hub
	authState ClientAuthState
	channels  []Channel
	outMsg    chan []byte
	query     url.Values
	log       *slog.Logger
	closeSlow func()
}

func NewClient(id uuid.UUID, hub *Hub, conn *websocket.Conn, authState ClientAuthState, logger *slog.Logger, query url.Values) *client {
	connID := uuid.New()
	clientLogger := logger.With(slog.String("client_id", id.String()), slog.String("conn_id", connID.String()), slog.String("auth_state", authState.String()))

	return &client{
		id:        id,
		connID:    connID,
		conn:      conn,
		authState: authState,
		channels:  make([]Channel, 0),
		outMsg:    make(chan []byte, outMsgBufferSize),
		hub:       hub,
		log:       clientLogger,
		query:     query,
	}
}

func (c *client) String() string {
	return c.id.String()
}

func (c *client) JoinChannels(channels []Channel) {
	c.channels = channels
}

// // Then try to register the realm with the connection path that was
// // passed in.
// err = registerRealm(client, path, hub)
// if err != nil {
// 	log.Err(err).Msg("register-realm-error")
// 	client.conn.Close()
// }

func (c *client) ReadLoop(ctx context.Context) {
	defer func() {
		c.hub.ClientDisconnected <- c
	}()

	for {
		msgType, msg, err := c.conn.Read(ctx)
		if err != nil {
			c.log.Error("conn.Read", slog.Any("error", err))
			return
		}

		c.log.Info("recv", slog.String("msg_type", msgType.String()), slog.String("msg", string(msg)))

		if err := c.hub.processClientWebsocketMessage(c, msg); err != nil {
			c.log.Error("processClientWebsocketMessage", slog.String("msg_type", msgType.String()), slog.String("msg", string(msg)), slog.Any("error", err))
			// @TODO: some send err msg to client
			continue
		}
	}
}

func (c *client) WriteLoop(ctx context.Context) {
	for {
		select {
		case outMsg, ok := <-c.outMsg:
			if !ok {
				// c.log.Debug("WriteLoop recv from outMsg chan not ok, closed channel")
				return
			}

			c.log.Info("send", slog.String("msg_type", websocket.MessageText.String()), slog.Any("msg", string(outMsg)))

			if err := c.conn.Write(ctx, websocket.MessageText, outMsg); err != nil {
				c.log.Error("conn.Write", slog.Any("error", err))
				return
			}
		case <-ctx.Done():
			c.log.Info("WriteLoop ctx.Done")
			return
		}
	}
}
