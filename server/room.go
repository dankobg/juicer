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
	id                string
	clients           map[string]*client
	clientIDS         [2]string
	gameState         *gameState
	mu                sync.Mutex
	hub               *hub
	wg                *sync.WaitGroup
	disconnectTimeout time.Duration
	firstMoveTimeout  time.Duration
	timerEvent        chan timerEvent
	players           map[string]*player
}

func (r *room) String() string {
	return r.id
}

func NewRoom(hub *hub, c1, c2 *client) (*room, error) {
	whiteID := c1.id
	blackID := c2.id

	if rand.IntN(2) == 1 {
		whiteID = c2.id
		blackID = c1.id
	}

	gs, err := NewGameState(whiteID, blackID, gameTypeStandard, gameModeBlitz)
	if err != nil {
		return nil, err
	}

	room := &room{
		id:                random.AlphaNumeric(32),
		gameState:         gs,
		clients:           make(map[string]*client),
		clientIDS:         [2]string{c1.id, c2.id},
		hub:               hub,
		wg:                &sync.WaitGroup{},
		disconnectTimeout: time.Second * 10,
		firstMoveTimeout:  time.Second * 15,
		timerEvent:        make(chan timerEvent),
		players: map[string]*player{
			whiteID: {id: whiteID, color: "w"},
			blackID: {id: blackID, color: "b"},
		},
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

func (r *room) startFirstMoveTimer(playerID string) {
	go func() {
		r.timerEvent <- timerEvent{playerID: playerID, kind: timerKindFirstMove, action: timerActionStart}
	}()
}

func (r *room) stopFirstMoveTimer(playerID string) {
	go func() {
		r.timerEvent <- timerEvent{playerID: playerID, kind: timerKindFirstMove, action: timerActionStop}
	}()
}

func (r *room) startDisconnectTimer(playerID string) {
	go func() {
		r.timerEvent <- timerEvent{playerID: playerID, kind: timerKindDisconnect, action: timerActionStart}
	}()
}

func (r *room) stopDisconnectTimer(playerID string) {
	go func() {
		r.timerEvent <- timerEvent{playerID: playerID, kind: timerKindDisconnect, action: timerActionStop}
	}()
}

func (r *room) handleTimerEvent(event timerEvent) {
	p := r.players[event.playerID]

	switch event.kind {
	case timerKindDisconnect:
		if event.action == timerActionStart {
			r.hub.log.Debug("starting disconnect timer", slog.String("client_id", p.id), slog.String("color", p.color))
			p.startDisconnectTimer(r.disconnectTimeout)
		}
		if event.action == timerActionStop {
			r.hub.log.Debug("stopping disconnect timer", slog.String("client_id", p.id), slog.String("color", p.color))
			p.stopDisconnectTimer()
		}
	case timerKindFirstMove:
		if event.action == timerActionStart {
			r.hub.log.Debug("starting first move timer", slog.String("client_id", p.id), slog.String("color", p.color))
			p.startFirstMoveTimer(r.firstMoveTimeout)
		}
		if event.action == timerActionStop {
			r.hub.log.Debug("stopping first move timer", slog.String("client_id", p.id), slog.String("color", p.color))
			p.stopFirstMoveTimer()

			if p.color == "w" {
				b := r.players[r.gameState.BlackID]
				r.hub.log.Debug("waiting for first move", slog.String("client_id", b.id), slog.String("color", b.color))
				b.startFirstMoveTimer(r.firstMoveTimeout)
			}
		}
	}
}

func (r *room) startGame(ctx context.Context) {
	r.hub.log.Debug("room starting game", slog.String("room_id", r.id), slog.String("game_id", r.gameState.GameID))

	r.wg.Add(1)
	go func() {
		defer func() {
			for _, p := range r.players {
				p.stopFirstMoveTimer()
				p.stopDisconnectTimer()
			}
			r.wg.Done()
		}()

		whitep := r.players[r.gameState.WhiteID]
		blackp := r.players[r.gameState.BlackID]

		for {
			select {
			case event := <-r.timerEvent:
				r.handleTimerEvent(event)

			case <-tick(whitep.firstMoveTimer):
				r.hub.log.Debug("white first move timer timed out")
				r.abortGame()
				return

			case <-tick(blackp.firstMoveTimer):
				r.hub.log.Debug("black first move timer timed out")
				r.abortGame()
				return

			case <-tick(whitep.disconnectTimer):
				r.hub.log.Debug("white disconnect timer timed out")
				r.abortGame()
				return

			case <-tick(blackp.disconnectTimer):
				r.hub.log.Debug("black disconnect timer timed out")
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
