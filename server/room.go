package server

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/random"
)

type gameInfo struct {
	GameID string
	RoomID string
}

type room struct {
	id        string
	clients   map[string]*client
	gameState *gameState
	mu        sync.Mutex
	hub       *hub

	waitTimeout  time.Duration
	waitTimer    *time.Timer
	startWaiting chan struct{}
	stopWaiting  chan struct{}
	wg           *sync.WaitGroup
}

func (r *room) String() string {
	return r.id
}

func NewRoom(hub *hub, c1, c2 *client) (*room, error) {
	white := &player{ID: c1.id}
	black := &player{ID: c2.id}

	if rand.IntN(2) == 1 {
		white.ID = c2.id
		black.ID = c1.id
	}

	gs, err := NewGameState(white, black, gameTypeStandard, gameModeBlitz)
	if err != nil {
		return nil, err
	}

	room := &room{
		id:          random.AlphaNumeric(32),
		gameState:   gs,
		clients:     make(map[string]*client),
		hub:         hub,
		waitTimeout: time.Second * 10,
		wg:          &sync.WaitGroup{},
	}

	room.addClient(c1)
	room.addClient(c2)

	return room, nil
}

func (r *room) addClient(client *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.id] = client
}

func (r *room) removeClient(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, clientID)
}

func (r *room) startDisconnectTimer() {
	r.hub.log.Debug("starting room disconnect timer")
	r.startWaiting <- struct{}{}
}

func (r *room) stopDisconnectTimer() {
	r.hub.log.Debug("stopping room disconnect timer")
	r.stopWaiting <- struct{}{}
}

func (r *room) startGame(ctx context.Context, gameErr chan error) {
	r.hub.log.Debug("room starting game", slog.String("room_id", r.id), slog.String("game_id", r.gameState.GameID))

	tick := func(tkr *time.Timer) <-chan time.Time {
		if tkr != nil {
			return tkr.C
		}
		return nil
	}

	go func() {
		defer func() {
			if r.waitTimer != nil {
				r.waitTimer.Stop()
				r.waitTimer = nil
			}
		}()

		for {
			select {
			case <-r.startWaiting:
				r.hub.log.Debug("room starting wait timer")
				if r.waitTimer == nil {
					r.waitTimer = time.NewTimer(r.waitTimeout)
				} else {
					r.waitTimer.Reset(r.waitTimeout)
				}

			case <-r.stopWaiting:
				r.hub.log.Debug("room stopping wait timer")
				if r.waitTimer != nil {
					r.waitTimer.Stop()
					r.waitTimer = nil
				}

			case <-tick(r.waitTimer):
				r.hub.log.Debug("room wait timer timed out, aborting game")
				r.abortGame()
				return

			case <-ctx.Done():
				r.hub.log.Debug("room context done")
				r.abortGame()
				return

			case err := <-gameErr:
				r.hub.log.Debug("game had error", slog.Any("error", err))
				r.abortGame()
				return
			}
		}
	}()
}

func (r *room) abortGame() {
	r.hub.log.Debug("abort game called")
	r.gameState.MatchState = matchStateAborted

	gameAbortedMsg := &pb.Message{
		Event: &pb.Message_GameAborted{GameAborted: &pb.GameAborted{GameId: r.gameState.GameID, RoomId: r.id, Reason: "timeout"}},
	}
	r.hub.broadcastRoom <- &roomMessage{Message: gameAbortedMsg}
}
