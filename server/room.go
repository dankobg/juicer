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
	id           string
	clients      map[string]*client
	clientIDS    [2]string
	gameState    *gameState
	mu           sync.Mutex
	hub          *hub
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
		id:           random.AlphaNumeric(32),
		gameState:    gs,
		clients:      make(map[string]*client),
		clientIDS:    [2]string{c1.id, c2.id},
		hub:          hub,
		waitTimeout:  time.Second * 10,
		wg:           &sync.WaitGroup{},
		startWaiting: make(chan struct{}),
		stopWaiting:  make(chan struct{}),
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

func (r *room) removeClients() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.clients {
		delete(r.clients, c.id)
	}
}

func (r *room) startDisconnectTimer() {
	go func() {
		r.startWaiting <- struct{}{}
	}()
}

func (r *room) stopDisconnectTimer() {
	go func() {
		r.stopWaiting <- struct{}{}
	}()
}

func (r *room) startGame(ctx context.Context) {
	r.hub.log.Debug("room starting game", slog.String("room_id", r.id), slog.String("game_id", r.gameState.GameID))

	r.wg.Add(1)

	go func() {
		defer func() {
			if r.waitTimer != nil {
				r.waitTimer.Stop()
				r.waitTimer = nil
			}
			r.wg.Done()
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

			case <-func() <-chan time.Time {
				if r.waitTimer != nil {
					return r.waitTimer.C
				}
				return nil
			}():
				r.hub.log.Debug("room wait timer timed out, aborting game")
				r.abortGame()
				return

			case <-ctx.Done():
				r.hub.log.Debug("room context done")
				r.abortGame()
				return
			}
		}
	}()

	r.wg.Wait()
	r.hub.log.Debug("room cleaned up")
}

func (r *room) abortGame() {
	r.hub.log.Debug("abort game called")

	gameAbortedMsg := &pb.Message{
		Event: &pb.Message_GameAborted{GameAborted: &pb.GameAborted{GameId: r.gameState.GameID, RoomId: r.id, Reason: "timeout"}},
	}
	for _, c := range r.clients {
		if _, ok := r.hub.lobbyClients[c.id]; !ok {
			r.hub.addClientToLobby(c)
		}
		r.hub.broadcastClient <- &clientMessage{ClientID: c.id, Message: gameAbortedMsg}
	}

	for _, cid := range r.clientIDS {
		r.removeClient(cid)
		r.hub.removeClientsFromGame(cid)
	}

	r.gameState.MatchState = matchStateAborted
	r.hub.removeRoom(r.id)
	r.hub.broadcastHubInfo()
}
