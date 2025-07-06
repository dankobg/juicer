package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/opt"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/dankobg/juicer/store"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Channel is similar to a "room" or a "realm" for communication
// e.g. `lobby`, `game.{game_id}`, `gametv.{game_id}`, `gametv.{game_id}.chat` etc.
type Channel string

func (ch Channel) String() string {
	return string(ch)
}

// ClientMessage is msg sent to a client across all channels unless `channel` is specified
type ClientMessage struct {
	clientID uuid.UUID
	channel  Channel
	msg      []byte
}

// ChannelMessage is msg sent to everyone in that channel
type ChannelMessage struct {
	channel Channel
	msg     []byte
}

const (
	matchmakingInterval = time.Second * 10
)

// mappings from proto enums to db id (for now these won't change, even if so, very rarely)
type mappings struct {
	variants       map[pb.Variant]uuid.UUID
	timeKinds      map[pb.GameTimeKind]uuid.UUID
	timeCategories map[pb.GameTimeCategory]uuid.UUID
	results        map[pb.GameResult]uuid.UUID
	resultStatuses map[pb.GameResultStatus]uuid.UUID
	states         map[pb.GameState]uuid.UUID
}

func initMappings() mappings {
	return mappings{
		variants:       make(map[pb.Variant]uuid.UUID),
		timeKinds:      make(map[pb.GameTimeKind]uuid.UUID),
		timeCategories: make(map[pb.GameTimeCategory]uuid.UUID),
		results:        make(map[pb.GameResult]uuid.UUID),
		resultStatuses: make(map[pb.GameResultStatus]uuid.UUID),
		states:         make(map[pb.GameState]uuid.UUID),
	}
}

var pubsubTopics = []string{
	"ipc.*",    // internal
	"wsc.*",    // websocket conn
	"lobby",    // lobby
	"game.*",   // game related
	"gametv.*", // game spectators
	"client.*", // specific client
	"conn.*",   // specific conn
}

var QuickGames = []api.QuickGame{
	{Name: "Hyperbullet", ClockSecs: 30, IncrementSecs: 0},
	{Name: "Bullet", ClockSecs: 60, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 180, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 180, IncrementSecs: 1},
	{Name: "Blitz", ClockSecs: 300, IncrementSecs: 0},
	{Name: "Blitz", ClockSecs: 300, IncrementSecs: 2},
	{Name: "Rapid", ClockSecs: 600, IncrementSecs: 0},
	{Name: "Rapid", ClockSecs: 600, IncrementSecs: 5},
	{Name: "Rapid", ClockSecs: 900, IncrementSecs: 0},
	{Name: "Rapid", ClockSecs: 900, IncrementSecs: 5},
	{Name: "Classical", ClockSecs: 1800, IncrementSecs: 0},
	{Name: "Classical", ClockSecs: 2700, IncrementSecs: 10},
}

type UserInfoFetcher interface {
	FetchUserInfo(ctx context.Context, userID uuid.UUID) (UserInfo, error)
}

type UserInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

type Hub struct {
	ClientConnected     chan *client
	ClientDisconnected  chan *client
	clients             map[*client][]Channel
	clientsByUserID     map[uuid.UUID]*client
	games               map[uuid.UUID]*gameState
	channels            map[Channel]map[*client]struct{}
	topics              []string
	subs                map[string]*redis.PubSub
	subMessages         map[string]<-chan *redis.Message
	broadcastClient     chan ClientMessage
	broadcastChannel    chan ChannelMessage
	store               store.Store
	rdb                 *redis.Client
	log                 *slog.Logger
	wg                  sync.WaitGroup
	mu                  sync.RWMutex
	mappings            mappings
	categoryThressholds []categoryThreshold
	userInfoFetcher     UserInfoFetcher
}

func NewHub(store store.Store, rdb *redis.Client, sl *slog.Logger) *Hub {
	ctx := context.Background()
	hub := &Hub{
		ClientConnected:     make(chan *client),
		ClientDisconnected:  make(chan *client),
		clients:             make(map[*client][]Channel),
		clientsByUserID:     make(map[uuid.UUID]*client),
		games:               make(map[uuid.UUID]*gameState),
		topics:              pubsubTopics,
		channels:            make(map[Channel]map[*client]struct{}),
		subs:                make(map[string]*redis.PubSub),
		subMessages:         make(map[string]<-chan *redis.Message),
		broadcastClient:     make(chan ClientMessage, 500),
		broadcastChannel:    make(chan ChannelMessage, 500),
		store:               store,
		rdb:                 rdb,
		log:                 sl,
		wg:                  sync.WaitGroup{},
		mu:                  sync.RWMutex{},
		mappings:            initMappings(),
		categoryThressholds: make([]categoryThreshold, 0),
	}
	for _, topic := range hub.topics {
		pubsub := rdb.PSubscribe(ctx, topic)
		hub.subs[topic] = pubsub
		hub.subMessages[topic] = pubsub.Channel()
	}
	return hub
}

func (h *Hub) SetUserInfoFetcher(fetcher UserInfoFetcher) {
	h.userInfoFetcher = fetcher
}

func (h *Hub) fetchGameLookupTables(ctx context.Context) error {
	gameVariants, e1 := h.store.GameVariant().List(ctx)
	gameTimeKinds, e2 := h.store.GameTimeKind().List(ctx)
	gameTimeCategories, e3 := h.store.GameTimeCategory().List(ctx)
	gameResults, e4 := h.store.GameResult().List(ctx)
	gameResultStatuses, e5 := h.store.GameResultStatus().List(ctx)
	gameStates, e6 := h.store.GameState().List(ctx)
	if err := errors.Join(e1, e2, e3, e4, e5, e6); err != nil {
		return err
	}

	// for now i put this here also...
	for _, v := range gameTimeCategories {
		var limit time.Duration = math.MaxUint32
		if v.UpperTimeLimitSecs != nil {
			limit = time.Second * time.Duration(*v.UpperTimeLimitSecs)
		}
		switch v.Name {
		case "hyperbullet":
			h.categoryThressholds = append(h.categoryThressholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET, upperLimit: limit})
		case "bullet":
			h.categoryThressholds = append(h.categoryThressholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET, upperLimit: limit})
		case "blitz":
			h.categoryThressholds = append(h.categoryThressholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ, upperLimit: limit})
		case "rapid":
			h.categoryThressholds = append(h.categoryThressholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID, upperLimit: limit})
		case "classical":
			h.categoryThressholds = append(h.categoryThressholds, categoryThreshold{timeCategory: pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL, upperLimit: limit})
		}
	}

	for _, v := range gameVariants {
		switch v.Name {
		case "standard":
			h.mappings.variants[pb.Variant_VARIANT_STANDARD] = v.ID
		case "atomic":
			h.mappings.variants[pb.Variant_VARIANT_ATOMIC] = v.ID
		case "crazyhouse":
			h.mappings.variants[pb.Variant_VARIANT_CRAZYHOUSE] = v.ID
		case "chess960":
			h.mappings.variants[pb.Variant_VARIANT_CHESS960] = v.ID
		case "king-of-the-hill":
			h.mappings.variants[pb.Variant_VARIANT_KING_OF_THE_HILL] = v.ID
		case "three-check":
			h.mappings.variants[pb.Variant_VARIANT_THREE_CHECK] = v.ID
		case "horde":
			h.mappings.variants[pb.Variant_VARIANT_HORDE] = v.ID
		case "racing-kings":
			h.mappings.variants[pb.Variant_VARIANT_RACING_KINGS] = v.ID
		}
	}

	for _, v := range gameTimeKinds {
		switch v.Name {
		case "realtime":
			h.mappings.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_REALTIME] = v.ID
		case "correspondance":
			h.mappings.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_CORRESPONDANCE] = v.ID
		case "unlimited":
			h.mappings.timeKinds[pb.GameTimeKind_GAME_TIME_KIND_UNLIMITED] = v.ID
		}
	}

	for _, v := range gameTimeCategories {
		switch v.Name {
		case "hyperbullet":
			h.mappings.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_HYPERBULLET] = v.ID
		case "bullet":
			h.mappings.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_BULLET] = v.ID
		case "blitz":
			h.mappings.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_BLITZ] = v.ID
		case "rapid":
			h.mappings.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_RAPID] = v.ID
		case "classical":
			h.mappings.timeCategories[pb.GameTimeCategory_GAME_TIME_CATEGORY_CLASSICAL] = v.ID
		}
	}

	for _, v := range gameResults {
		switch v.Name {
		case "white-won":
			h.mappings.results[pb.GameResult_GAME_RESULT_WHITE_WON] = v.ID
		case "black-won":
			h.mappings.results[pb.GameResult_GAME_RESULT_BLACK_WON] = v.ID
		case "draw":
			h.mappings.results[pb.GameResult_GAME_RESULT_DRAW] = v.ID
		case "interrupted":
			h.mappings.results[pb.GameResult_GAME_RESULT_INTERRUPTED] = v.ID
		}
	}

	for _, v := range gameResultStatuses {
		switch v.Name {
		case "checkmate":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE] = v.ID
		case "insufficient-material":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_INSUFFICIENT_MATERIAL] = v.ID
		case "threefold-repetition":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_THREEFOLD_REPETITION] = v.ID
		case "fivefold-repetition":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FIVEFOLD_REPETITION] = v.ID
		case "fifty-move-rule":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FIFTY_MOVE_RULE] = v.ID
		case "seventyfive-move-rule":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_SEVENTYFIVE_MOVE_RULE] = v.ID
		case "stalemate":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE] = v.ID
		case "resignation":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION] = v.ID
		case "draw-agreed":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_DRAW_AGREED] = v.ID
		case "flagged":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED] = v.ID
		case "adjudication":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_ADJUDICATION] = v.ID
		case "timed-out":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT] = v.ID
		case "aborted":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED] = v.ID
		case "interrupted":
			h.mappings.resultStatuses[pb.GameResultStatus_GAME_RESULT_STATUS_INTERRUPTED] = v.ID
		}
	}

	for _, v := range gameStates {
		switch v.Name {
		case "idle":
			h.mappings.states[pb.GameState_GAME_STATE_IDLE] = v.ID
		case "waiting-start":
			h.mappings.states[pb.GameState_GAME_STATE_WAITING_START] = v.ID
		case "in-progress":
			h.mappings.states[pb.GameState_GAME_STATE_IN_PROGRESS] = v.ID
		case "finished":
			h.mappings.states[pb.GameState_GAME_STATE_FINISHED] = v.ID
		case "interrupted":
			h.mappings.states[pb.GameState_GAME_STATE_INTERRUPTED] = v.ID
		}
	}

	return nil
}

// Run starts the pubsub and machmaking, as well as broadcast events
func (h *Hub) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Hub recovered", r)
		}
	}()

	if err := h.fetchGameLookupTables(ctx); err != nil {
		h.log.Error("failed to fetch game lookup table mappings", slog.Any("error", err))
		return nil
	}
	h.log.Info("hub is running")
	go h.startPubSub(ctx)
	go h.startMatchmaking(ctx)
loop:
	for {
		select {
		case c := <-h.ClientConnected:
			h.onClientConnected(c)
		case c := <-h.ClientDisconnected:
			h.onClientDisconnected(c)
		case m, ok := <-h.broadcastChannel:
			if !ok {
				continue
			}
			h.broadcastToChannel(m.channel, m.msg)
		case m, ok := <-h.broadcastClient:
			if !ok {
				continue
			}
			h.broadcastToClient(m.clientID, "", m.msg)
		case <-ctx.Done():
			break loop
		}
	}
	return nil
}

// Stop stops the hub
func (h *Hub) Stop(ctx context.Context) error {
	return nil
}

func (h *Hub) logPubsubMsg(msg *redis.Message) {
	h.log.Debug("pubsub recv", slog.String("channel", msg.Channel), slog.String("pattern", msg.Pattern), slog.String("payload", msg.Payload))
}

// startPubSub starts the redis pubsub on main topics
func (h *Hub) startPubSub(ctx context.Context) {
	h.log.Info("pubsub started")
	for {
		select {
		case msg := <-h.subMessages["ipc"]:
			h.logPubsubMsg(msg)
		case msg := <-h.subMessages["wsc.*"]:
			h.logPubsubMsg(msg)
			h.handleClientMessage(msg)
		case msg := <-h.subMessages["lobby"]:
			h.logPubsubMsg(msg)
			h.broadcastChannel <- ChannelMessage{channel: "lobby", msg: []byte(msg.Payload)}
		case msg := <-h.subMessages["game.*"]:
			h.logPubsubMsg(msg)
			gameID, err := extractGameMessageInfo(msg.Channel)
			if err != nil {
				h.log.Error("pubsub extractGameMessageInfo", slog.String("channel", msg.Channel), slog.String("pattern", msg.Pattern), slog.String("payload", msg.Payload), slog.Any("error", err))
				continue
			}
			h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gameID), msg: []byte(msg.Payload)}
		case msg := <-h.subMessages["gametv.*"]:
			h.logPubsubMsg(msg)
			gameID, err := extractGameMessageInfo(msg.Pattern)
			if err != nil {
				h.log.Error("pubsub extractGameMessageInfo", slog.String("channel", msg.Channel), slog.String("pattern", msg.Pattern), slog.String("payload", msg.Payload), slog.Any("error", err))
				continue
			}
			h.broadcastChannel <- ChannelMessage{channel: Channel("gametv." + gameID), msg: []byte(msg.Payload)}
		case msg := <-h.subMessages["client.*"]:
			h.logPubsubMsg(msg)
			clientID, channel, err := extractClientSubtopics(msg.Channel)
			if err != nil {
				h.log.Error("pubsub extractClientSubtopics", slog.String("channel", msg.Channel), slog.String("pattern", msg.Pattern), slog.String("payload", msg.Payload), slog.Any("error", err))
				continue
			}
			h.broadcastClient <- ClientMessage{clientID: clientID, channel: Channel(channel), msg: []byte(msg.Payload)}
		case <-ctx.Done():
			h.log.Debug("pubsub ctx done")
			return
		}
	}
}

// startMatchmaking runs the matchmaking which periodically tries to pair players
func (h *Hub) startMatchmaking(ctx context.Context) {
	h.log.Info("matchmaking started")
	ticker := time.NewTicker(matchmakingInterval)
loop:
	for {
		select {
		case <-ticker.C:
			h.tryMatchPlayers(ctx)
		case <-ctx.Done():
			break loop
		}
	}
}

// handleClientMessage handles the republished msg from browser websocket
func (h *Hub) handleClientMessage(m *redis.Message) {
	clientID, authState, err := extractWSCSubtopics(m.Channel)
	if err != nil {
		h.log.Error("handleClientMessage extractWSCSubtopics", slog.String("channel", m.Channel), slog.String("pattern", m.Pattern), slog.String("payload", m.Payload), slog.Any("error", err))
		return
	}
	msg := &pb.Message{}
	if err := protojson.Unmarshal([]byte(m.Payload), msg); err != nil {
		h.log.Error("handleClientMessage protojson unmarshal Message", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
	}
	switch msg.GetEvent().(type) {
	case *pb.Message_SeekGame:
		h.handleSeekGameMsg(msg.GetSeekGame(), clientID, authState)
	case *pb.Message_CancelSeekGame:
		h.handleCancelSeekGameMsg(msg.GetCancelSeekGame(), clientID, authState)
	case *pb.Message_GameAbort:
		h.handleGameAbortMsg(msg.GetGameAbort(), clientID, authState)
	case *pb.Message_GameOfferDraw:
		h.handleGameOfferDrawMsg(msg.GetGameOfferDraw(), clientID, authState)
	case *pb.Message_GameResign:
		h.handleGameResignMsg(msg.GetGameResign(), clientID, authState)
	case *pb.Message_GameDeclineDraw:
		h.handleGameDeclineDraw(msg.GetGameDeclineDraw(), clientID, authState)
	case *pb.Message_GameAcceptDraw:
		h.handleAcceptDrawMsg(msg.GetGameAcceptDraw(), clientID, authState)
	case *pb.Message_GameChat:
		h.handleGameChatMsg(msg.GetGameChat(), clientID, authState)
	case *pb.Message_PlayMoveUci:
		h.handlePlayMoveUCIMsg(msg.GetPlayMoveUci(), clientID, authState)
	}
}

func (h *Hub) handleSeekGameMsg(msg *pb.SeekGame, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil SeekGame msg")
		return
	}
	ctx := context.Background()
	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := fmt.Sprintf("seek_game.%d.%d_%d", authState, msg.GetTimeControl().Clock.Seconds, msg.GetTimeControl().Increment.Seconds)
		if err := p.ZAdd(ctx, key, redis.Z{Member: clientID, Score: float64(time.Now().UnixNano())}).Err(); err != nil {
			h.log.Error("seek_game add to queue", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		if err := p.HSet(ctx, "clients_seeking", clientID, key).Err(); err != nil {
			h.log.Error("seek_game add seek key for client", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		if err := p.Publish(ctx, key, clientID).Err(); err != nil {
			h.log.Error("seek_game publish joined queue", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		return nil
	}); err != nil {
		h.log.Error("seek_game pipeline", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
	}
	h.broadcastHubInfoToClient(uuid.MustParse(clientID))
}

func (h *Hub) handleCancelSeekGameMsg(msg *pb.CancelSeekGame, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil CancelSeekGame msg")
		return
	}
	ctx := context.Background()
	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := fmt.Sprintf("seek_game.%d.%d_%d", authState, msg.GetTimeControl().Clock.Seconds, msg.GetTimeControl().Increment.Seconds)
		if err := p.ZRem(ctx, key, clientID).Err(); err != nil {
			h.log.Error("cancel_seek_game remove from guest queue", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		if err := p.HDel(ctx, "clients_seeking", clientID).Err(); err != nil {
			h.log.Error("cancel_seek_game remove seek key for client", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		if err := p.Publish(ctx, key, clientID).Err(); err != nil {
			h.log.Error("cancel_seek_game publish cancel seek game", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
			return err
		}
		return nil
	}); err != nil {
		h.log.Error("cancel_seek_game pipeline", slog.String("client_id", clientID), slog.String("auth_state", authState.String()), slog.Any("error", err))
	}
}

func (h *Hub) handleGameAbortMsg(msg *pb.GameAbort, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameAbort msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameAbort: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("GameAbort: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("GameAbort: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("GameAbort: can't abort a finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	var result pb.GameResult
	var resultStatus pb.GameResultStatus
	var state pb.GameState
	ownColor := pb.Color_COLOR_WHITE
	if clientID == gs.blackID.String() {
		ownColor = pb.Color_COLOR_BLACK
	}
	if ownColor == pb.Color_COLOR_WHITE {
		if gs.chess.Position.Ply <= 1 {
			result = pb.GameResult_GAME_RESULT_INTERRUPTED
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			result = pb.GameResult_GAME_RESULT_BLACK_WON
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_FINISHED
		}
	}
	if ownColor == pb.Color_COLOR_BLACK {
		if gs.chess.Position.Ply <= 2 {
			result = pb.GameResult_GAME_RESULT_INTERRUPTED
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			result = pb.GameResult_GAME_RESULT_WHITE_WON
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_FINISHED
		}
	}
	updateGameDto := dto.GameChangeset{
		ResultID:       opt.New(h.mappings.results[result]),
		ResultStatusID: opt.New(h.mappings.resultStatuses[resultStatus]),
		StateID:        opt.New(h.mappings.states[state]),
		EndTime:        opt.New(time.Now()),
	}
	if _, err := h.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
		h.log.Error("GameAbort, failed to persist game", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().RemoveActiveGame(context.Background(), gs.gameID); err != nil {
		h.log.Error("GameAbort, failed to remove active game presence", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.whiteID); err != nil {
		h.log.Error("GameAbort, failed to remove active game presence 1", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.blackID); err != nil {
		h.log.Error("GameAbort, failed to remove active game presence 2", slog.Any("error", err))
		return
	}
	h.removeGame(gs.gameID)
	gameFinishedMsg := &pb.Message{
		Event: &pb.Message_GameFinished{GameFinished: &pb.GameFinished{GameResult: result, GameResultStatus: resultStatus, GameState: state}},
	}
	bb, err := protojson.Marshal(gameFinishedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameFinished", slog.Any("error", err))
		return
	}
	h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
}

func (h *Hub) handleGameOfferDrawMsg(msg *pb.GameOfferDraw, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameOfferDraw msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameOfferDraw: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("GameOfferDraw: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("GameOfferDraw: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("GameOfferDraw: can't offer draw in finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	otherID := gs.blackID
	if clientID == gs.blackID.String() {
		otherID = gs.whiteID
	}
	gameOfferDrawMsg := &pb.Message{
		Event: &pb.Message_GameOfferDraw{GameOfferDraw: msg},
	}
	bb, err := protojson.Marshal(gameOfferDrawMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameOfferDraw", slog.Any("error", err))
		return
	}
	h.broadcastClient <- ClientMessage{clientID: otherID, msg: bb}
}

func (h *Hub) handleGameResignMsg(msg *pb.GameResign, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameResign msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameResign: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("GameResign: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("GameResign: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("GameResign: can't resign a finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	var result pb.GameResult
	var resultStatus pb.GameResultStatus
	var state pb.GameState
	ownColor := pb.Color_COLOR_WHITE
	if clientID == gs.blackID.String() {
		ownColor = pb.Color_COLOR_BLACK
	}
	if ownColor == pb.Color_COLOR_WHITE {
		if gs.chess.Position.Ply <= 1 {
			result = pb.GameResult_GAME_RESULT_INTERRUPTED
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			result = pb.GameResult_GAME_RESULT_BLACK_WON
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION
			state = pb.GameState_GAME_STATE_FINISHED
		}
	}
	if ownColor == pb.Color_COLOR_BLACK {
		if gs.chess.Position.Ply < 2 {
			result = pb.GameResult_GAME_RESULT_INTERRUPTED
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			state = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			result = pb.GameResult_GAME_RESULT_WHITE_WON
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION
			state = pb.GameState_GAME_STATE_FINISHED
		}
	}
	updateGameDto := dto.GameChangeset{
		ResultID:       opt.New(h.mappings.results[result]),
		ResultStatusID: opt.New(h.mappings.resultStatuses[resultStatus]),
		StateID:        opt.New(h.mappings.states[state]),
		EndTime:        opt.New(time.Now()),
	}
	if _, err := h.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
		h.log.Error("GameResign, failed to persist game", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().RemoveActiveGame(context.Background(), gs.gameID); err != nil {
		h.log.Error("GameResign, failed to remove active game presence", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.whiteID); err != nil {
		h.log.Error("GameResign, failed to remove active game presence 1", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.blackID); err != nil {
		h.log.Error("GameResign, failed to remove active game presence 2", slog.Any("error", err))
		return
	}
	h.removeGame(gs.gameID)
	gameFinishedMsg := &pb.Message{
		Event: &pb.Message_GameFinished{GameFinished: &pb.GameFinished{GameResult: result, GameResultStatus: resultStatus, GameState: state}},
	}
	bb, err := protojson.Marshal(gameFinishedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameFinished", slog.Any("error", err))
		return
	}
	h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
}

func (h *Hub) handleGameDeclineDraw(msg *pb.GameDeclineDraw, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameDeclineDraw msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameDeclineDraw: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("GameDeclineDraw: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("GameDeclineDraw: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("GameDeclineDraw: can't offer draw in finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	otherID := gs.blackID
	if clientID == gs.blackID.String() {
		otherID = gs.whiteID
	}
	gameDeclineDrawMsg := &pb.Message{
		Event: &pb.Message_GameDeclineDraw{GameDeclineDraw: msg},
	}
	bb, err := protojson.Marshal(gameDeclineDrawMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameDeclineDraw", slog.Any("error", err))
		return
	}
	h.broadcastClient <- ClientMessage{clientID: otherID, msg: bb}
}

func (h *Hub) handleAcceptDrawMsg(msg *pb.GameAcceptDraw, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameAcceptDraw msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameAcceptDraw: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("GameAcceptDraw: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("GameAcceptDraw: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("GameAcceptDraw: can't offer draw in a finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	result := pb.GameResult_GAME_RESULT_DRAW
	resultStatus := pb.GameResultStatus_GAME_RESULT_STATUS_DRAW_AGREED
	state := pb.GameState_GAME_STATE_FINISHED
	updateGameDto := dto.GameChangeset{
		ResultID:       opt.New(h.mappings.results[result]),
		ResultStatusID: opt.New(h.mappings.resultStatuses[resultStatus]),
		StateID:        opt.New(h.mappings.states[state]),
		EndTime:        opt.New(time.Now()),
	}
	if _, err := h.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
		h.log.Error("GameAcceptDraw, failed to persist game", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().RemoveActiveGame(context.Background(), gs.gameID); err != nil {
		h.log.Error("GameAcceptDraw, failed to remove active game presence", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.whiteID); err != nil {
		h.log.Error("GameAcceptDraw, failed to remove active game presence 1", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.blackID); err != nil {
		h.log.Error("GameAcceptDraw, failed to remove active game presence 2", slog.Any("error", err))
		return
	}
	h.removeGame(gs.gameID)
	gameFinishedMsg := &pb.Message{
		Event: &pb.Message_GameFinished{GameFinished: &pb.GameFinished{GameResult: result, GameResultStatus: resultStatus, GameState: state}},
	}
	bb, err := protojson.Marshal(gameFinishedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameFinished", slog.Any("error", err))
		return
	}
	h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
}

func (h *Hub) handleGameChatMsg(msg *pb.GameChat, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil GameChat msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("GameChat: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	values := map[string]any{"client_id": clientID, "msg": msg.Message}
	res, err := h.rdb.XAdd(context.Background(), &redis.XAddArgs{Stream: fmt.Sprintf("game-chat.%s", gameID.String()), Values: values}).Result()
	if err != nil {
		h.log.Error("GameChat: stream add", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gameChatReceiveMsg := &pb.Message{
		Event: &pb.Message_GameChatReceive{GameChatReceive: &pb.GameChatReceive{ClientId: clientID, Id: res, Message: msg.Message}},
	}
	bb, err := protojson.Marshal(gameChatReceiveMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameChatReceive", slog.Any("error", err))
		return
	}
	h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
}

func (h *Hub) handlePlayMoveUCIMsg(msg *pb.PlayMoveUCI, clientID string, authState ClientAuthState) {
	if msg == nil {
		h.log.Error("nil PlayMoveUCI msg")
		return
	}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), uuid.MustParse(clientID))
	if err != nil {
		h.log.Error("PlayMoveUCI: no gameID presence", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err != nil {
		return
	}
	if authState != gs.authState {
		h.log.Error("PlayMoveUCI: authState mismatch", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if !(gs.whiteID.String() == clientID || gs.blackID.String() == clientID) {
		h.log.Error("PlayMoveUCI: player not in a game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		h.log.Error("PlayMoveUCI: play a finished game", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()))
		return
	}
	pcolor := pb.Color_COLOR_WHITE
	if clientID == gs.blackID.String() {
		pcolor = pb.Color_COLOR_BLACK
	}
	if pcolor == pb.Color_COLOR_WHITE && gs.chess.Position.Turn.IsBlack() || pcolor == pb.Color_COLOR_BLACK && gs.chess.Position.Turn.IsWhite() {
		h.log.Error("PlayMoveUCI: wrong turn", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()), slog.String("move", msg.Move))
		return
	}
	m, err := gs.chess.MakeMoveUCI(msg.GetMove())
	if err != nil {
		h.log.Error("PlayMoveUCI: invalid move", slog.String("client_id", clientID), slog.String("game_id", gameID.String()), slog.String("auth_state", authState.String()), slog.String("move", msg.Move))
		return
	}
	uci := m.ToUCI()
	san := m.ToSAN(gs.chess.Position, gs.chess.Position.Check, gs.chess.IsCheckmate(), gs.chess.LegalMoves)
	fen := gs.chess.Position.Fen()
	gs.moves = append(gs.moves, &pb.HistoryMoveInfo{
		Fen:      fen,
		Check:    gs.chess.Position.Check,
		Move:     &pb.HistoryMove{Uci: uci, San: san},
		PlayedAt: timestamppb.Now(),
	})
	if gs.chess.Position.Ply == 1 {
		go func() {
			gs.timerAction <- timerAction{timerEvent: stopWhiteFirstMoveTimer}
			gs.timerAction <- timerAction{timerEvent: startBlackFirstMoveTimer}
		}()
	}
	if gs.chess.Position.Ply == 2 {
		go func() {
			gs.timerAction <- timerAction{timerEvent: stopBlackFirstMoveTimer}
		}()
	}
	if gs.chess.Position.Ply >= 2 {
		gs.toggleClockAfterMove()
	}
	legalMoves := make([]string, 0, len(gs.chess.LegalMoves))
	for _, m := range gs.chess.LegalMoves {
		legalMoves = append(legalMoves, fmt.Sprint(m.String()))
	}
	lan := m.ToLAN(gs.chess.Position, gs.chess.Position.Check, gs.chess.IsCheckmate())
	receieveMoveMsg := &pb.Message{
		Event: &pb.Message_ReceiveMove{ReceiveMove: &pb.ReceiveMove{
			Uci:        uci,
			Lan:        lan,
			San:        san,
			Fen:        fen,
			Ply:        uint32(gs.chess.Position.Ply),
			LegalMoves: legalMoves,
			Clocks: &pb.Clocks{
				White: durationpb.New(gs.gameClock.White.Remaining()),
				Black: durationpb.New(gs.gameClock.Black.Remaining()),
			},
		}},
	}
	bb, err := protojson.Marshal(receieveMoveMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_ReceiveMove", slog.Any("error", err))
		return
	}
	h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
	if gs.chess.IsTerminated() {
		var result pb.GameResult
		var resultStatus pb.GameResultStatus
		state := pb.GameState_GAME_STATE_FINISHED
		switch gs.chess.Status() {
		case engine.StatusInsufficientMaterial:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_INSUFFICIENT_MATERIAL
		case engine.StatusThreeFoldRepetition:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_THREEFOLD_REPETITION
		case engine.StatusFiveFoldRepetition:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_FIVEFOLD_REPETITION
		case engine.StatusFiftyMoveRule:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_FIFTY_MOVE_RULE
		case engine.StatusSeventyFiveMoveRule:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_SEVENTYFIVE_MOVE_RULE
		case engine.StatusCheckmate:
			if pcolor == pb.Color_COLOR_WHITE {
				result = pb.GameResult_GAME_RESULT_WHITE_WON
			} else {
				result = pb.GameResult_GAME_RESULT_BLACK_WON
			}
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE
		case engine.StatusStalemate:
			result = pb.GameResult_GAME_RESULT_DRAW
			resultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE
		}
		updateGameDto := dto.GameChangeset{
			ResultID:       opt.New(h.mappings.results[result]),
			ResultStatusID: opt.New(h.mappings.resultStatuses[resultStatus]),
			StateID:        opt.New(h.mappings.states[state]),
			EndTime:        opt.New(time.Now()),
			LastMove:       opt.New(time.Now()),
		}
		if _, err := h.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
			h.log.Error("PlayMoveUCI, failed to persist game", slog.Any("error", err))
			return
		}
		if err := h.store.Presence().RemoveActiveGame(context.Background(), gs.gameID); err != nil {
			h.log.Error("PlayMoveUCI, failed to remove active game presence", slog.Any("error", err))
			return
		}
		if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.whiteID); err != nil {
			h.log.Error("PlayMoveUCI, failed to remove active game presence 1", slog.Any("error", err))
			return
		}
		if err := h.store.Presence().DelPlayerGameID(context.Background(), gs.blackID); err != nil {
			h.log.Error("PlayMoveUCI, failed to remove active game presence 2", slog.Any("error", err))
			return
		}
		h.removeGame(gs.gameID)
		gameFinishedMsg := &pb.Message{
			Event: &pb.Message_GameFinished{GameFinished: &pb.GameFinished{GameResult: result, GameResultStatus: resultStatus, GameState: state}},
		}
		bb, err := protojson.Marshal(gameFinishedMsg)
		if err != nil {
			h.log.Error("failed to protojson marshal Message_GameFinished", slog.Any("error", err))
			return
		}
		h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
	} else {
		updateGameDto := dto.GameChangeset{
			LastMove: opt.New(time.Now()),
		}
		if _, err := h.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
			h.log.Error("PlayMoveUCI, failed to persist game", slog.Any("error", err))
			return
		}
	}
}

func (h *Hub) getHubInfo() ([]byte, error) {
	hubInfoMsg := &pb.Message{
		Event: &pb.Message_HubInfo{HubInfo: &pb.HubInfo{Lobby: int32(h.lobbyClientsCount()), Playing: int32(h.gameClientsCount())}},
	}
	bb, err := protojson.Marshal(hubInfoMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_HubInfo", slog.Any("error", err))
		return nil, err
	}
	return bb, nil
}

func (h *Hub) broadcastHubInfoToLobby() {
	bb, err := h.getHubInfo()
	if err == nil {
		h.broadcastChannel <- ChannelMessage{channel: "lobby", msg: bb}
	}
}

func (h *Hub) broadcastHubInfoToClient(clientID uuid.UUID) {
	bb, err := h.getHubInfo()
	if err == nil {
		h.broadcastClient <- ClientMessage{clientID: clientID, msg: bb}
	}
}

func (h *Hub) sendInitChannelInfo(c *client) {
	for _, channel := range c.channels {
		if channel == "lobby" {
			h.sendInitLobbyInfo(c)
		} else if strings.HasPrefix(channel.String(), "game.") {
			gameID := strings.Split(channel.String(), ".")[1]
			gameuuid := uuid.MustParse(gameID)
			h.sendInitGameInfo(c, gameuuid)
		}
	}
}

func (h *Hub) sendInitLobbyInfo(c *client) {
	h.broadcastHubInfoToLobby()
}

func (h *Hub) sendInitGameInfo(c *client, gameID uuid.UUID) {
	// TODO: check again what to do
	gs, ok := h.games[gameID]
	if !ok {
		h.log.Error("sendInitGameInfo no game")
		return
	}
	if !(gs.whiteID == c.id || gs.blackID == c.id) {
		h.log.Error("sendInitGameInfo not in game")
		return
	}
	go func() {
		if gs.getPlayerColorFromID(c.id) == pb.Color_COLOR_WHITE {
			gs.timerAction <- timerAction{timerEvent: stopWhiteReconnectTimer}
		} else {
			gs.timerAction <- timerAction{timerEvent: stopBlackReconnectTimer}
		}
	}()
	otherID := gs.whiteID
	if c.id == gs.whiteID {
		otherID = gs.blackID
	}
	clientConnectedMsg := &pb.Message{
		Event: &pb.Message_ClientConnected{ClientConnected: &pb.ClientConnected{Id: c.id.String()}},
	}
	bb1, err := protojson.Marshal(clientConnectedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_ClientConnected", slog.Any("error", err))
		return
	}
	h.broadcastClient <- ClientMessage{clientID: otherID, msg: bb1}
	pcolor := pb.Color_COLOR_WHITE
	if gs.blackID == c.id {
		pcolor = pb.Color_COLOR_BLACK
	}
	fen := gs.chess.Position.Fen()
	ply := gs.chess.Position.Ply
	clocks := &pb.Clocks{
		White: durationpb.New(gs.gameClock.White.Remaining()),
		Black: durationpb.New(gs.gameClock.Black.Remaining()),
	}
	legalMoves := make([]string, 0, len(gs.chess.LegalMoves))
	for _, m := range gs.chess.LegalMoves {
		legalMoves = append(legalMoves, fmt.Sprint(m.String()))
	}
	matchFoundInner := &pb.MatchFound{GameId: gs.gameID.String(), ClientId: c.id.String(), Color: pcolor, Fen: fen, Ply: uint32(ply), Clocks: clocks, LegalMoves: legalMoves, GameState: gs.state, TimeControl: gs.timeControl, ReconnectTimeoutMs: int64(defaultReconnectTimeout), FirstMoveTimeoutMs: int64(defaultFirstMoveTimeout), HistoryMoveInfos: gs.moves, StartTime: timestamppb.New(gs.startTime)}
	if gs.authState == ClientAuth {
		oppUserInfo, err := h.userInfoFetcher.FetchUserInfo(context.Background(), otherID)
		if err != nil {
			h.log.Error("failed to fetch opponent user info", slog.Any("error", err))
			return
		}
		matchFoundInner.OpponentInfo = &pb.OpponentInfo{Username: oppUserInfo.Username, AvatarUrl: oppUserInfo.AvatarURL, Rating: 1500}
	}
	matchRejoinedMsg := &pb.Message{Event: &pb.Message_MatchFound{MatchFound: matchFoundInner}}
	bb2, err := protojson.Marshal(matchRejoinedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_MatchFound", slog.Any("error", err))
		return
	}
	h.broadcastClient <- ClientMessage{clientID: c.id, msg: bb2}
	lastChatID, err := h.getGameLastChatId(gs.gameID.String())
	if err != nil {
		h.log.Error("failed to get last chat id", slog.Any("error", err))
	}
	lastChatIDParam := c.query.Get("last_chat_id")
	if lastChatID == lastChatIDParam {
		return
	}
	start := "-"
	if lastChatIDParam != "" {
		start = "(" + lastChatIDParam
	}
	msgs, err := h.rdb.XRangeN(context.Background(), "game-chat."+gameID.String(), start, "+", 100).Result()
	if err != nil {
		h.log.Error("failed to fetch chat", slog.Any("err", err))
		return
	}
	gameChat := make([]*pb.GameChatReceive, 0)
	for _, m := range msgs {
		cid, ok1 := m.Values["client_id"].(string)
		msg, ok2 := m.Values["msg"].(string)
		if !ok1 || !ok2 {
			h.log.Error("failed to parse client_id or msg", slog.Any("err", err))
			return
		}
		gameChat = append(gameChat, &pb.GameChatReceive{
			Id:       m.ID,
			ClientId: cid,
			Message:  msg,
		})
	}
	gameChatRetrieveMsg := &pb.Message{
		Event: &pb.Message_GameChatRetrieve{GameChatRetrieve: &pb.GameChatRetrieve{GameChat: gameChat}},
	}
	bb, err := protojson.Marshal(gameChatRetrieveMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_GameChatRetrieve", slog.Any("error", err))
		return
	}
	h.broadcastClient <- ClientMessage{clientID: c.id, msg: bb}
}

func (h *Hub) assignChannelsToClient(c *client) {
	gameID, err := h.store.Presence().GetPlayerGameID(context.TODO(), c.id)
	if err != nil {
		c.channels = []Channel{"lobby"}
	} else {
		c.channels = []Channel{Channel("game." + gameID.String())}
	}
}

func (h *Hub) addClient(c *client) {
	h.mu.Lock()
	h.clientsByUserID[c.id] = c
	h.clients[c] = make([]Channel, 0)
	for _, channel := range c.channels {
		if h.channels[channel] == nil {
			h.channels[channel] = make(map[*client]struct{})
		}
		h.channels[channel][c] = struct{}{}
		h.clients[c] = append(h.clients[c], channel)
	}
	h.mu.Unlock()
}

func (h *Hub) removeClient(c *client) {
	close(c.egress)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, channel := range h.clients[c] {
		delete(h.channels[channel], c)
		if len(h.channels[channel]) == 0 {
			delete(h.channels, channel)
		}
	}
	delete(h.clients, c)
	delete(h.clientsByUserID, c.id)
}

func (h *Hub) managePlayerAddedToGame(c *client, gameID uuid.UUID) {
	// @TODO: check later
	c.channels = []Channel{Channel("game." + gameID.String())}

	h.mu.Lock()
	h.clients[c] = c.channels
	delete(h.channels["lobby"], c)
	for _, channel := range c.channels {
		if h.channels[channel] == nil {
			h.channels[channel] = make(map[*client]struct{})
		}
		if _, exists := h.channels[channel][c]; !exists {
			h.channels[channel][c] = struct{}{}
		}
	}
	h.mu.Unlock()
}

func (h *Hub) addGame(gs *gameState) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.games[gs.gameID] = gs
}

func (h *Hub) removeGame(gameID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.games, gameID)
}

func (h *Hub) lobbyClientsCount() int {
	return len(h.channels["lobby"])
}

func (h *Hub) gameClientsCount() int {
	var gameClientsCount int
	for channel := range h.channels {
		if strings.HasPrefix(channel.String(), "game.") {
			gameClientsCount += 2
		}
	}
	return gameClientsCount
}

func (h *Hub) tryMatchPlayers(ctx context.Context) {
	for _, qg := range QuickGames {
		h.tryMatchPlayersForSpecificTimeControl(ctx, ClientAuth, qg.ClockSecs, qg.IncrementSecs)
		h.tryMatchPlayersForSpecificTimeControl(ctx, ClientGuest, qg.ClockSecs, qg.IncrementSecs)
	}
}

func (h *Hub) tryMatchPlayersForSpecificTimeControl(ctx context.Context, authState ClientAuthState, clockSecs, incSecs int32) {
	key := fmt.Sprintf("seek_game.%d.%d_%d", authState, clockSecs, incSecs)
	res, err := h.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		h.log.Error("failed to fetch players from priority queue", slog.String("queue", key), slog.Any("error", err))
		return
	}
	queueSize := len(res)
	if queueSize < 2 {
		// h.log.Debug("not enough players to pair", slog.String("queue", key),)
		return
	}
	h.log.Debug("players in queue", slog.String("queue", key), slog.Int("count", queueSize))
	pairs := make(chan [2]string, queueSize/2)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i+1 < queueSize; i += 2 {
			pairs <- [2]string{res[i], res[i+1]}
		}
		close(pairs)
	}()
	workers := queueSize / 2
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pair := range pairs {
				h.processPair(ctx, pair, key, authState, clockSecs, incSecs)
			}
		}()
	}
	wg.Wait()
	if queueSize%2 != 0 {
		h.log.Debug("unmatched player remains in queue", slog.String("queue", key), slog.String("client_id", res[queueSize-1]))
	}
}

func (h *Hub) processPair(ctx context.Context, pair [2]string, key string, authState ClientAuthState, clockSecs, incSecs int32) {
	c1Id, err1 := uuid.Parse(pair[0])
	c2Id, err2 := uuid.Parse(pair[1])
	if err1 != nil || err2 != nil {
		h.log.Error("failed to parse client UUIDs", slog.String("c1", pair[0]), slog.String("c2", pair[1]))
		return
	}
	h.mu.RLock()
	c1, ok1 := h.clientsByUserID[c1Id]
	c2, ok2 := h.clientsByUserID[c2Id]
	if !ok1 || !ok2 {
		h.log.Error("failed to get client by id", slog.String("c1", pair[0]), slog.String("c2", pair[1]))
		return
	}
	h.mu.RUnlock()
	if c1.authState != c2.authState {
		h.log.Error("clients auth_state mismatch")
		return
	}
	timeControl := &pb.GameTimeControl{Clock: durationpb.New(time.Second * time.Duration(clockSecs)), Increment: durationpb.New(time.Second * time.Duration(incSecs))}
	gs, err := h.newGame(c1, c2, timeControl)
	if err != nil {
		h.log.Error("failed to create a game", slog.Any("error", err))
		return
	}
	h.addGame(gs)
	if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.ZRem(ctx, key, pair[0], pair[1]).Err(); err != nil {
			h.log.Error("failed to remove players from queue", slog.String("c1", pair[0]), slog.String("c2", pair[1]), slog.String("game_id", gs.gameID.String()), slog.Any("error", err))
			return err
		}
		if err := p.Publish(ctx, "game_found", gs.gameID.String()).Err(); err != nil {
			h.log.Error("failed to publish game found", slog.String("game_id", gs.gameID.String()), slog.Any("error", err))
		}
		return nil
	}); err != nil {
		h.removeGame(gs.gameID)
		h.log.Error("tryMatchPlayers pipeline", slog.String("game_id", gs.gameID.String()), slog.Any("error", err))
	}
	h.managePlayerAddedToGame(c1, gs.gameID)
	h.managePlayerAddedToGame(c2, gs.gameID)
	fen := gs.chess.Position.Fen()
	ply := gs.chess.Position.Ply
	clocks := &pb.Clocks{
		White: durationpb.New(gs.gameClock.White.Remaining()),
		Black: durationpb.New(gs.gameClock.Black.Remaining()),
	}
	legalMoves := make([]string, 0, len(gs.chess.LegalMoves))
	for _, m := range gs.chess.LegalMoves {
		legalMoves = append(legalMoves, fmt.Sprint(m.String()))
	}
	matchFoundInner1 := &pb.MatchFound{GameId: gs.gameID.String(), ClientId: c1.id.String(), Color: gs.getPlayerByID(c1.id).color, Fen: fen, Ply: uint32(ply), Clocks: clocks, LegalMoves: legalMoves, GameState: gs.state, TimeControl: gs.timeControl, ReconnectTimeoutMs: int64(defaultReconnectTimeout), FirstMoveTimeoutMs: int64(defaultFirstMoveTimeout), HistoryMoveInfos: gs.moves, StartTime: timestamppb.New(gs.startTime)}
	matchFoundInner2 := &pb.MatchFound{GameId: gs.gameID.String(), ClientId: c2.id.String(), Color: gs.getPlayerByID(c2.id).color, Fen: fen, Ply: uint32(ply), Clocks: clocks, LegalMoves: legalMoves, GameState: gs.state, TimeControl: gs.timeControl, ReconnectTimeoutMs: int64(defaultReconnectTimeout), FirstMoveTimeoutMs: int64(defaultFirstMoveTimeout), HistoryMoveInfos: gs.moves, StartTime: timestamppb.New(gs.startTime)}
	if authState == ClientAuth {
		userInfo1, err := h.userInfoFetcher.FetchUserInfo(ctx, c1.id)
		if err != nil {
			h.log.Error("failed to fetch user info 1", slog.Any("error", err2))
			return
		}
		userInfo2, err := h.userInfoFetcher.FetchUserInfo(ctx, c2.id)
		if err != nil {
			h.log.Error("failed to fetch user info 2", slog.Any("error", err2))
			return
		}
		matchFoundInner1.OpponentInfo = &pb.OpponentInfo{Username: userInfo2.Username, AvatarUrl: userInfo2.AvatarURL, Rating: 1500}
		matchFoundInner2.OpponentInfo = &pb.OpponentInfo{Username: userInfo1.Username, AvatarUrl: userInfo1.AvatarURL, Rating: 1500}
	}
	matchFound1 := &pb.Message{Event: &pb.Message_MatchFound{MatchFound: matchFoundInner1}}
	matchFound2 := &pb.Message{Event: &pb.Message_MatchFound{MatchFound: matchFoundInner2}}
	b1, err1 := protojson.Marshal(matchFound1)
	if err1 != nil {
		h.log.Error("failed to protojson marshal Message_MatchFound 1", slog.Any("error", err1))
	}
	b2, err2 := protojson.Marshal(matchFound2)
	if err2 != nil {
		h.log.Error("failed to protojson marshal Message_MatchFound 2", slog.Any("error", err2))
	}
	createGameDto := dto.GameChangeset{
		ID:                   opt.New(gs.gameID),
		VariantID:            opt.New(h.mappings.variants[gs.variant]),
		TimeKindID:           opt.New(h.mappings.timeKinds[gs.timeKind]),
		TimeCategoryID:       opt.New(h.mappings.timeCategories[gs.timeCategory]),
		TimeControlClock:     opt.New(clockSecs),
		TimeControlIncrement: opt.New(incSecs),
		ReconnectTimeout:     opt.New(int32(gs.reconnectTimeout.Seconds())),
		FirstMoveTimeout:     opt.New(int32(gs.firstMoveTimeout.Seconds())),
		StateID:              opt.New(h.mappings.states[gs.state]),
		StartTime:            opt.New(time.Now()),
		Fen:                  opt.New(gs.chess.Position.Fen()),
	}
	if c1.authState == ClientGuest {
		createGameDto.GuestWhiteID = opt.New(gs.whiteID)
		createGameDto.GuestBlackID = opt.New(gs.blackID)
	} else {
		createGameDto.WhiteID = opt.New(gs.whiteID)
		createGameDto.BlackID = opt.New(gs.blackID)
	}
	if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		createGameDto.ResultID = opt.New(h.mappings.results[gs.result])
	}
	if gs.resultStatus != pb.GameResultStatus_GAME_RESULT_STATUS_UNSPECIFIED {
		createGameDto.ResultStatusID = opt.New(h.mappings.resultStatuses[gs.resultStatus])
	}
	movesChangesets := make([]dto.GameMoveChangeset, 0, len(gs.moves)-1)
	for _, hist := range gs.moves {
		if hist.Move == nil {
			continue
		}
		movesChangesets = append(movesChangesets, dto.GameMoveChangeset{
			GameID:   opt.New(gs.gameID),
			Fen:      opt.New(gs.chess.Position.Fen()),
			Uci:      opt.New(hist.Move.Uci),
			San:      opt.New(hist.Move.San),
			PlayedAt: opt.New(hist.PlayedAt.AsTime()),
		})
	}
	createdGame, err := h.store.Game().Create(ctx, createGameDto, movesChangesets)
	if err != nil {
		h.log.Debug("failed to create game", slog.Any("error", err))
		return
	}
	if _, err := h.store.Presence().SetActiveGame(ctx, createdGame); err != nil {
		h.log.Debug("failed to cache active game presence", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().SetPlayerGameID(ctx, gs.whiteID, gs.gameID); err != nil {
		h.log.Debug("failed to set white player game id presence", slog.Any("error", err))
		return
	}
	if err := h.store.Presence().SetPlayerGameID(ctx, gs.blackID, gs.gameID); err != nil {
		h.log.Debug("failed to set black player game id presence", slog.Any("error", err))
		return
	}
	go gs.Run(ctx)
	h.broadcastClient <- ClientMessage{clientID: c1.id, msg: b1}
	h.broadcastClient <- ClientMessage{clientID: c2.id, msg: b2}
	h.log.Debug("match found success")
}

func (h *Hub) onClientConnected(client *client) {
	h.log.Debug("client connected", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()))
	h.assignChannelsToClient(client)
	h.addClient(client)
	h.sendInitChannelInfo(client)
	clientConnectedMsg := &pb.Message{Event: &pb.Message_ClientConnected{ClientConnected: &pb.ClientConnected{Id: client.id.String()}}}
	b, err := protojson.Marshal(clientConnectedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_ClientConnected ", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", b).Err(); err != nil {
			h.log.Error("failed to publish client_connected message", slog.String("client_id", client.id.String()), slog.Any("error", err))
		}
	}
}

func (h *Hub) onClientDisconnected(client *client) {
	ctx := context.Background()
	h.log.Debug("client disconnected", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()))
	h.removeClient(client)
	cskey, cserr := h.rdb.HGet(ctx, "clients_seeking", client.id.String()).Result()
	if cserr == nil {
		if _, err := h.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
			if err := p.ZRem(ctx, cskey, client.id.String()).Err(); err != nil {
				h.log.Error("cancel_seek_game remove from guest queue", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()), slog.Any("error", err))
				return err
			}
			if err := p.HDel(ctx, "clients_seeking", client.id.String()).Err(); err != nil {
				h.log.Error("seek_game remove seek key for client", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()), slog.Any("error", err))
				return err
			}
			if err := p.Publish(ctx, cskey, client.id.String()).Err(); err != nil {
				h.log.Error("cancel_seek_game publish cancel seek game", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()), slog.Any("error", err))
				return err
			}
			return nil
		}); err != nil {
			h.log.Error("remove seeks pipeline", slog.String("client_id", client.id.String()), slog.String("auth_state", client.authState.String()), slog.Any("error", err))
		}
	}
	clientDisconnectedMsg := &pb.ClientDisconnected{Id: client.id.String()}
	gameID, err := h.store.Presence().GetPlayerGameID(context.Background(), client.id)
	if err != nil {
		h.log.Debug("ClientDisconnected: no gameID presence", slog.String("client_id", client.id.String()), slog.String("game_id", gameID.String()), slog.String("auth_state", client.authState.String()))
		return
	}
	gs, err := h.getGameState(gameID)
	if err == nil {
		if client.authState != gs.authState {
			h.log.Error("ClientDisconnected: authState mismatch", slog.String("client_id", client.id.String()), slog.String("game_id", gameID.String()), slog.String("auth_state", client.authState.String()))
			return
		}
		if !(gs.whiteID == client.id || gs.blackID == client.id) {
			h.log.Error("ClientDisconnected: player not in a game", slog.String("client_id", client.id.String()), slog.String("game_id", gameID.String()), slog.String("auth_state", client.authState.String()))
			return
		}
		if gs.result != pb.GameResult_GAME_RESULT_UNSPECIFIED {
			h.log.Error("ClientDisconnected: can't abort a finished game", slog.String("client_id", client.id.String()), slog.String("game_id", gameID.String()), slog.String("auth_state", client.authState.String()))
			return
		}
		go func() {
			if gs.getPlayerColorFromID(client.id) == pb.Color_COLOR_WHITE {
				gs.timerAction <- timerAction{timerEvent: startWhiteReconnectTimer}
			} else {
				gs.timerAction <- timerAction{timerEvent: startBlackReconnectTimer}
			}
		}()
		cdmsg := &pb.Message{Event: &pb.Message_ClientDisconnected{ClientDisconnected: clientDisconnectedMsg}}
		b, err := protojson.Marshal(cdmsg)
		if err != nil {
			h.log.Error("failed to protojson marshal Message_ClientDisconnected", slog.String("client_id", client.id.String()), slog.Any("error", err))
		}
		h.broadcastChannel <- ChannelMessage{channel: Channel("game." + gameID.String()), msg: b}
	}
	b, err := protojson.Marshal(clientDisconnectedMsg)
	if err != nil {
		h.log.Error("failed to protojson marshal Message_ClientDisconnected ipc", slog.String("client_id", client.id.String()), slog.Any("error", err))
	} else {
		if err := h.rdb.Publish(context.Background(), "ipc", b).Err(); err != nil {
			h.log.Error("failed to publish client_disconnected message", slog.String("client_id", client.id.String()), slog.Any("error", err))
		}
	}
}

func (h *Hub) broadcastToClient(clientID uuid.UUID, channel Channel, b []byte) {
	h.log.Debug("broadcastToClient", slog.String("client_id", clientID.String()), slog.String("channel", channel.String()), slog.String("msg", string(b)))
	c, ok := h.clientsByUserID[clientID]
	if !ok {
		h.log.Debug("broadcastToClient client not found", slog.String("client_id", clientID.String()), slog.String("msg", string(b)))
		return
	}
	select {
	case c.egress <- b:
	default:
		h.removeClient(c)
	}
}

func (h *Hub) broadcastToChannel(channel Channel, b []byte) {
	h.log.Debug("broadcastToChannel", slog.String("channel", channel.String()), slog.String("msg", string(b)))
	for c := range h.channels[channel] {
		select {
		case c.egress <- b:
		default:
			h.removeClient(c)
		}
	}
}

func (h *Hub) processClientMessage(c *client, b []byte) error {
	topic := fmt.Sprintf("wsc.%s.%d", c.id, c.authState)
	if err := h.rdb.Publish(context.Background(), topic, b).Err(); err != nil {
		h.log.Error("failed to publish client message", slog.String("client_id", c.id.String()), slog.String("topic", topic), slog.Any("error", err))
	}
	return nil
}

// extractWSCSubtopics extracts the client_id and auth_state
func extractWSCSubtopics(topic string) (string, ClientAuthState, error) {
	subtopics := strings.Split(topic, ".")
	if len(subtopics) != 3 {
		return "", ClientGuest, fmt.Errorf("invalid subtopics length, expected 2, got: %d", len(subtopics))
	}
	clientID, authStateStr := subtopics[1], subtopics[2]
	if !(authStateStr == "0" || authStateStr == "1") {
		return "", ClientGuest, fmt.Errorf("invalid subtopics auth_state. must be 0 or 1")
	}
	authState := ClientGuest
	if authStateStr == "1" {
		authState = ClientAuth
	}
	return clientID, authState, nil
}

// extractClientSubtopics extracts the client_id and the channel
func extractClientSubtopics(topic string) (uuid.UUID, string, error) {
	subtopics := strings.SplitN(topic, ".", 3)
	if len(subtopics) != 2 && len(subtopics) != 3 {
		return uuid.Nil, "", fmt.Errorf("invalid subtopics length, expected 2 or 3, got: %d", len(subtopics))
	}
	if len(subtopics) == 2 {
		clientIDRaw := subtopics[1]
		clientID, err := uuid.Parse(clientIDRaw)
		if err != nil {
			return uuid.Nil, "", fmt.Errorf("failed to parse client uuid")
		}
		return clientID, "", nil
	}
	clientIDRaw, channel := subtopics[1], subtopics[2]
	clientID, err := uuid.Parse(clientIDRaw)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to parse client uuid")
	}
	return clientID, channel, nil
}

// extractGameMessageInfo extracts the game_id
func extractGameMessageInfo(topic string) (string, error) {
	subtopics := strings.Split(topic, ".")
	if len(subtopics) != 2 {
		return "", fmt.Errorf("invalid subtopics length, expected 2 or 3, got: %d", len(subtopics))
	}
	gameID := subtopics[1]
	return gameID, nil
}

func (h *Hub) newGame(c1, c2 *client, timeControl *pb.GameTimeControl, gameOpts ...GameOption) (*gameState, error) {
	whiteID, blackID := c1.id, c2.id
	if rand.Intn(2) == 1 {
		whiteID, blackID = c2.id, c1.id
	}
	players := [2]PlayerInfo{
		{ID: whiteID, Name: "White", Color: pb.Color_COLOR_WHITE},
		{ID: blackID, Name: "Black", Color: pb.Color_COLOR_BLACK},
	}
	gs, err := NewGameState(players, timeControl, h, c1.authState, h.categoryThressholds, gameOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new gamestate: %w", err)
	}
	return gs, nil
}

func (h *Hub) getGameState(gameID uuid.UUID) (*gameState, error) {
	gs, ok := h.games[gameID]
	if !ok {
		// @TODO: load from redis
		return nil, fmt.Errorf("no game found")
	}
	return gs, nil
}

func (h *Hub) getGameLastChatId(gameID string) (string, error) {
	arr, err := h.rdb.XRevRangeN(context.Background(), "game-chat."+gameID, "+", "-", 1).Result()
	if err != nil {
		return "", fmt.Errorf("failed to fetch last chat msg id: %w", err)
	}
	if len(arr) == 0 {
		return "", nil
	}
	return arr[0].ID, nil
}
