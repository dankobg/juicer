package server

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 55 * time.Second
	maxMessageSize = 1024
)

type Client struct {
	ID     string
	RoomID string
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan *Message
	Log    *slog.Logger
}

func (c *Client) String() string {
	return c.ID
}

func (c *Client) sendErrorMsg(err error) {
	c.Send <- &Message{Type: "error", Data: []byte(err.Error())}
}

// readPump reads messages from the websocket connection and passes to hub for processing.
// It is ran in a per-connection goroutine to ensure there is at most one reader
// on a connection by doing all reads from that goroutine
func (c *Client) readPump() {
	defer func() {
		c.Hub.ClientDisconnected <- c
		if err := c.Conn.Close(); err != nil {
			c.Log.Error("readpump cleanup conn close", slog.Any("error", err))
		}
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.Log.Error("conn setreaddeadline", slog.Any("error", err))
	}
	c.Conn.SetPongHandler(func(string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			c.Log.Error("ponghandler conn setreaddeadline", slog.Any("error", err))
			return err
		}
		return nil
	})

	for {
		msg := &Message{}
		if err := c.Conn.ReadJSON(msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Log.Error("websocket unexpected close", slog.Any("error", err))
			}
			break
		}

		c.Log.Debug("readpump recv", slog.String("client_id", c.ID), slog.String("type", msg.Type), slog.String("data", string(msg.Data)))
		msg.ClientID = c.ID
		if err := c.Hub.ProcessMessage(msg); err != nil {
			c.Log.Error("hub process message", slog.Any("error", err))
			c.sendErrorMsg(err)
			continue
		}
	}
}

// writePump writes messages to the websocket connection.
// It is ran in a per-connection goroutine to ensure there is at most one writer
// on a connection by doing all reads from that goroutine
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			c.Log.Error("writepump cleanup conn close", slog.Any("error", err))
		}
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.Log.Error("conn setwritedeadline", slog.Any("error", err))
			}
			if !ok {
				c.Log.Debug("hub closed the channel")
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.Log.Debug("conn write close message", slog.Any("error", err))
				}
				return
			}
			if err := c.Conn.WriteJSON(&msg); err != nil {
				c.Log.Error("writepump conn write", slog.Any("error", err))
			}

			// write queued messages to the websocket conn
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if err := c.Conn.WriteJSON(<-c.Send); err != nil {
					c.Log.Error("writepump conn write buffered", slog.Any("error", err))
				}
			}

		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.Log.Error("ping ticker conn setwritedeadline", slog.Any("error", err))
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Log.Error("ping ticker conn writemessage", slog.Any("error", err))
				return
			}
		}
	}
}
