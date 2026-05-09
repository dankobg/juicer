package gameplay

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/dankobg/juicer/engine"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/goforj/godump"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrChessEngineInit = errors.New("failed to init chess engine")

	ErrTimeControlMissing          = errors.New("time control is required")
	ErrTimeControlClockInvalid     = errors.New("clock must be > 0")
	ErrTimeControlIncrementInvalid = errors.New("increment must be >= 0")

	ErrPlayerIdMissing    = errors.New("player id can't be empty")
	ErrPlayerColorMissing = errors.New("player color can't be empty")
	ErrPlayersSameColors  = errors.New("player can't have the same color")

	ErrGameAlreadyConcluded = errors.New("game already concluded")
	ErrPlayerNotInGame      = errors.New("player not in game")
	ErrNotYourTurn          = errors.New("not your turn")
	ErrInvalidMove          = errors.New("invalid move attempt")
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
	ID       uuid.UUID
	Username string
	Color    pb.Color
	Guest    bool
}

type GameState struct {
	Chess            *engine.Chess
	GameID           int64
	White            *Player
	Black            *Player
	Players          map[uuid.UUID]*Player
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
	GameMoves        []*pb.GameMove
	running          atomic.Bool

	GameCommand         chan GameCommand
	GameEvent           chan GameEvent
	activeGameTimer     *time.Timer
	firstMoveTimer      *time.Timer
	whiteReconnectTimer *time.Timer
	blackReconnectTimer *time.Timer
}

func NewGameState(gameID int64, players [2]Player, timeControl *pb.GameTimeControl, thresholds []CategoryThreshold, gameEvent chan GameEvent, opts ...GameOption) (*GameState, error) {
	if err := validatePlayers(players); err != nil {
		return nil, err
	}

	timeCategory, err := determineGameTimeCategoryFromTimeControl(timeControl, thresholds)
	if err != nil {
		return nil, err
	}

	gopts := &gameOpts{
		gameID:           gameID,
		fen:              engine.FENStartingPosition,
		gameVariant:      pb.GameVariant_GAME_VARIANT_STANDARD,
		gameTimeControl:  timeControl,
		gameTimeKind:     pb.GameTimeKind_GAME_TIME_KIND_REALTIME,
		gameTimeCategory: timeCategory,
		gameState:        pb.GameState_GAME_STATE_ACTIVE,
		reconnectTimeout: defaultReconnectTimeout,
		firstMoveTimeout: defaultFirstMoveTimeout,
		startTime:        new(time.Now()),
	}
	for _, o := range opts {
		o.apply(gopts)
	}

	chess, err := engine.NewChess(gopts.fen)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChessEngineInit, err)
	}

	var white, black *Player

	playersByID := make(map[uuid.UUID]*Player)

	for _, p := range players {
		if p.Color == pb.Color_COLOR_WHITE {
			white = &Player{ID: p.ID, Username: p.Username, Color: p.Color, Guest: p.Guest}
			playersByID[p.ID] = white
		} else {
			black = &Player{ID: p.ID, Username: p.Username, Color: p.Color, Guest: p.Guest}
			playersByID[p.ID] = black
		}
	}

	gameMoves := []*pb.GameMove{{Fen: gopts.fen, Check: false}}

	gs := &GameState{
		Chess:            chess,
		GameID:           gopts.gameID,
		White:            white,
		Black:            black,
		Players:          playersByID,
		GameVariant:      gopts.gameVariant,
		GameTimeKind:     gopts.gameTimeKind,
		GameTimeCategory: gopts.gameTimeCategory,
		GameTimeControl:  gopts.gameTimeControl,
		GameState:        gopts.gameState,
		GameResult:       gopts.gameResult,
		GameResultStatus: gopts.gameResultStatus,
		FirstMoveTimeout: gopts.firstMoveTimeout,
		ReconnectTimeout: gopts.reconnectTimeout,
		Rated:            gopts.rated,
		StartTime:        gopts.startTime,
		LastMove:         gopts.lastMove,
		EndTime:          gopts.endTime,
		GameMoves:        gameMoves,
		running:          atomic.Bool{},

		GameCommand: make(chan GameCommand, 32),
		GameEvent:   gameEvent,
		// activeGameTimer     *time.Timer
		// firstMoveTimer      *time.Timer
		// whiteReconnectTimer *time.Timer
		// blackReconnectTimer *time.Timer
	}

	return gs, nil
}

func (gs *GameState) GetPlayerByID(id uuid.UUID) *Player {
	return gs.Players[id]
}

func (gs *GameState) HasGamePlayer(id uuid.UUID) bool {
	return gs.White.ID == id || gs.Black.ID == id
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

func (gs *GameState) Start(ctx context.Context) {
	if gs.running.Swap(true) {
		return
	}

	go func() {
		for {
			select {
			case cmd := <-gs.GameCommand:
				switch c := cmd.(type) {
				case AbortGameCmd:
					err := gs.abortGame(c)
					if err != nil {
						fmt.Println("------------------- abortGame err: ", err)
					}
				case ResignGameCmd:
					err := gs.resignGame(c)
					if err != nil {
						fmt.Println("------------------- resignGame err: ", err)
					}
				case OfferDrawCmd:
					err := gs.offerDraw(c)
					if err != nil {
						fmt.Println("------------------- offerDraw err: ", err)
					}
				case AcceptDrawCmd:
					err := gs.acceptDraw(c)
					if err != nil {
						fmt.Println("------------------- acceptDraw err: ", err)
					}
				case DeclineDrawCmd:
					err := gs.declineDraw(c)
					if err != nil {
						fmt.Println("------------------- declineDraw err: ", err)
					}
				case PlayMoveUCICmd:
					res, err := gs.playMoveUCI(c)
					if err != nil {
						fmt.Println("------------------- playMoveUCI err: ", err)

						gs.GameEvent <- PlayMoveUCIErrorEvent{Err: err}
					}

					gs.GameEvent <- res
				}

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (gs *GameState) abortGame(c AbortGameCmd) error {
	godump.Dump("-------------------------- abortGame: ", c)

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return ErrGameAlreadyConcluded
	}

	return nil
}

func (gs *GameState) resignGame(c ResignGameCmd) error {
	godump.Dump("-------------------------- resignGame: ", c)

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return ErrGameAlreadyConcluded
	}

	return nil
}

func (gs *GameState) offerDraw(c OfferDrawCmd) error {
	godump.Dump("-------------------------- offerDraw: ", c)

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return ErrGameAlreadyConcluded
	}

	return nil
}

func (gs *GameState) acceptDraw(c AcceptDrawCmd) error {
	godump.Dump("-------------------------- acceptDraw: ", c)

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return ErrGameAlreadyConcluded
	}

	return nil
}

func (gs *GameState) declineDraw(c DeclineDrawCmd) error {
	godump.Dump("-------------------------- declineDraw: ", c)

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return ErrGameAlreadyConcluded
	}

	return nil
}

func (gs *GameState) playMoveUCI(c PlayMoveUCICmd) (PlayMoveUCIEvent, error) {
	godump.Dump("-------------------------- playMoveUCI: ", c)

	if !gs.HasGamePlayer(c.UserID) {
		return PlayMoveUCIEvent{}, ErrPlayerNotInGame
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return PlayMoveUCIEvent{}, ErrGameAlreadyConcluded
	}

	player := gs.GetPlayerByID(c.UserID)
	if player.Color == pb.Color_COLOR_WHITE && gs.Chess.Position.Turn.IsBlack() || player.Color == pb.Color_COLOR_BLACK && gs.Chess.Position.Turn.IsWhite() {
		return PlayMoveUCIEvent{}, ErrNotYourTurn
	}

	move, err := gs.Chess.MakeMoveUCI(c.UCI)
	if err != nil {
		return PlayMoveUCIEvent{}, fmt.Errorf("%w: %w", ErrInvalidMove, err)
	}

	playedAt := time.Now()

	uci := move.ToUCI()
	lan := move.ToLAN(gs.Chess.Position, gs.Chess.Position.Check, gs.Chess.IsCheckmate())
	san := move.ToSAN(gs.Chess.Position, gs.Chess.Position.Check, gs.Chess.IsCheckmate(), gs.Chess.LegalMoves)
	fen := gs.Chess.Position.Fen()

	gs.LastMove = new(playedAt)
	gs.GameMoves = append(gs.GameMoves, &pb.GameMove{
		Fen:      fen,
		Uci:      new(uci),
		San:      new(san),
		Check:    gs.Chess.Position.Check,
		PlayedAt: timestamppb.New(playedAt),
	})

	if gs.Chess.Position.Ply == 1 {
		// stop white first move timer
		// start black first move timer
	}

	if gs.Chess.Position.Ply == 2 {
		// stop black first move timer
	}

	if gs.Chess.Position.Ply >= 2 {
		// gs.toggleClockAfterMove()
	}

	legalMoves := make([]string, len(gs.Chess.LegalMoves))
	for i, legalMove := range gs.Chess.LegalMoves {
		legalMoves[i] = fmt.Sprint(legalMove.String())
	}

	if gs.Chess.IsTerminated() {
		switch gs.Chess.Status() {
		case engine.StatusInsufficientMaterial:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_INSUFFICIENT_MATERIAL
		case engine.StatusThreeFoldRepetition:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_THREEFOLD_REPETITION
		case engine.StatusFiveFoldRepetition:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_FIVEFOLD_REPETITION
		case engine.StatusFiftyMoveRule:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_FIFTY_MOVE_RULE
		case engine.StatusSeventyFiveMoveRule:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_SEVENTYFIVE_MOVE_RULE
		case engine.StatusCheckmate:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			if player.Color == pb.Color_COLOR_WHITE {
				gs.GameResult = pb.GameResult_GAME_RESULT_WHITE_WON
			} else {
				gs.GameResult = pb.GameResult_GAME_RESULT_BLACK_WON
			}

			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE
		case engine.StatusStalemate:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE
		}
	}

	res := PlayMoveUCIEvent{
		GameID:                gs.GameID,
		UserID:                c.UserID,
		Players:               gs.Players,
		Uci:                   uci,
		San:                   san,
		Lan:                   lan,
		WhiteClockRemainingMs: 420_000,
		BlackClockRemainingMs: 69_000,
		Position:              gs.Chess.Position.Copy(),
		GameResult:            gs.GameResult,
		GameResultStatus:      gs.GameResultStatus,
		GameState:             gs.GameState,
		LastMove:              gs.LastMove,
		StartTime:             gs.StartTime,
		EndTime:               gs.EndTime,
		Repetitions:           gs.Chess.Repetitions,
		LegalMoves:            legalMoves,
	}

	return res, nil
}

func (gs *GameState) Stop() {
	if gs.activeGameTimer != nil {
		gs.activeGameTimer.Stop()
	}

	if gs.firstMoveTimer != nil {
		gs.firstMoveTimer.Stop()
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
			return ErrPlayerIdMissing
		}

		if p.Color == pb.Color_COLOR_UNSPECIFIED {
			return ErrPlayerColorMissing
		}
	}

	if players[0].Color == players[1].Color {
		return ErrPlayersSameColors
	}

	return nil
}

func determineGameTimeCategoryFromTimeControl(gtc *pb.GameTimeControl, thresholds []CategoryThreshold) (pb.GameTimeCategory, error) {
	if gtc == nil {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, ErrTimeControlMissing
	}

	if gtc.GetClockMs() <= 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, ErrTimeControlClockInvalid
	}

	if gtc.GetIncrementMs() < 0 {
		return pb.GameTimeCategory_GAME_TIME_CATEGORY_UNSPECIFIED, ErrTimeControlIncrementInvalid
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
