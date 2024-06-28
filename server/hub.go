package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/redis/go-redis/v9"
)

type lobbyMessage struct {
	*pb.Message
}

type clientMessage struct {
	*pb.Message
	RoomID   string
	ClientID string
}

type roomMessage struct {
	*pb.Message
	RoomID string
}

type hub struct {
	lobbyClients       map[string]*client
	clientsInGame      map[string]gameInfo
	rooms              map[string]*room
	clientConnected    chan *client
	clientDisconnected chan *client
	broadcastLobby     chan *lobbyMessage
	broadcastRoom      chan *roomMessage
	broadcastClient    chan *clientMessage
	log                *slog.Logger
	rdb                *redis.Client
	mu                 *sync.RWMutex
	wg                 *sync.WaitGroup
}

func NewHub(logger *slog.Logger, rdb *redis.Client) *hub {
	return &hub{
		lobbyClients:       make(map[string]*client),
		clientsInGame:      make(map[string]gameInfo),
		rooms:              make(map[string]*room),
		clientConnected:    make(chan *client),
		clientDisconnected: make(chan *client),
		broadcastLobby:     make(chan *lobbyMessage, 256),
		broadcastRoom:      make(chan *roomMessage, 256),
		broadcastClient:    make(chan *clientMessage, 256),
		log:                logger,
		rdb:                rdb,
		mu:                 &sync.RWMutex{},
		wg:                 &sync.WaitGroup{},
	}
}

func (h *hub) Stop(ctx context.Context) error {
	return nil
}

func (h *hub) lobbyClientsCount() int {
	return len(h.lobbyClients)
}

func (h *hub) roomsCount() int {
	return len(h.rooms)
}

func (h *hub) clientsInGameCount() int {
	return len(h.clientsInGame)
}

func (h *hub) addClientToLobby(client *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lobbyClients[client.id] = client
}

func (h *hub) removeClientFromLobby(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.lobbyClients, clientID)
}

func (h *hub) addRoom(room *room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.rooms[room.id] = room
}

func (h *hub) removeRoom(roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, roomID)
}

func (h *hub) broadcastToLobby(msg *lobbyMessage) {
	h.log.Debug("lobby broadcast", slog.String("msg", msg.String()))

	for _, client := range h.lobbyClients {
		select {
		case client.send <- msg.Message:
		default:
			h.removeClientFromLobby(client.id)
			close(client.send)
		}
	}
}

func (h *hub) broadcastToRoom(msg *roomMessage) {
	h.log.Debug("room broadcast", slog.String("room_id", msg.RoomID), slog.String("msg", msg.String()))

	room, ok := h.rooms[msg.RoomID]
	if !ok {
		h.log.Debug("room send, room_id does not exist", slog.String("room_id", msg.RoomID), slog.String("msg", msg.String()))
		return
	}

	for _, roomClient := range room.clients {
		select {
		case roomClient.send <- msg.Message:
		default:
			h.removeClientFromLobby(roomClient.id)
			close(roomClient.send)
		}
	}
}

func (h *hub) broadcastToClient(msg *clientMessage) {
	h.log.Debug("client broadcast", slog.String("client_id", msg.ClientID), slog.String("msg", msg.String()))

	if msg.RoomID != "" {
		if r, ok := h.rooms[msg.RoomID]; ok {
			if c, ok := r.clients[msg.ClientID]; ok {
				select {
				case c.send <- msg.Message:
				default:
					h.removeClientFromLobby(c.id)
					close(c.send)
				}
			}
		}
	} else {
		if c, ok := h.lobbyClients[msg.ClientID]; ok {
			select {
			case c.send <- msg.Message:
			default:
				h.removeClientFromLobby(c.id)
				close(c.send)
			}
		}
	}
}

func (h *hub) handleBroadcastToLobby(msg *lobbyMessage) {
	h.broadcastToLobby(msg)
}

func (h *hub) handleBroadcastToRoom(msg *roomMessage) {
	h.broadcastToRoom(msg)
}

func (h *hub) handleBroadcastToClient(msg *clientMessage) {
	h.broadcastToClient(msg)
}

func (h *hub) handleClientConnected(client *client) {
	info, inGame := h.clientsInGame[client.id]
	if inGame {
		client.roomID = info.RoomID
		if r, ok := h.rooms[info.RoomID]; ok {
			r.addClient(client)
			r.stopDisconnectTimer(client.id)

			for _, c := range r.clients {
				if c.id == client.id {
					continue
				}
				clientConnectedMsg := &pb.Message{
					Event: &pb.Message_ClientConnected{ClientConnected: &pb.ClientConnected{Id: client.id}},
				}
				h.broadcastClient <- &clientMessage{RoomID: r.id, ClientID: c.id, Message: clientConnectedMsg}
			}
			matchFoundMsg := &pb.Message{
				Event: &pb.Message_MatchFound{MatchFound: &pb.MatchFound{GameId: info.GameID, RoomId: info.RoomID}},
			}
			h.broadcastClient <- &clientMessage{RoomID: r.id, ClientID: client.id, Message: matchFoundMsg}
		}
	} else {
		h.addClientToLobby(client)
	}

	h.log.Debug("client connected",
		slog.String("client_id", client.id),
		slog.Bool("in_game", inGame),
		slog.String("room_id", info.RoomID),
		slog.String("game_id", info.GameID),
		slog.String("remote_addr", client.conn.RemoteAddr().String()),
		slog.Int("clients_count", h.lobbyClientsCount()),
		slog.Int("rooms_count", h.roomsCount()),
	)

	h.broadcastHubInfo()

	clientConnectedMsg := &pb.Message{
		Event: &pb.Message_ClientConnected{ClientConnected: &pb.ClientConnected{Id: client.id}},
	}
	h.broadcastLobby <- &lobbyMessage{Message: clientConnectedMsg}
}

func (h *hub) handleClientDisconnected(client *client) {
	h.removeClientFromLobby(client.id)

	info, inGame := h.clientsInGame[client.id]
	if inGame {
		client.roomID = info.RoomID
		if r, ok := h.rooms[info.RoomID]; ok {
			r.removeClient(client.id)
			r.startDisconnectTimer(client.id)

			for _, c := range r.clients {
				if c.id == client.id {
					continue
				}
				clientDisconnectedMsg := &pb.Message{
					Event: &pb.Message_ClientDisconnected{ClientDisconnected: &pb.ClientDisconnected{Id: client.id}},
				}
				h.broadcastClient <- &clientMessage{RoomID: r.id, ClientID: c.id, Message: clientDisconnectedMsg}
			}
		}
	}

	h.log.Debug("client disconnected",
		slog.String("client_id", client.id),
		slog.Bool("in_game", inGame),
		slog.String("room_id", info.RoomID),
		slog.String("game_id", info.GameID),
		slog.String("remote_addr", client.conn.RemoteAddr().String()),
		slog.Int("clients_count", h.lobbyClientsCount()),
		slog.Int("rooms_count", h.roomsCount()),
	)

	// close(client.send)

	h.broadcastHubInfo()

	clientDisconnectedMsg := &pb.Message{
		Event: &pb.Message_ClientDisconnected{ClientDisconnected: &pb.ClientDisconnected{Id: client.id}},
	}
	h.broadcastLobby <- &lobbyMessage{Message: clientDisconnectedMsg}
}

func (h *hub) processClientMessage(client *client, msg *pb.Message) error {
	cmsg := &clientMessage{
		ClientID: client.id,
		Message:  msg,
	}

	switch msg.Event.(type) {
	case *pb.Message_Echo:
		h.onEcho(cmsg)
	case *pb.Message_Chat:
		h.onChat(cmsg)
	case *pb.Message_SeekGame:
		h.onSeekGame(cmsg)
	case *pb.Message_CancelSeekGame:
		h.onCancelSeekGame(cmsg)
	case *pb.Message_AbortGame:
		h.onAbortGame(cmsg)
	case *pb.Message_OfferDraw:
		h.onOfferDraw(cmsg)
	case *pb.Message_AcceptDraw:
		h.onAcceptDraw(cmsg)
	case *pb.Message_PlayMoveUci:
		h.onPlayMoveUCI(cmsg)
	}

	return nil
}

func (h *hub) onEcho(msg *clientMessage) {
	h.broadcastClient <- msg
}

func (h *hub) onChat(msg *clientMessage) {
	// ctx := context.TODO()

	chatMsg := msg.Message.GetChat()
	if chatMsg == nil {
		h.log.Error("nil pb message", slog.String("msg", msg.String()))
		return
	}

	info, inGame := h.clientsInGame[msg.ClientID]
	if inGame {
		if r, ok := h.rooms[info.RoomID]; ok {
			chatMsg := &pb.Message{
				Event: &pb.Message_Chat{Chat: &pb.Chat{Message: chatMsg.Message}},
			}
			var otherPlayerID string
			for _, cid := range r.clientIDS {
				if cid == msg.ClientID {
					continue
				}
				otherPlayerID = cid
			}
			h.broadcastClient <- &clientMessage{ClientID: otherPlayerID, RoomID: r.id, Message: chatMsg}
		}
	}
}

func (h *hub) onSeekGame(msg *clientMessage) {
	ctx := context.TODO()

	seekGameMsg := msg.Message.GetSeekGame()
	if seekGameMsg == nil {
		h.log.Error("nil pb message", slog.String("msg", msg.String()))
		return
	}

	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZAdd(ctx, "seeking_game", redis.Z{Score: 1500, Member: msg.ClientID}).Err(); err != nil {
			h.log.Error("seeking_game add to priority queue", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		if err := p.Publish(ctx, "seeking_game:joined", msg.ClientID).Err(); err != nil {
			h.log.Error("seeking_game publish joined", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		return nil
	}); err != nil {
		h.log.Error("seeking_game pipeline", slog.String("client_id", msg.ClientID), slog.Any("error", err))
	}

	h.sendSeekingPlayersCount()

	h.log.Debug("seek game success", slog.String("client_id", msg.ClientID))
}

func (h *hub) onCancelSeekGame(msg *clientMessage) {
	ctx := context.TODO()

	cancelSeekGameMsg := msg.Message.GetCancelSeekGame()
	if cancelSeekGameMsg == nil {
		h.log.Error("nil message", slog.String("event", msg.String()))
		return
	}

	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, "seeking_game", msg.ClientID).Err(); err != nil {
			h.log.Error("cancel_seeking_game remove from priority queue", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		if err := p.Publish(ctx, "seeking_game:left", msg.ClientID).Err(); err != nil {
			h.log.Error("cancel_seeking_game publish left", slog.String("client_id", msg.ClientID), slog.Any("error", err))
		}
		return nil
	}); err != nil {
		h.log.Error("cancel_seeking_game pipeline", slog.String("client_id", msg.ClientID), slog.Any("error", err))
	}

	h.sendSeekingPlayersCount()

	h.log.Debug("cancel seek game success", slog.String("client_id", msg.ClientID))
}

func (h *hub) onAbortGame(msg *clientMessage) {
	// ctx := context.TODO()

	abortGameMsg := msg.Message.GetAbortGame()
	if abortGameMsg == nil {
		h.log.Error("nil message", slog.String("event", msg.String()))
		return
	}

	info, inGame := h.clientsInGame[msg.ClientID]
	if inGame {
		if r, ok := h.rooms[info.RoomID]; ok {
			r.finishGame <- finishGame{clientID: msg.ClientID, result: resultAborted, status: resultStatusAborted}
		}
	}
}

func (h *hub) onOfferDraw(msg *clientMessage) {
	// ctx := context.TODO()

	offerDrawMsg := msg.Message.GetOfferDraw()
	if offerDrawMsg == nil {
		h.log.Error("nil message", slog.String("event", msg.String()))
		return
	}

	info, inGame := h.clientsInGame[msg.ClientID]
	if inGame {
		if r, ok := h.rooms[info.RoomID]; ok {
			offerDrawMsg := &pb.Message{
				Event: &pb.Message_OfferDraw{OfferDraw: &pb.OfferDraw{}},
			}
			var otherPlayerID string
			for _, cid := range r.clientIDS {
				if cid == msg.ClientID {
					continue
				}
				otherPlayerID = cid
			}
			h.broadcastClient <- &clientMessage{ClientID: otherPlayerID, RoomID: r.id, Message: offerDrawMsg}
		}
	}
}

func (h *hub) onAcceptDraw(msg *clientMessage) {
	// ctx := context.TODO()

	acceptDrawMsg := msg.Message.GetAcceptDraw()
	if acceptDrawMsg == nil {
		h.log.Error("nil message", slog.String("event", msg.String()))
		return
	}

	info, inGame := h.clientsInGame[msg.ClientID]
	if inGame {
		if r, ok := h.rooms[info.RoomID]; ok {
			r.finishGame <- finishGame{clientID: msg.ClientID, result: resultDraw, status: resultStatusDrawAgreed}
		}
	}
}

func (h *hub) onPlayMoveUCI(msg *clientMessage) {
	// ctx := context.TODO()

	playMoveUciMsg := msg.Message.GetPlayMoveUci()
	if playMoveUciMsg == nil {
		h.log.Error("nil message", slog.String("event", msg.String()))
		return
	}

}

func (h *hub) broadcastHubInfo() {
	hubInfoMsg := &pb.Message{
		Event: &pb.Message_HubInfo{HubInfo: &pb.HubInfo{Lobby: int32(h.lobbyClientsCount()), Rooms: int32(h.roomsCount()), Playing: int32(h.clientsInGameCount())}},
	}
	h.broadcastLobby <- &lobbyMessage{Message: hubInfoMsg}
}

func (h *hub) Run(ctx context.Context) {
	go func() {
		h.runMatchmaking(ctx)
	}()

	for {
		select {
		case client := <-h.clientConnected:
			h.handleClientConnected(client)
		case client := <-h.clientDisconnected:
			h.handleClientDisconnected(client)
		case msg := <-h.broadcastLobby:
			h.handleBroadcastToLobby(msg)
		case msg := <-h.broadcastRoom:
			h.handleBroadcastToRoom(msg)
		case msg := <-h.broadcastClient:
			h.handleBroadcastToClient(msg)
		case <-ctx.Done():
			h.log.Debug("hub run ctx Done")
		}
	}
}

func (h *hub) tryMatchPlayers(ctx context.Context) {
	res, err := h.rdb.ZRangeByScore(ctx, "seeking_game", &redis.ZRangeBy{Min: "1400", Max: "1600", Count: 2}).Result()
	if err != nil {
		h.log.Error("failed to fetch players from priority queue", slog.Any("error", err))
		return
	}

	if len(res) != 2 {
		// h.Log.Debug("not a pair of 2 players", slog.Int("fetched", len(res))) // too noisy even for dbg?
		return
	}

	h.log.Debug("fetched players list", slog.Any("client_ids", res))

	h.mu.RLock()
	c1, ok1 := h.lobbyClients[res[0]]
	c2, ok2 := h.lobbyClients[res[1]]
	h.mu.RUnlock()

	if !ok1 || !ok2 {
		h.log.Error("client not found in lobby")
		return
	}

	room, err := NewRoom(h, c1, c2)
	if err != nil {
		h.log.Error("failed to create room", slog.Any("error", err))
		return
	}
	c1.roomID = room.id
	c2.roomID = room.id
	h.addRoom(room)

	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, "seeking_game", res[0], res[1]).Err(); err != nil {
			h.log.Error("match found remove from priority queue", slog.String("client1_id", c1.id), slog.String("client2_id", c2.id), slog.String("room_id", room.id), slog.Any("error", err))
			return err
		}

		if err := p.Publish(ctx, "game_found", room.id).Err(); err != nil {
			h.log.Error("match found publish pair success", slog.String("client1_id", c1.id), slog.String("client2_id", c2.id), slog.String("game_id", room.gameState.GameID), slog.String("room_id", room.id), slog.Any("error", err))
			return err
		}

		return nil
	}); err != nil {
		h.removeRoom(room.id)
		h.log.Error("match found pipeline", slog.String("room_id", room.id), slog.Any("error", err))
	}

	gameInfo := gameInfo{
		RoomID: room.id,
		GameID: room.gameState.GameID,
	}

	h.addClientsToGame(gameInfo, c1.id, c2.id)
	h.removeClientFromLobby(c1.id)
	h.removeClientFromLobby(c2.id)

	go room.startGame(ctx)

	matchFoundMsg := &pb.Message{
		Event: &pb.Message_MatchFound{MatchFound: &pb.MatchFound{GameId: room.gameState.GameID, RoomId: room.id}},
	}
	h.broadcastRoom <- &roomMessage{RoomID: room.id, Message: matchFoundMsg}

	h.broadcastHubInfo()

	h.log.Debug("match found success", slog.String("room_id", room.id), slog.String("game_id", room.gameState.GameID))
}

func (h *hub) addClientsToGame(info gameInfo, c1, c2 string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clientsInGame[c1] = info
	h.clientsInGame[c2] = info
}

func (h *hub) removeClientsFromGame(clientIDs ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, id := range clientIDs {
		delete(h.clientsInGame, id)
	}
}

func (h *hub) sendSeekingPlayersCount() {
	count, err := h.seekingGamePlayersCount()
	if err != nil {
		h.log.Debug("seeking players zcard", slog.Any("error", err))
		return
	}

	seekingCountMsg := &pb.Message{
		Event: &pb.Message_SeekingCount{SeekingCount: &pb.SeekingCount{Count: int32(count)}},
	}
	h.broadcastLobby <- &lobbyMessage{Message: seekingCountMsg}
}

func (h *hub) runMatchmaking(ctx context.Context) {
	h.log.Debug("matchmaking started")

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			h.tryMatchPlayers(ctx)
		case <-ctx.Done():
			h.log.Debug("matchmaking context done")
		}
	}
}

func (h *hub) seekingGamePlayersCount() (int64, error) {
	ctx := context.TODO()
	return h.rdb.ZCard(ctx, "seeking_game").Result()
}
