package gameplay

import (
	"errors"
	"fmt"
	"time"

	"github.com/dankobg/juicer/engine"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
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

type player struct {
	id    uuid.UUID
	name  string
	color pb.Color
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

type GameState struct {
	chess            *engine.Chess
	gameID           int64
	whiteID          uuid.UUID
	blackID          uuid.UUID
	white            *player
	black            *player
	guest            bool
	gameVariant      pb.GameVariant
	gameTimeKind     pb.GameTimeKind
	gameTimeCategory pb.GameTimeCategory
	gameTimeControl  *pb.GameTimeControl
	gameResult       pb.GameResult
	gameResultStatus pb.GameResultStatus
	gameState        pb.GameState
	reconnectTimeout time.Duration
	firstMoveTimeout time.Duration
	startTime        time.Time
	lastMove         time.Time
	historyMoveInfos []*pb.HistoryMoveInfo

	// ##########################################

	timerAction chan timerAction

	whiteGameTimeRemaining time.Duration
	blackGameTimeRemaining time.Duration
	whiteIncrement         time.Duration
	blackIncrement         time.Duration
	activeGameTimer        *time.Timer

	whiteFirstMoveTimer *time.Timer
	blackFirstMoveTimer *time.Timer
	whiteReconnectTimer *time.Timer
	blackReconnectTimer *time.Timer
}

type PlayerInfo struct {
	ID    uuid.UUID
	Name  string
	Color pb.Color
}

func NewGameState(gameID int64, playersInfo [2]PlayerInfo, timeControl *pb.GameTimeControl, guest bool, thresholds []categoryThreshold, opts ...GameOption) (*GameState, error) {
	if err := validatePlayersInfo(playersInfo); err != nil {
		return nil, err
	}

	timeCategory, err := determineGameTimeCategoryFromTimeControl(timeControl, thresholds)
	if err != nil {
		return nil, err
	}

	gopts := &gameOpts{
		fen:              engine.FENStartingPosition,
		variant:          pb.GameVariant_GAME_VARIANT_STANDARD,
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
	var white, black *player
	for _, p := range playersInfo {
		if p.Color == pb.Color_COLOR_WHITE {
			whiteID = p.ID
			white = &player{id: p.ID, name: p.Name, color: p.Color}
		} else {
			blackID = p.ID
			black = &player{id: p.ID, name: p.Name, color: p.Color}
		}
	}

	historyMoveInfos := []*pb.HistoryMoveInfo{{Fen: gopts.fen, Check: false, PlayedAt: nil, Move: nil}}

	gs := &GameState{
		chess:            chess,
		gameID:           gameID,
		whiteID:          whiteID,
		blackID:          blackID,
		white:            white,
		black:            black,
		guest:            guest,
		gameVariant:      gopts.variant,
		gameTimeKind:     gopts.timeKind,
		gameTimeCategory: gopts.timeCategory,
		gameTimeControl:  timeControl,
		gameState:        pb.GameState_GAME_STATE_WAITING_START,
		firstMoveTimeout: gopts.firstMoveTimeout,
		reconnectTimeout: gopts.reconnectTimeout,
		timerAction:      make(chan timerAction),
		historyMoveInfos: historyMoveInfos,
		// maybe later i can have different time control or increment per player, but who cares...
		whiteGameTimeRemaining: gopts.timeControl.GetClock().AsDuration(),
		blackGameTimeRemaining: gopts.timeControl.GetClock().AsDuration(),
		whiteIncrement:         gopts.timeControl.GetIncrement().AsDuration(),
		blackIncrement:         gopts.timeControl.GetIncrement().AsDuration(),
	}

	return gs, nil
}

func (gs *GameState) Stop() {
	if gs.activeGameTimer != nil {
		gs.activeGameTimer.Stop()
	}
	if gs.whiteFirstMoveTimer != nil {
		gs.whiteFirstMoveTimer.Stop()
	}
	if gs.blackFirstMoveTimer != nil {
		gs.blackFirstMoveTimer.Stop()
	}
	if gs.whiteReconnectTimer != nil {
		gs.whiteReconnectTimer.Stop()
	}
	if gs.blackReconnectTimer != nil {
		gs.blackReconnectTimer.Stop()
	}
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
