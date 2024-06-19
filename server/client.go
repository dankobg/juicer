package server

import (
	"log/slog"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 55 * time.Second
	maxMessageSize = 1024
)

type clientAuthStatus uint8

const (
	clientAuth clientAuthStatus = iota + 1
	clientAnonymous
)

func (cas clientAuthStatus) String() string {
	switch cas {
	case clientAuth:
		return "auth"
	case clientAnonymous:
		return "anonymous"
	default:
		return "unknown"
	}
}

func (cas clientAuthStatus) anonymous() bool {
	return cas == clientAnonymous
}

func (cas clientAuthStatus) authed() bool {
	return cas == clientAuth
}

type client struct {
	id     string
	roomID string
	hub    *hub
	conn   *websocket.Conn
	send   chan *pb.Message
	log    *slog.Logger
}

func (c *client) String() string {
	return c.id
}

func (c *client) sendErrorMsg(err error) {
	errMsg := &pb.Message{
		Event: &pb.Message_Error{Error: &pb.Error{Message: err.Error()}},
	}
	c.send <- errMsg
}

// readPump reads messages from the websocket connection and passes to hub for processing.
// It is ran in a per-connection goroutine to ensure there is at most one reader
// on a connection by doing all reads from that goroutine
func (c *client) readPump() {
	defer func() {
		c.hub.clientDisconnected <- c
		if err := c.conn.Close(); err != nil {
			c.log.Error("readpump cleanup conn close", slog.Any("error", err))
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.log.Error("conn setreaddeadline", slog.Any("error", err))
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			c.log.Error("ponghandler conn setreaddeadline", slog.Any("error", err))
			return err
		}
		return nil
	})

	for {
		mt, b, err := c.conn.ReadMessage()
		msgType := wsMsgType(mt)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error("websocket unexpected close", slog.Any("error", err))
			}
			break
		}

		msg := &pb.Message{}
		if err := protojson.Unmarshal(b, msg); err != nil {
			c.log.Error("protojson unmarshal", slog.String("type", msgType), slog.String("data", string(b)), slog.Any("error", err))
			c.sendErrorMsg(err)
			continue
		}

		c.log.Debug("readpump recv", slog.String("client_id", c.id), slog.String("msg_type", msgType), slog.String("msg", msg.String()))

		if err := c.hub.processClientMessage(c, msg); err != nil {
			c.log.Error("hub process message", slog.Any("error", err))
			c.sendErrorMsg(err)
			continue
		}
	}
}

// writePump writes messages to the websocket connection.
// It is ran in a per-connection goroutine to ensure there is at most one writer
// on a connection by doing all reads from that goroutine
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			c.log.Error("writepump cleanup conn close", slog.Any("error", err))
		}
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.log.Debug("write pump msg", slog.String("msg", msg.String()))

			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error("conn setwritedeadline", slog.Any("error", err))
			}
			if !ok {
				c.log.Debug("hub closed the channel")
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.log.Debug("conn write close message", slog.Any("error", err))
				}
				return
			}
			b, err := protojson.Marshal(msg)
			if err != nil {
				c.log.Error("writepump marshal", slog.Any("error", err))
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				c.log.Error("writepump conn write", slog.Any("error", err))
			}

			// write queued messages to the websocket conn
			n := len(c.send)
			for i := 0; i < n; i++ {
				b, err := protojson.Marshal(<-c.send)
				if err != nil {
					c.log.Error("writepump queued marshal", slog.Any("error", err))
					continue
				}
				if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
					c.log.Error("writepump conn write buffered", slog.Any("error", err))
				}
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error("ping ticker conn setwritedeadline", slog.Any("error", err))
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.log.Error("ping ticker conn writemessage", slog.Any("error", err))
				return
			}
		}
	}
}

func wsMsgType(typ int) string {
	switch typ {
	case 1:
		return "text"
	case 2:
		return "binary"
	case 8:
		return "close"
	case 9:
		return "ping"
	case 10:
		return "pong"
	default:
		panic("unknown msg type")
	}
}
