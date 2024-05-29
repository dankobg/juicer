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

// readPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine and
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Hub.ClientDisconnected <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Log.Error("UnexpectedCloseError", err)
			}
			break
		}

		c.Log.Info("readPump RECV MSG", slog.String("type", msg.Type), slog.String("data", string(msg.Data)))
		c.Hub.Broadcast <- &Message{Type: "server_echo", Data: msg.Data}
	}
}

// writePump pumps messages from the hub to the websocket connection.
// The application runs writePump in a per-connection goroutine and
// ensures that there is at most one writer on a connection by executing all
// writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Log.Info("hub closed the channel")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteJSON(&msg); err != nil {
				slog.Error("failed to write json msg to conn", err)
				return
			}

			// add queued messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if err := c.Conn.WriteJSON(<-c.Send); err != nil {
					return
				}
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("ticker conn WriteDeadlineErr", err)
				return
			}
		}
	}
}
