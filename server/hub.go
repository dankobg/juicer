package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

type Message struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

type Hub struct {
	mu                 sync.RWMutex
	Clients            map[string]*Client
	Rooms              map[string]*Room
	ClientConnected    chan *Client
	ClientDisconnected chan *Client
	Broadcast          chan *Message
	Log                *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		mu:                 sync.RWMutex{},
		Clients:            make(map[string]*Client),
		Rooms:              make(map[string]*Room),
		ClientConnected:    make(chan *Client),
		ClientDisconnected: make(chan *Client),
		Broadcast:          make(chan *Message, 256),
		Log:                logger,
	}
}

func (h *Hub) ClientsCount() int {
	return len(h.Clients)

}

func (h *Hub) RoomsCount() int {
	return len(h.Rooms)
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Clients[client.ID] = client
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Clients, client.ID)
}

func (h *Hub) addRoom(room *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Rooms[room.ID] = room
}

func (h *Hub) removeRoom(room *Room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Rooms, room.ID)
}

func (h *Hub) HandleClientConnected(client *Client) {
	h.Log.Info("client connected", slog.String("client_id", client.ID), slog.String("room_id", client.RoomID), slog.String("remote_addr", client.Conn.RemoteAddr().String()))

	h.addClient(client)

	h.Log.Info("Hub clients info", slog.Int("clients_count", h.ClientsCount()), slog.Int("rooms_count", h.RoomsCount()))

	clientJoinedMsg := &Message{Type: "client_joined", Data: []byte(fmt.Sprintf(`{"msg":"client joined","id":"%s"}`, client.ID))}
	clientsCountMsg := &Message{Type: "clients_count", Data: []byte(fmt.Sprintf(`{"lobby":"%d","rooms":"%d"}`, h.ClientsCount(), h.RoomsCount()))}
	for _, lobbyClient := range h.Clients {
		h.Broadcast <- clientsCountMsg
		if lobbyClient.ID != client.ID {
			h.Broadcast <- clientJoinedMsg
		}
	}
}

func (h *Hub) HandleClientDisconnected(client *Client) {
	h.Log.Info("client disconnected", slog.String("client_id", client.ID), slog.String("room_id", client.RoomID), slog.String("remote_addr", client.Conn.RemoteAddr().String()))

	h.removeClient(client)
	close(client.Send)

	clientLeftMsg := &Message{Type: "client_left", Data: []byte(fmt.Sprintf(`{"msg": "client left", "id": "%s"}`, client.ID))}
	clientsCountMsg := &Message{Type: "clients_count", Data: []byte(fmt.Sprintf(`{"lobby":"%d","rooms":"%d"}`, h.ClientsCount(), h.RoomsCount()))}
	for range h.Clients {
		h.Broadcast <- clientLeftMsg
		h.Broadcast <- clientsCountMsg
	}
}

func (h *Hub) HandleBroadcast(msg *Message) {
	h.Log.Info("broadcast msg received", slog.String("type", msg.Type), slog.String("data", string(msg.Data)))

	for _, client := range h.Clients {
		select {
		case client.Send <- msg:
		default:
			h.removeClient(client)
			close(client.Send)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.ClientConnected:
			h.HandleClientConnected(client)
		case client := <-h.ClientDisconnected:
			h.HandleClientDisconnected(client)
		case msg := <-h.Broadcast:
			h.HandleBroadcast(msg)
		}
	}
}
