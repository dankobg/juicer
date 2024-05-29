package server

type Room struct {
	ID        string
	Clients   map[string]*Client
	GameState struct{}
}

func (r *Room) String() string {
	return r.ID
}

func NewRoom(c1, c2 *Client) *Room {
	return nil

	// p1 := &engine.Player{
	// 	Name:  c1.NetPlayer.Name,
	// 	Color: engine.Color(c1.NetPlayer.Color),
	// }

	// p2 := &engine.Player{
	// 	Name:  c2.NetPlayer.Name,
	// 	Color: engine.Color(c2.NetPlayer.Color),
	// }

	// gs, err := NewGameState(p1, p2)
	// if err != nil {
	// 	fmt.Printf("NEW GS err: %+v\n", err)
	// }

	// clients := make(map[string]*Client)
	// clients[c1.ID] = c1
	// clients[c2.ID] = c2

	// return &Room{
	// 	ID:        uuid.NewString(),
	// 	GameState: gs,
	// 	Clients:   clients,
	// }
}

// **********************************************************
// **********************************************************
// **********************************************************
// **********************************************************
// **********************************************************
// **********************************************************

// package gameserver

// import (
// 	"fmt"
// 	"juicer/engine"

// 	"github.com/google/uuid"
// )

// type Room struct {
// 	ID        string
// 	GameState *GameState
// 	Clients   map[string]*Client
// }

// func (r *Room) String() string {
// 	return r.ID
// }

// func NewRoom(c1, c2 *Client) *Room {
// 	p1 := &engine.Player{
// 		Name:  c1.NetPlayer.Name,
// 		Color: engine.Color(c1.NetPlayer.Color),
// 	}

// 	p2 := &engine.Player{
// 		Name:  c2.NetPlayer.Name,
// 		Color: engine.Color(c2.NetPlayer.Color),
// 	}

// 	gs, err := NewGameState(p1, p2)
// 	if err != nil {
// 		fmt.Printf("NEW GS err: %+v\n", err)
// 	}

// 	clients := make(map[string]*Client)
// 	clients[c1.ID] = c1
// 	clients[c2.ID] = c2

// 	return &Room{
// 		ID:        uuid.NewString(),
// 		GameState: gs,
// 		Clients:   clients,
// 	}
// }

// func (r *Room) startClocks() {
// 	fmt.Printf("start clock")
// }

// func (r *Room) StartGame() {
// 	r.startClocks()
// }
