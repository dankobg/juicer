package server

import (
	"math/rand/v2"
	"sync"
	"time"

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
	roomId := random.AlphaNumeric(32)

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
		id:          roomId,
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
	r.startWaiting <- struct{}{}
}

func (r *room) stopDisconnectTimer() {
	r.stopWaiting <- struct{}{}
}

func (r *room) startGame() error {
	r.hub.log.Debug("starting")

	// gs.wg.Add(1)

	tick := func(tkr *time.Timer) <-chan time.Time {
		if tkr != nil {
			return tkr.C
		}
		return nil
	}

	go func() {
		defer func() {
			// gs.wg.Done()

			if r.waitTimer != nil {
				r.waitTimer.Stop()
			}
		}()

		for {
			select {
			case <-r.startWaiting:
				r.hub.log.Debug("start wait timer")
				if r.waitTimer == nil {
					r.waitTimer = time.NewTimer(r.waitTimeout)
				} else {
					r.waitTimer.Reset(r.waitTimeout)
				}

			case <-r.stopWaiting:
				r.hub.log.Debug("stop wait timer")
				if r.waitTimer != nil {
					r.waitTimer.Stop()
					r.waitTimer = nil
				}

			case <-tick(r.waitTimer):
				r.hub.log.Debug("timed out, aborting game")
				// r.abortGame()
				return
			}
		}
	}()

	// gs.wg.Wait()
	return nil
}

func (r *room) abortGame() {
	r.gameState.MatchState = matchStateAborted
	r.hub.broadcastRoom <- &roomMessage{}

	info := gameInfo{
		RoomID: r.id,
		GameID: r.gameState.GameID,
	}

	r.hub.abortGame(info)
}
