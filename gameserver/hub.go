package gameserver

import (
	"hash/maphash"
	"math/rand"
	"net"
	"sync"
)

type Hub struct {
	Cons  map[string]*Client
	Rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Cons:  make(map[string]*Client),
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) handleNewConnection(con net.Conn, clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Cons[clientID] = &Client{
		ID:   clientID,
		Conn: con,
	}
}

func (h *Hub) handleDisconnection(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c, has := h.Cons[clientID]; has {
		c.Conn.Close()
		delete(h.Cons, clientID)
	}
}

func (h *Hub) SetupGame(c1, c2 *Client) error {
	r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))

	if r.Intn(2) == 0 {
		c1.NetPlayer.Color = "w"
		c2.NetPlayer.Color = "b"
	} else {
		c1.NetPlayer.Color = "b"
		c2.NetPlayer.Color = "w"
	}

	h.addRoom(c1, c2)

	return nil
}

func (h *Hub) ClientsCount() int {
	return len(h.Cons)
}

func (h *Hub) RoomsCount() int {
	return len(h.Rooms)
}

func (h *Hub) addRoom(c1, c2 *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := NewRoom(c1, c2)
	h.Rooms[room.ID] = room
}

// func (h *Hub) removeRoom(id string) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	if _, ok := h.Rooms[id]; ok {
// 		delete(h.Rooms, id)
// 	}
// }
