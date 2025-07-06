package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dankobg/juicer/core/clock"
	"github.com/dankobg/juicer/db/dto"
	"github.com/dankobg/juicer/engine"
	"github.com/dankobg/juicer/opt"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	AverageGameMovesEstimate = 40
	defaultReconnectTimeout  = time.Second * 15
	defaultFirstMoveTimeout  = time.Second * 10
)

type categoryThreshold struct {
	upperLimit   time.Duration
	timeCategory pb.GameTimeCategory
}

type timerEvent = uint

const (
	startGameTimer timerEvent = iota
	pauseGameTimer
	stopGameTimer
	toggleGameTimer
	incrementWhiteGameTimer
	incrementBlackGameTimer
	startWhiteReconnectTimer
	startBlackReconnectTimer
	stopWhiteReconnectTimer
	stopBlackReconnectTimer
	startWhiteFirstMoveTimer
	startBlackFirstMoveTimer
	stopWhiteFirstMoveTimer
	stopBlackFirstMoveTimer
)

type timerAction struct {
	timerEvent timerEvent
	increment  *time.Duration
}

type player struct {
	id    uuid.UUID
	name  string
	color pb.Color
}

type gameState struct {
	chess               *engine.Chess
	gameID              uuid.UUID
	whiteID             uuid.UUID
	blackID             uuid.UUID
	authState           ClientAuthState
	variant             pb.Variant
	timeKind            pb.GameTimeKind
	timeCategory        pb.GameTimeCategory
	timeControl         *pb.GameTimeControl
	result              pb.GameResult
	resultStatus        pb.GameResultStatus
	state               pb.GameState
	reconnectTimeout    time.Duration
	firstMoveTimeout    time.Duration
	startTime           time.Time
	lastMove            time.Time
	moves               []*pb.HistoryMoveInfo
	timerAction         chan timerAction
	gameClock           *clock.Clock
	whiteFirstMoveTimer *time.Timer
	blackFirstMoveTimer *time.Timer
	whiteReconnectTimer *time.Timer
	blackReconnectTimer *time.Timer
	players             map[uuid.UUID]*player
	finishGame          chan *pb.GameFinished
	hub                 *Hub
}

type PlayerInfo struct {
	ID    uuid.UUID
	Name  string
	Color pb.Color
}

func NewGameState(playersInfo [2]PlayerInfo, timeControl *pb.GameTimeControl, hub *Hub, authState ClientAuthState, thresholds []categoryThreshold, opts ...GameOption) (*gameState, error) {
	if err := validatePlayersInfo(playersInfo); err != nil {
		return nil, err
	}
	timeCategory, err := determineGameTimeCategoryFromTimeControl(timeControl, thresholds)
	if err != nil {
		return nil, err
	}
	gopts := &gameOpts{
		fen:              engine.FENStartingPosition,
		variant:          pb.Variant_VARIANT_STANDARD,
		timeControl:      timeControl,
		timeKind:         pb.GameTimeKind_GAME_TIME_KIND_REALTIME,
		timeCategory:     timeCategory,
		reconnectTimeout: defaultReconnectTimeout,
		firstMoveTimeout: defaultFirstMoveTimeout,
	}
	for _, o := range opts {
		o.apply(gopts)
	}
	chess, err := engine.NewChess(gopts.fen)
	if err != nil {
		return nil, fmt.Errorf("failed to init chess engine: %w", err)
	}
	var whiteID, blackID uuid.UUID
	players := make(map[uuid.UUID]*player, 2)
	for _, p := range playersInfo {
		players[p.ID] = &player{
			id:    p.ID,
			name:  p.Name,
			color: p.Color,
		}
		if p.Color == pb.Color_COLOR_WHITE {
			whiteID = p.ID
		} else {
			blackID = p.ID
		}
	}
	moves := []*pb.HistoryMoveInfo{{Fen: gopts.fen, Check: false, PlayedAt: nil, Move: nil}}
	gs := &gameState{
		chess:            chess,
		gameID:           uuid.New(),
		whiteID:          whiteID,
		blackID:          blackID,
		authState:        authState,
		variant:          gopts.variant,
		timeKind:         gopts.timeKind,
		timeCategory:     gopts.timeCategory,
		timeControl:      timeControl,
		state:            pb.GameState_GAME_STATE_IN_PROGRESS,
		reconnectTimeout: gopts.reconnectTimeout,
		firstMoveTimeout: gopts.firstMoveTimeout,
		players:          players,
		finishGame:       make(chan *pb.GameFinished),
		timerAction:      make(chan timerAction),
		hub:              hub,
		gameClock:        clock.NewClock(gopts.timeControl.GetClock().AsDuration(), gopts.timeControl.GetIncrement().AsDuration()),
		moves:            moves,
		startTime:        time.Now(),
	}
	return gs, nil
}

func (gs *gameState) getPlayerByID(id uuid.UUID) *player {
	return gs.players[id]
}

func (gs *gameState) getPlayerIDFromColor(color pb.Color) uuid.UUID {
	if color == pb.Color_COLOR_WHITE {
		return gs.whiteID
	}
	return gs.blackID
}

func (gs *gameState) getPlayerColorFromID(id uuid.UUID) pb.Color {
	if gs.whiteID == id {
		return pb.Color_COLOR_WHITE
	}
	return pb.Color_COLOR_BLACK
}

func (gs *gameState) getPlayerByColor(color pb.Color) *player {
	return gs.players[gs.getPlayerIDFromColor(color)]
}

func determineGameTimeCategoryFromTimeControl(gtc *pb.GameTimeControl, thresholds []categoryThreshold) (pb.GameTimeCategory, error) {
	if gtc == nil {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("time control is required")
	}
	if gtc.GetClock().AsDuration() <= 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("clock must be > 0")
	}
	if gtc.GetIncrement().AsDuration() < 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("increment must be >= 0")
	}
	totalTime := gtc.GetClock().AsDuration() + gtc.GetIncrement().AsDuration()*AverageGameMovesEstimate
	for _, threshold := range thresholds {
		if totalTime < threshold.upperLimit {
			return threshold.timeCategory, nil
		}
	}
	return thresholds[len(thresholds)-1].timeCategory, nil
}

func validatePlayersInfo(playersInfo [2]PlayerInfo) error {
	for _, p := range playersInfo {
		if p.ID == uuid.Nil {
			return errors.New("player id can't be empty")
		}
		if p.Color == pb.Color_COLOR_UNSPECIFIED {
			return errors.New("player color can't be empty")
		}
	}
	if playersInfo[0].Color == playersInfo[1].Color {
		return errors.New("player can't have the same color")
	}
	return nil
}

func (gs *gameState) Run(ctx context.Context) {
	gs.hub.log.Debug("gamestate running", slog.String("game_id", gs.gameID.String()), slog.String("white", gs.whiteID.String()), slog.String("black", gs.blackID.String()))

	go func() {
		gs.timerAction <- timerAction{timerEvent: startWhiteFirstMoveTimer}
	}()

	go func() {
		for {
			select {
			// listen for timer action events
			case action := <-gs.timerAction:
				switch action.timerEvent {
				case startWhiteFirstMoveTimer:
					gs.hub.log.Debug("starting white first move timer", slog.String("game_id", gs.gameID.String()))
					gs.whiteFirstMoveTimer = time.NewTimer(gs.firstMoveTimeout)
				case startBlackFirstMoveTimer:
					gs.hub.log.Debug("starting black first move timer", slog.String("game_id", gs.gameID.String()))
					gs.blackFirstMoveTimer = time.NewTimer(gs.firstMoveTimeout)
				case stopWhiteFirstMoveTimer:
					gs.hub.log.Debug("stopping white first move timer", slog.String("game_id", gs.gameID.String()))
					if gs.whiteFirstMoveTimer != nil {
						gs.whiteFirstMoveTimer.Stop()
					}
					gs.whiteFirstMoveTimer = nil
				case stopBlackFirstMoveTimer:
					gs.hub.log.Debug("stopping black first move timer", slog.String("game_id", gs.gameID.String()))
					if gs.blackFirstMoveTimer != nil {
						gs.blackFirstMoveTimer.Stop()
					}
					gs.blackFirstMoveTimer = nil
				case startWhiteReconnectTimer:
					gs.hub.log.Debug("starting white reconnect timer", slog.String("game_id", gs.gameID.String()))
					gs.whiteReconnectTimer = time.NewTimer(gs.reconnectTimeout)
				case startBlackReconnectTimer:
					gs.hub.log.Debug("starting black reconnect timer", slog.String("game_id", gs.gameID.String()))
					gs.blackReconnectTimer = time.NewTimer(gs.reconnectTimeout)
				case stopWhiteReconnectTimer:
					gs.hub.log.Debug("stopping white reconnect timer", slog.String("game_id", gs.gameID.String()))
					if gs.whiteReconnectTimer != nil {
						gs.whiteReconnectTimer.Stop()
					}
					gs.whiteReconnectTimer = nil
				case stopBlackReconnectTimer:
					gs.hub.log.Debug("stopping black reconnect timer", slog.String("game_id", gs.gameID.String()))
					if gs.blackReconnectTimer != nil {
						gs.blackReconnectTimer.Stop()
					}
					gs.blackReconnectTimer = nil
				case startGameTimer:
					gs.hub.log.Debug("starting game timer", slog.String("game_id", gs.gameID.String()))
					gs.gameClock.Start()
				case pauseGameTimer:
					gs.hub.log.Debug("pausing game timer", slog.String("game_id", gs.gameID.String()))
					gs.gameClock.Pause()
				case toggleGameTimer:
					gs.hub.log.Debug("toggling game timer", slog.String("game_id", gs.gameID.String()))
					gs.gameClock.Toggle()
				case stopGameTimer:
					gs.hub.log.Debug("stopping game timer", slog.String("game_id", gs.gameID.String()))
					gs.gameClock.Reset()
				case incrementWhiteGameTimer:
					gs.hub.log.Debug("incrementing white game timer", slog.String("game_id", gs.gameID.String()), slog.Duration("increment", *action.increment))
					gs.gameClock.White.Add(*action.increment)
				case incrementBlackGameTimer:
					gs.hub.log.Debug("incrementing black game timer", slog.String("game_id", gs.gameID.String()), slog.Duration("increment", *action.increment))
					gs.gameClock.Black.Add(*action.increment)
				}
			// listen for game timers expiration
			case <-gs.gameClock.White.Tick():
				gs.hub.log.Debug("white game timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_BLACK_WON, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED, GameState: pb.GameState_GAME_STATE_FINISHED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-gs.gameClock.Black.Tick():
				gs.hub.log.Debug("black game timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_WHITE_WON, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED, GameState: pb.GameState_GAME_STATE_FINISHED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-clock.Tick(gs.whiteFirstMoveTimer):
				gs.hub.log.Debug("white first move timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_INTERRUPTED, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT, GameState: pb.GameState_GAME_STATE_INTERRUPTED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-clock.Tick(gs.blackFirstMoveTimer):
				gs.hub.log.Debug("black first move timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_INTERRUPTED, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT, GameState: pb.GameState_GAME_STATE_INTERRUPTED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-clock.Tick(gs.whiteReconnectTimer):
				gs.hub.log.Debug("white reconnect timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_BLACK_WON, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT, GameState: pb.GameState_GAME_STATE_FINISHED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-clock.Tick(gs.blackReconnectTimer):
				gs.hub.log.Debug("black reconnect timer expired", slog.String("game_id", gs.gameID.String()))
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_WHITE_WON, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT, GameState: pb.GameState_GAME_STATE_FINISHED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case gameFinishedMsg := <-gs.finishGame:
				gs.hub.log.Debug("adjudicating game", slog.String("game_id", gs.gameID.String()))
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			case <-ctx.Done():
				gs.hub.log.Debug("gamestate timer expiration listener ctx done")
				gameFinishedMsg := &pb.GameFinished{GameResult: pb.GameResult_GAME_RESULT_INTERRUPTED, GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED, GameState: pb.GameState_GAME_STATE_INTERRUPTED}
				gs.performGameFinishedTasks(gameFinishedMsg)
				return
			}
		}
	}()
}

func (gs *gameState) performGameFinishedTasks(msg *pb.GameFinished) {
	updateGameDto := dto.GameChangeset{
		ResultID:       opt.New(gs.hub.mappings.results[msg.GetGameResult()]),
		ResultStatusID: opt.New(gs.hub.mappings.resultStatuses[msg.GetGameResultStatus()]),
		StateID:        opt.New(gs.hub.mappings.states[msg.GetGameState()]),
		EndTime:        opt.New(time.Now()),
	}
	if _, err := gs.hub.store.Game().Update(context.Background(), gs.gameID, updateGameDto); err != nil {
		gs.hub.log.Error("GameFinish, failed to persist game", slog.Any("error", err))
		return
	}
	if err := gs.hub.store.Presence().RemoveActiveGame(context.Background(), gs.gameID); err != nil {
		gs.hub.log.Error("GameFinish, failed to remove active game presence", slog.Any("error", err))
		return
	}
	if err := gs.hub.store.Presence().DelPlayerGameID(context.Background(), gs.whiteID); err != nil {
		gs.hub.log.Error("GameFinish, failed to remove active game presence 1", slog.Any("error", err))
		return
	}
	if err := gs.hub.store.Presence().DelPlayerGameID(context.Background(), gs.blackID); err != nil {
		gs.hub.log.Error("GameFinish, failed to remove active game presence 2", slog.Any("error", err))
		return
	}
	gs.hub.removeGame(gs.gameID)
	message := &pb.Message{Event: &pb.Message_GameFinished{GameFinished: msg}}
	bb, err := protojson.Marshal(message)
	if err != nil {
		gs.hub.log.Error("failed to protojson marshal Message_GameFinished", slog.Any("error", err))
		return
	}
	gs.hub.broadcastChannel <- ChannelMessage{channel: Channel("game." + gs.gameID.String()), msg: bb}
}

func (gs *gameState) Stop() {
	gs.gameClock.Reset()
	if gs.whiteFirstMoveTimer != nil {
		gs.whiteFirstMoveTimer.Stop()
	}
	if gs.whiteReconnectTimer != nil {
		gs.whiteReconnectTimer.Stop()
	}
	if gs.blackFirstMoveTimer != nil {
		gs.blackFirstMoveTimer.Stop()
	}
	if gs.blackReconnectTimer != nil {
		gs.blackReconnectTimer.Stop()
	}
}

func (gs *gameState) toggleClockAfterMove() {
	go func() {
		gs.timerAction <- timerAction{timerEvent: toggleGameTimer}
	}()
	gs.lastMove = time.Now()
}
