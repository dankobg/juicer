package server

import (
	"fmt"
	"net/http"

	"github.com/dankobg/juicer/random"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ApiHandler) ws(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	id := random.AlphaNumeric(32)

	client := &Client{
		ID:   id,
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
