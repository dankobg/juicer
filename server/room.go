package server

import (
	"hash/maphash"
	"math/rand"

	"github.com/dankobg/juicer/random"
)

type Room struct {
	ID        string
	Clients   map[string]*Client
	GameState *GameState
}

func (r *Room) String() string {
	return r.ID
}

func NewRoom(c1, c2 *Client) (*Room, error) {
	rnd := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))

	roomId := random.AlphaNumeric(32)

	white := &Player{ID: c1.ID}
	black := &Player{ID: c2.ID}

	if rnd.Intn(2) == 1 {
		white.ID = c2.ID
		black.ID = c1.ID
	}

	gs, err := NewGameState(white, black, GameTypeStandard, GameModeBlitz)
	if err != nil {
		return nil, err
	}

	clients := make(map[string]*Client)
	clients[c1.ID] = c1
	clients[c2.ID] = c2

	room := &Room{
		ID:        roomId,
		Clients:   clients,
		GameState: gs,
	}

	return room, nil
}
