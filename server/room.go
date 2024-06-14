package server

import (
	"math/rand/v2"
	"sync"

	"github.com/dankobg/juicer/random"
)

type Room struct {
	ID        string
	Clients   map[string]*Client
	GameState *GameState
	mu        sync.Mutex
}

func (r *Room) String() string {
	return r.ID
}

func NewRoom(c1, c2 *Client) (*Room, error) {
	roomId := random.AlphaNumeric(32)

	white := &Player{ID: c1.ID}
	black := &Player{ID: c2.ID}

	if rand.IntN(2) == 1 {
		white.ID = c2.ID
		black.ID = c1.ID
	}

	gs, err := NewGameState(white, black, GameTypeStandard, GameModeBlitz)
	if err != nil {
		return nil, err
	}

	room := &Room{
		ID:        roomId,
		GameState: gs,
		Clients:   make(map[string]*Client),
	}

	room.addClient(c1)
	room.addClient(c2)

	return room, nil
}

func (r *Room) addClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[client.ID] = client
}

func (r *Room) removeClient(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, clientID)
}

func (r *Room) StartGame() error {
	return nil
}
