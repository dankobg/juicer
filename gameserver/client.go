package gameserver

import (
	"net"

	"github.com/google/uuid"
)

type Client struct {
	ID        string
	Conn      net.Conn
	NetPlayer *NetPlayer
}

func (c *Client) String() string {
	return c.ID
}

func NewClient(conn net.Conn, np *NetPlayer) *Client {
	return &Client{
		ID:        uuid.NewString(),
		Conn:      conn,
		NetPlayer: np,
	}
}

type NetPlayer struct {
	Name      string
	Color     string
	Anonymous bool
}

func NewNetPlayer(name string) *NetPlayer {
	return &NetPlayer{
		Name: name,
	}
}

func (np *NetPlayer) SetAnonymous() {
	np.Anonymous = true
}
