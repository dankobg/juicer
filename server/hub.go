package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	Type     string          `json:"t"`
	ClientID string          `json:"c"`
	Data     json.RawMessage `json:"d"`
}

type Hub struct {
	mu                 sync.RWMutex
	Clients            map[string]*Client
	Rooms              map[string]*Room
	ClientConnected    chan *Client
	ClientDisconnected chan *Client
	Broadcast          chan *Message
	BroadcastRoom      chan *Message
	Log                *slog.Logger
	Rdb                *redis.Client
}

func NewHub(logger *slog.Logger, rdb *redis.Client) *Hub {
	return &Hub{
		mu:                 sync.RWMutex{},
		Clients:            make(map[string]*Client),
		Rooms:              make(map[string]*Room),
		ClientConnected:    make(chan *Client),
		ClientDisconnected: make(chan *Client),
		Broadcast:          make(chan *Message, 256),
		BroadcastRoom:      make(chan *Message, 256),
		Log:                logger,
		Rdb:                rdb,
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
	h.addClient(client)

	h.Log.Debug("client connected",
		slog.String("client_id", client.ID),
		slog.String("room_id", client.RoomID),
		slog.String("remote_addr", client.Conn.RemoteAddr().String()),
		slog.Int("clients_count", h.ClientsCount()),
		slog.Int("rooms_count", h.RoomsCount()),
	)

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
	h.removeClient(client)

	h.Log.Debug("client disconnected",
		slog.String("client_id", client.ID),
		slog.String("room_id", client.RoomID),
		slog.String("remote_addr", client.Conn.RemoteAddr().String()),
		slog.Int("clients_count", h.ClientsCount()),
		slog.Int("rooms_count", h.RoomsCount()),
	)

	close(client.Send)

	clientLeftMsg := &Message{Type: "client_left", Data: []byte(fmt.Sprintf(`{"msg": "client left", "id": "%s"}`, client.ID))}
	clientsCountMsg := &Message{Type: "clients_count", Data: []byte(fmt.Sprintf(`{"lobby":"%d","rooms":"%d"}`, h.ClientsCount(), h.RoomsCount()))}
	for range h.Clients {
		h.Broadcast <- clientLeftMsg
		h.Broadcast <- clientsCountMsg
	}
}

func (h *Hub) ProcessMessage(msg *Message) error {
	switch msg.Type {
	case "seek_game":
		h.onClientSeekGame(msg)
	case "cancel_seek_game":
		h.onClientCancelSeekGame(msg)
	}

	return nil
}

func (h *Hub) HandleBroadcast(msg *Message) {
	h.Log.Debug("broadcast recv", slog.String("type", msg.Type), slog.String("data", string(msg.Data)))

	for _, client := range h.Clients {
		select {
		case client.Send <- msg:
		default:
			h.removeClient(client)
			close(client.Send)
		}
	}
}

func (h *Hub) HandleBroadcastToRoom(roomID string, msg *Message) {
	h.Log.Debug("broadcast room recv", slog.String("room_id", roomID), slog.String("type", msg.Type), slog.String("data", string(msg.Data)))

	room, ok := h.Rooms[roomID]
	if !ok {
		h.Log.Debug("broadcast room missing", slog.String("room_id", roomID), slog.String("type", msg.Type), slog.String("data", string(msg.Data)))
		return
	}

	for _, roomClient := range room.Clients {
		select {
		case roomClient.Send <- msg:
		default:
			h.removeClient(roomClient)
			close(roomClient.Send)
		}
	}
}

func (h *Hub) onClientSeekGame(msg *Message) {
	ctx := context.TODO()

	type clientSeekingGame struct {
		GameMode string `json:"game_mode"`
	}
	data := clientSeekingGame{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.Log.Error("seeking_game unmarshal msg", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		return
	}

	if _, err := h.Rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZAdd(ctx, "seeking_game", redis.Z{Score: 1500, Member: msg.ClientID}).Err(); err != nil {
			h.Log.Error("seeking_game add to priority queue", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		if err := p.Publish(ctx, "seeking_game:joined", msg.ClientID).Err(); err != nil {
			h.Log.Error("seeking_game publish joined", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		return nil
	}); err != nil {
		h.Log.Error("seeking_game pipeline", slog.String("client_id", msg.ClientID), slog.Any("error", err))
	}

	h.Log.Debug("seek game success", slog.String("client_id", msg.ClientID))
}

func (h *Hub) onClientCancelSeekGame(msg *Message) {
	ctx := context.TODO()

	type clientCancelSeekingGame struct {
		ClientID string `json:"client_id"`
	}
	data := clientCancelSeekingGame{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.Log.Error("cancel_seeking_game unmarshal msg", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		return
	}

	if _, err := h.Rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, "seeking_game", msg.ClientID).Err(); err != nil {
			h.Log.Error("cancel_seeking_game remove from priority queue", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		if err := p.Publish(ctx, "seeking_game:left", msg.ClientID).Err(); err != nil {
			h.Log.Error("cancel_seeking_game publish left", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		return nil
	}); err != nil {
		h.Log.Error("cancel_seeking_game pipeline", slog.String("client_id", msg.ClientID), slog.Any("error", err))
	}

	h.Log.Debug("cancel seek game success", slog.String("client_id", msg.ClientID))
}

func (h *Hub) Run(ctx context.Context) {
	go func() {
		h.RunMatchmaking(ctx)
	}()

	for {
		select {
		case client := <-h.ClientConnected:
			h.HandleClientConnected(client)
		case client := <-h.ClientDisconnected:
			h.HandleClientDisconnected(client)
		case msg := <-h.Broadcast:
			h.HandleBroadcast(msg)
		// case msg := <-h.BroadcastRoom:
		// 	h.HandleBroadcastToRoom("???????????????????????????",msg)
		case <-ctx.Done():
			h.Log.Debug("hub run ctx Done")
		}
	}
}

func (h *Hub) TryMatchPlayers(ctx context.Context) {
	h.Log.Debug("trying to match players for a game")

	res, err := h.Rdb.ZRangeByScore(ctx, "seeking_game", &redis.ZRangeBy{Min: "1400", Max: "1600", Count: 2}).Result()
	if err != nil {
		h.Log.Error("failed to fetch players from priority queue", slog.Any("error", err))
		return
	}

	h.Log.Debug("fetched players list", slog.Any("client_ids", res))

	if len(res) != 2 {
		h.Log.Debug("not a pair of 2 players", slog.Int("fetched", len(res)))
		return
	}

	h.mu.RLock()
	c1, ok1 := h.Clients[res[0]]
	c2, ok2 := h.Clients[res[1]]
	h.mu.RUnlock()

	if !ok1 || !ok2 {
		h.Log.Error("client not found in lobby")
		return
	}

	room, err := NewRoom(c1, c2)
	if err != nil {
		h.Log.Error("failed to create room", slog.Any("error", err))
		return
	}
	c1.RoomID = room.ID
	c2.RoomID = room.ID
	h.addRoom(room)

	if _, err := h.Rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, "seeking_game", res[0], res[1]).Err(); err != nil {
			h.Log.Error("match found remove from priority queue", slog.String("client1_id", c1.ID), slog.String("client2_id", c2.ID), slog.String("room_id", room.ID), slog.Any("error", err))
			return err
		}
		if err := room.StartGame(); err != nil {
			h.Log.Error("match found startgame", slog.String("client1_id", c1.ID), slog.String("client2_id", c2.ID), slog.String("room_id", room.ID), slog.Any("error", err))
			return err
		}
		if err := p.Publish(ctx, "game_found", room.ID).Err(); err != nil {
			h.Log.Error("match found publish pair success", slog.String("client1_id", c1.ID), slog.String("client2_id", c2.ID), slog.String("room_id", room.ID), slog.Any("error", err))
			return err
		}
		return nil
	}); err != nil {
		h.removeRoom(room)
		h.Log.Error("match found pipeline", slog.String("room_id", room.ID), slog.Any("error", err))
	}

	h.Log.Debug("match found success", slog.String("room_id", room.ID))
}

func (h *Hub) RunMatchmaking(ctx context.Context) {
	h.Log.Debug("matchmaking started")

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			h.TryMatchPlayers(ctx)
		case <-ctx.Done():
			fmt.Println("matchmaking context done")
		}
	}
}
