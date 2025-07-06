package core

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/url"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
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
	writeWait        = 10 * time.Second
	pingPeriod       = 30 * time.Second
	pongWait         = pingPeriod + time.Second*5
	readMsgLimit     = 1024
	egressBufferSize = 100
)

type client struct {
	id        uuid.UUID
	connID    uuid.UUID
	conn      *websocket.Conn
	authState ClientAuthState
	channels  []Channel
	egress    chan []byte
	hub       *Hub
	log       *slog.Logger
	query     url.Values
}

func NewClient(id uuid.UUID, hub *Hub, conn *websocket.Conn, authState ClientAuthState, sl *slog.Logger, query url.Values) *client {
	l := sl.With(slog.String("client_id", id.String()), slog.String("auth_state", authState.String()))

	return &client{
		id:        id,
		connID:    uuid.New(),
		conn:      conn,
		authState: authState,
		channels:  make([]Channel, 0),
		egress:    make(chan []byte, egressBufferSize),
		hub:       hub,
		log:       l,
		query:     query,
	}
}

func (c *client) String() string {
	return c.id.String()
}

func (c *client) sendMessage(message *pb.Message) {
	b, err := protojson.Marshal(message)
	if err != nil {
		c.log.Error("failed to marshal msg", slog.Any("error", err))
		return
	}
	c.egress <- b
}

func (c *client) sendProblemMsg(err error) {
	errMsg := &pb.Message{Event: &pb.Message_Problem{Problem: &pb.Problem{Message: err.Error()}}}
	c.sendMessage(errMsg)
}

// ReadLoop reads messages from the websocket connection and passes to hub for processing.
// It is ran in a per-connection goroutine to ensure there is at most one reader
// on a connection by doing all reads from that goroutine
func (c *client) ReadLoop() {
	defer func() {
		c.hub.ClientDisconnected <- c

		clientDisconnectedMsg := &pb.Message{Event: &pb.Message_ClientDisconnected{ClientDisconnected: &pb.ClientDisconnected{Id: c.id.String()}}}
		b, err := protojson.Marshal(clientDisconnectedMsg)
		if err != nil {
			c.log.Error("failed to protojson marshal Message_ClientDisconnected message", slog.Any("error", err))
		} else {
			if err := c.hub.rdb.Publish(context.Background(), "ipc", b).Err(); err != nil {
				c.log.Error("failed to publish ipc Message_ClientDisconnected message", slog.Any("error", err))
			}
		}

		if err := c.conn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			c.log.Error("readLoop: conn close", slog.Any("error", err))
		}
	}()

	c.conn.SetReadLimit(readMsgLimit)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.log.Error("readLoop: conn setreaddeadline", slog.Any("error", err))
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			c.log.Error("readLoop: pong conn setreaddeadline", slog.Any("error", err))
			return err
		}
		return nil
	})

	for {
		mt, b, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error("readLoop: ws unexpected close", slog.Any("error", err))
			} else {
				c.log.Debug("readLoop: ws expected close", slog.Any("error", err))
			}
			break
		}
		if mt != websocket.TextMessage {
			c.log.Error("readLoop: unsupported msg type", slog.String("message_type", formatMsgType(mt)), slog.Any("error", err))
			break
		}
		if err := c.hub.processClientMessage(c, b); err != nil {
			c.log.Error("readLoop: processClientMessage", slog.Any("error", err))
			c.sendProblemMsg(err)
			continue
		}
	}
}

// WriteLoop writes messages to the websocket connection.
// It is ran in a per-connection goroutine to ensure there is at most one writer
// on a connection by doing all reads from that goroutine
func (c *client) WriteLoop() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			c.log.Error("writeLoop: conn close", slog.Any("error", err))
		}
	}()

	for {
		select {
		case msg, ok := <-c.egress:
			c.log.Debug("writeLoop: recv from egress", slog.String("msg", string(msg)))

			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error("writeLoop: conn setwritedeadline", slog.Any("error", err))
			}
			if !ok {
				c.log.Debug("writeLoop: hub closed the channel")
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					c.log.Debug("writeLoop: conn writemessage", slog.Any("error", err))
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.log.Error("writeLoop: conn writemessage", slog.Any("error", err))
			}

			// write queued messages to the websocket conn
			n := len(c.egress)
			for range n {
				if err := c.conn.WriteMessage(websocket.TextMessage, <-c.egress); err != nil {
					c.log.Error("writeLoop: conn WriteMessage queued", slog.Any("error", err))
				}
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error("writeLoop: ping conn setwritedeadline", slog.Any("error", err))
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.log.Error("writeLoop: ping conn writemessage", slog.Any("error", err))
				return
			}
		}
	}
}

func formatMsgType(msgType int) string {
	switch msgType {
	case websocket.TextMessage:
		return "text"
	case websocket.BinaryMessage:
		return "binary"
	case websocket.CloseMessage:
		return "close"
	case websocket.PingMessage:
		return "ping"
	case websocket.PongMessage:
		return "pong"
	default:
		return "unknown"
	}
}
