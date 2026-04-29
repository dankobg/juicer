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

type CategoryThreshold struct {
	UpperLimit   time.Duration
	TimeCategory pb.GameTimeCategory
}

type Player struct {
	ID    uuid.UUID
	Name  string
	Color pb.Color
	Guest bool
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
	Chess            *engine.Chess
	GameID           int64
	White            *Player
	Black            *Player
	Players          map[uuid.UUID]*Player
	Guest            bool
	GameVariant      pb.GameVariant
	GameTimeKind     pb.GameTimeKind
	GameTimeCategory pb.GameTimeCategory
	// It could be much more complicated with many phases (e.g. increment only after move 30 etc)
	GameTimeControl  *pb.GameTimeControl
	GameResult       pb.GameResult
	GameResultStatus pb.GameResultStatus
	GameState        pb.GameState
	ReconnectTimeout time.Duration
	FirstMoveTimeout time.Duration
	LastMove         *time.Time
	StartTime        *time.Time
	EndTime          *time.Time
	Rated            bool
	HistoryMoveInfos []*pb.HistoryMoveInfo

	// ##########################################

	TimerAction     chan timerAction
	activeGameTimer *time.Timer

	whiteFirstMoveTimer *time.Timer
	blackFirstMoveTimer *time.Timer
	whiteReconnectTimer *time.Timer
	blackReconnectTimer *time.Timer
}

func NewGameState(gameID int64, players [2]Player, timeControl *pb.GameTimeControl, guest bool, thresholds []CategoryThreshold, opts ...GameOption) (*GameState, error) {
	if err := validatePlayers(players); err != nil {
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
		state:            pb.GameState_GAME_STATE_WAITING_START,
	}
	for _, o := range opts {
		o.apply(gopts)
	}

	chess, err := engine.NewChess(gopts.fen)
	if err != nil {
		return nil, fmt.Errorf("failed to init chess engine: %w", err)
	}

	var white, black *Player

	playersByID := make(map[uuid.UUID]*Player)

	for _, p := range players {
		if p.Color == pb.Color_COLOR_WHITE {
			white = &Player{ID: p.ID, Name: p.Name, Color: p.Color, Guest: p.Guest}
			playersByID[p.ID] = white
		} else {
			black = &Player{ID: p.ID, Name: p.Name, Color: p.Color, Guest: p.Guest}
			playersByID[p.ID] = black
		}
	}

	historyMoveInfos := []*pb.HistoryMoveInfo{{Fen: gopts.fen, Check: false, PlayedAt: nil, Move: nil}}

	gs := &GameState{
		Chess:            chess,
		GameID:           gameID,
		White:            white,
		Black:            black,
		Players:          playersByID,
		Guest:            guest,
		GameVariant:      gopts.variant,
		GameTimeKind:     gopts.timeKind,
		GameTimeCategory: gopts.timeCategory,
		GameTimeControl:  timeControl,
		GameState:        pb.GameState_GAME_STATE_IDLE,
		FirstMoveTimeout: gopts.firstMoveTimeout,
		ReconnectTimeout: gopts.reconnectTimeout,
		TimerAction:      make(chan timerAction),
		HistoryMoveInfos: historyMoveInfos,
		Rated:            gopts.rated,
	}

	return gs, nil
}

func (gs *GameState) GetPlayerByID(id uuid.UUID) *Player {
	return gs.Players[id]
}

func (gs *GameState) GetPlayerByColor(color pb.Color) *Player {
	switch color {
	case pb.Color_COLOR_WHITE:
		return gs.White
	case pb.Color_COLOR_BLACK:
		return gs.Black
	default:
		return nil
	}
}

func (gs *GameState) Start() {
	gs.GameState = pb.GameState_GAME_STATE_WAITING_START
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

func validatePlayers(players [2]Player) error {
	for _, p := range players {
		if p.ID == uuid.Nil {
			return errors.New("player id can't be empty")
		}

		if p.Color == pb.Color_COLOR_UNSPECIFIED {
			return errors.New("player color can't be empty")
		}
	}

	if players[0].Color == players[1].Color {
		return errors.New("player can't have the same color")
	}

	return nil
}

func determineGameTimeCategoryFromTimeControl(gtc *pb.GameTimeControl, thresholds []CategoryThreshold) (pb.GameTimeCategory, error) {
	if gtc == nil {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("time control is required")
	}

	if gtc.GetClockMs() <= 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("clock must be > 0")
	}

	if gtc.GetIncrementMs() < 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, errors.New("increment must be >= 0")
	}

	clock := time.Duration(gtc.GetClockMs()) * time.Millisecond
	increment := time.Duration(gtc.GetIncrementMs()) * time.Millisecond

	totalTime := clock + increment*time.Duration(AverageGameMovesEstimate)

	for _, threshold := range thresholds {
		if totalTime < threshold.UpperLimit {
			return threshold.TimeCategory, nil
		}
	}

	return thresholds[len(thresholds)-1].TimeCategory, nil
}
