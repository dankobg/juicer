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
	ErrPlayersSameColors  = errors.New("players can't have the same color")

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

type DrawOffer struct {
	OfferedBy uuid.UUID
	OfferedAt time.Time
	Ply       int
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
	GameTimeControl        *pb.GameTimeControl
	GameResult             pb.GameResult
	GameResultStatus       pb.GameResultStatus
	GameState              pb.GameState
	ReconnectTimeout       time.Duration
	FirstMoveTimeout       time.Duration
	LastMove               *time.Time
	StartTime              *time.Time
	EndTime                *time.Time
	Rated                  bool
	GameMoves              []*pb.GameMove
	Version                int
	PendingDrawOffer       *DrawOffer
	running                atomic.Bool
	GameCommand            chan GameCommand
	GameEvent              chan GameEvent
	WhiteRemainingGameTime time.Duration
	BlackRemainingGameTime time.Duration
	WhiteDisconnectedAt    *time.Time
	BlackDisconnectedAt    *time.Time
	activeGameTimer        *time.Timer
	firstMoveTimer         *time.Timer
	whiteReconnectTimer    *time.Timer
	blackReconnectTimer    *time.Timer
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
		Chess:                  chess,
		GameID:                 gopts.gameID,
		White:                  white,
		Black:                  black,
		Players:                playersByID,
		GameVariant:            gopts.gameVariant,
		GameTimeKind:           gopts.gameTimeKind,
		GameTimeCategory:       gopts.gameTimeCategory,
		GameTimeControl:        gopts.gameTimeControl,
		WhiteRemainingGameTime: time.Duration(gopts.gameTimeControl.ClockMs) * time.Millisecond,
		BlackRemainingGameTime: time.Duration(gopts.gameTimeControl.ClockMs) * time.Millisecond,
		GameState:              gopts.gameState,
		GameResult:             gopts.gameResult,
		GameResultStatus:       gopts.gameResultStatus,
		FirstMoveTimeout:       gopts.firstMoveTimeout,
		ReconnectTimeout:       gopts.reconnectTimeout,
		Rated:                  gopts.rated,
		StartTime:              gopts.startTime,
		LastMove:               gopts.lastMove,
		EndTime:                gopts.endTime,
		GameMoves:              gameMoves,
		running:                atomic.Bool{},
		GameCommand:            make(chan GameCommand, 64),
		GameEvent:              gameEvent,
	}

	return gs, nil
}

func (gs *GameState) GetPlayerByID(id uuid.UUID) *Player {
	return gs.Players[id]
}

func (gs *GameState) GetOtherPlayer(id uuid.UUID) *Player {
	for playerID, player := range gs.Players {
		if playerID != id {
			return player
		}
	}

	return nil
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

func Tick(t *time.Timer) <-chan time.Time {
	if t != nil {
		return t.C
	}

	return nil
}

func (gs *GameState) Start(ctx context.Context) {
	if gs.running.Swap(true) {
		return
	}

	gs.firstMoveTimer = time.NewTimer(gs.FirstMoveTimeout)

	go func() {
		for {
			select {
			case <-Tick(gs.firstMoveTimer):
				gs.GameEvent <- GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_INTERRUPTED,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT,
					GameState:        pb.GameState_GAME_STATE_INTERRUPTED,
					EndTime:          time.Now(),
				}

			case <-Tick(gs.activeGameTimer):
				gameResult := pb.GameResult_GAME_RESULT_BLACK_WON
				if gs.Chess.Position.Turn.IsBlack() {
					gameResult = pb.GameResult_GAME_RESULT_WHITE_WON
				}

				gs.GameEvent <- GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       gameResult,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED,
					GameState:        pb.GameState_GAME_STATE_FINISHED,
					EndTime:          time.Now(),
				}

			case <-Tick(gs.whiteReconnectTimer):
				gs.GameEvent <- GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_BLACK_WON,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT,
					GameState:        pb.GameState_GAME_STATE_FINISHED,
					EndTime:          time.Now(),
				}

			case <-Tick(gs.blackReconnectTimer):
				gs.GameEvent <- GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_WHITE_WON,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_TIMED_OUT,
					GameState:        pb.GameState_GAME_STATE_FINISHED,
					EndTime:          time.Now(),
				}

			case cmd := <-gs.GameCommand:
				switch c := cmd.(type) {
				case AbortGameCmd:
					if events, err := gs.abortGame(c); err != nil {
						gs.GameEvent <- AbortErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case ResignGameCmd:
					if events, err := gs.resignGame(c); err != nil {
						gs.GameEvent <- ResignErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case OfferDrawCmd:
					if events, err := gs.offerDraw(c); err != nil {
						gs.GameEvent <- OfferDrawErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case AcceptDrawCmd:
					if events, err := gs.acceptDraw(c); err != nil {
						gs.GameEvent <- AcceptDrawErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case DeclineDrawCmd:
					if events, err := gs.declineDraw(c); err != nil {
						gs.GameEvent <- DeclineDrawErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case PlayMoveUCICmd:
					if events, err := gs.playMoveUCI(c); err != nil {
						gs.GameEvent <- PlayMoveUCIErrorEvent{
							GameID: c.GameID,
							UserID: c.UserID,
							Err:    err,
						}
					} else {
						for _, event := range events {
							gs.GameEvent <- event
						}
					}

				case RejoinedGame:
					_, _ = gs.rejoinedGame(c)

				case LeftGame:
					_, _ = gs.leftGame(c)
				}

			case <-ctx.Done():
				gs.GameEvent <- GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_INTERRUPTED,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_INTERRUPTED,
					GameState:        pb.GameState_GAME_STATE_INTERRUPTED,
					EndTime:          time.Now(),
				}

				return
			}
		}
	}()
}

func (gs *GameState) rejoinedGame(c RejoinedGame) ([]GameEvent, error) {
	player := gs.GetPlayerByID(c.UserID)

	if !gs.HasGamePlayer(c.UserID) || player == nil {
		return nil, ErrPlayerNotInGame
	}

	if player.Color == pb.Color_COLOR_WHITE {
		gs.WhiteDisconnectedAt = nil
		if gs.whiteReconnectTimer != nil {
			gs.whiteReconnectTimer.Stop()
		}
		gs.whiteReconnectTimer = nil
	} else {
		gs.BlackDisconnectedAt = nil
		if gs.blackReconnectTimer != nil {
			gs.blackReconnectTimer.Stop()
		}
		gs.blackReconnectTimer = nil
	}

	return nil, nil
}

func (gs *GameState) leftGame(c LeftGame) ([]GameEvent, error) {
	player := gs.GetPlayerByID(c.UserID)

	if !gs.HasGamePlayer(c.UserID) || player == nil {
		return nil, ErrPlayerNotInGame
	}

	if player.Color == pb.Color_COLOR_WHITE {
		gs.WhiteDisconnectedAt = &c.LefAt
		gs.whiteReconnectTimer = time.NewTimer(gs.ReconnectTimeout)
	} else {
		gs.BlackDisconnectedAt = &c.LefAt
		gs.blackReconnectTimer = time.NewTimer(gs.ReconnectTimeout)
	}

	return nil, nil
}

func (gs *GameState) abortGame(c AbortGameCmd) ([]GameEvent, error) {
	abortedAt := time.Now()

	if !gs.HasGamePlayer(c.UserID) {
		return nil, ErrPlayerNotInGame
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	player := gs.GetPlayerByID(c.UserID)

	if player.Color == pb.Color_COLOR_WHITE {
		if gs.Chess.Position.Ply <= 1 {
			gs.GameResult = pb.GameResult_GAME_RESULT_INTERRUPTED
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			gs.GameResult = pb.GameResult_GAME_RESULT_BLACK_WON
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
		}
	} else {
		if gs.Chess.Position.Ply <= 2 {
			gs.GameResult = pb.GameResult_GAME_RESULT_INTERRUPTED
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			gs.GameResult = pb.GameResult_GAME_RESULT_WHITE_WON
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
		}
	}

	gs.EndTime = &abortedAt
	gs.Stop()

	events := []GameEvent{
		AbortEvent{
			GameID:           gs.GameID,
			GameResult:       gs.GameResult,
			GameResultStatus: gs.GameResultStatus,
			GameState:        gs.GameState,
			EndTime:          *gs.EndTime,
		},
	}

	return events, nil
}

func (gs *GameState) resignGame(c ResignGameCmd) ([]GameEvent, error) {
	resignedAt := time.Now()

	if !gs.HasGamePlayer(c.UserID) {
		return nil, ErrPlayerNotInGame
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	player := gs.GetPlayerByID(c.UserID)

	if player.Color == pb.Color_COLOR_WHITE {
		if gs.Chess.Position.Ply <= 1 {
			gs.GameResult = pb.GameResult_GAME_RESULT_INTERRUPTED
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			gs.GameResult = pb.GameResult_GAME_RESULT_BLACK_WON
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
		}
	}

	if player.Color == pb.Color_COLOR_BLACK {
		if gs.Chess.Position.Ply < 2 {
			gs.GameResult = pb.GameResult_GAME_RESULT_INTERRUPTED
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_ABORTED
			gs.GameState = pb.GameState_GAME_STATE_INTERRUPTED
		} else {
			gs.GameResult = pb.GameResult_GAME_RESULT_WHITE_WON
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_RESIGNATION
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
		}
	}

	gs.EndTime = &resignedAt
	gs.Stop()

	events := []GameEvent{
		ResignEvent{
			GameID:           gs.GameID,
			GameResult:       gs.GameResult,
			GameResultStatus: gs.GameResultStatus,
			GameState:        gs.GameState,
			EndTime:          *gs.EndTime,
		},
	}

	return events, nil
}

func (gs *GameState) offerDraw(c OfferDrawCmd) ([]GameEvent, error) {
	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	res := []GameEvent{OfferDrawEvent{}}

	return res, nil
}

func (gs *GameState) acceptDraw(c AcceptDrawCmd) ([]GameEvent, error) {
	acceptedDrawAt := time.Now()

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	if !gs.HasGamePlayer(c.UserID) {
		return nil, ErrPlayerNotInGame
	}

	gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
	gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_DRAW_AGREED
	gs.GameState = pb.GameState_GAME_STATE_FINISHED

	gs.EndTime = &acceptedDrawAt

	events := []GameEvent{
		AcceptDrawEvent{
			GameID:           gs.GameID,
			GameResult:       gs.GameResult,
			GameResultStatus: gs.GameResultStatus,
			GameState:        gs.GameState,
			EndTime:          *gs.EndTime,
		},
	}

	return events, nil
}

func (gs *GameState) declineDraw(c DeclineDrawCmd) ([]GameEvent, error) {
	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	events := []GameEvent{DeclineDrawEvent{}}

	return events, nil
}

func (gs *GameState) playMoveUCI(c PlayMoveUCICmd) ([]GameEvent, error) {
	playedAt := time.Now()

	if !gs.HasGamePlayer(c.UserID) {
		return nil, ErrPlayerNotInGame
	}

	if gs.GameResult != pb.GameResult_GAME_RESULT_UNSPECIFIED {
		return nil, ErrGameAlreadyConcluded
	}

	player := gs.GetPlayerByID(c.UserID)
	if player.Color == pb.Color_COLOR_WHITE && gs.Chess.Position.Turn.IsBlack() || player.Color == pb.Color_COLOR_BLACK && gs.Chess.Position.Turn.IsWhite() {
		return nil, ErrNotYourTurn
	}

	var terminated bool

	events := []GameEvent{}

	if gs.LastMove != nil && gs.Chess.Position.Ply >= 2 {
		elapsed := playedAt.Sub(*gs.LastMove)
		increment := time.Duration(gs.GameTimeControl.GetIncrementMs()) * time.Millisecond

		if gs.Chess.Position.Turn.IsWhite() {
			previousRemaining := gs.WhiteRemainingGameTime
			gs.WhiteRemainingGameTime -= elapsed

			if gs.WhiteRemainingGameTime <= 0 {
				flaggedAt := gs.LastMove.Add(previousRemaining)
				events = append(events, GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_BLACK_WON,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED,
					GameState:        pb.GameState_GAME_STATE_FINISHED,
					EndTime:          flaggedAt,
				})
				terminated = true
			} else {
				gs.WhiteRemainingGameTime += increment
			}
		} else {
			previousRemaining := gs.BlackRemainingGameTime
			gs.BlackRemainingGameTime -= elapsed

			if gs.BlackRemainingGameTime <= 0 {
				flaggedAt := gs.LastMove.Add(previousRemaining)
				events = append(events, GameFinishedEvent{
					GameID:           gs.GameID,
					GameResult:       pb.GameResult_GAME_RESULT_WHITE_WON,
					GameResultStatus: pb.GameResultStatus_GAME_RESULT_STATUS_FLAGGED,
					GameState:        pb.GameState_GAME_STATE_FINISHED,
					EndTime:          flaggedAt,
				})
				terminated = true
			} else {
				gs.BlackRemainingGameTime += increment
			}
		}
	}

	if terminated {
		return events, nil
	}

	move, err := gs.Chess.MakeMoveUCI(c.UCI)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidMove, err)
	}

	uci := move.ToUCI()
	lan := move.ToLAN(gs.Chess.Position, gs.Chess.Position.Check, gs.Chess.IsCheckmate())
	san := move.ToSAN(gs.Chess.Position, gs.Chess.Position.Check, gs.Chess.IsCheckmate(), gs.Chess.LegalMoves)
	fen := gs.Chess.Position.Fen()

	gs.LastMove = &playedAt
	gs.GameMoves = append(gs.GameMoves, &pb.GameMove{
		Fen:      fen,
		Uci:      new(uci),
		San:      new(san),
		Lan:      new(lan),
		Check:    gs.Chess.Position.Check,
		PlayedAt: timestamppb.New(playedAt),
	})
	gs.Version = int(c.Ack)

	legalMoves := make([]string, len(gs.Chess.LegalMoves))
	for i, legalMove := range gs.Chess.LegalMoves {
		legalMoves[i] = fmt.Sprint(legalMove.String())
	}

	if gs.Chess.Position.Ply == 1 {
		// stop white first move timer
		// start black first move timer
		if gs.firstMoveTimer == nil {
			gs.firstMoveTimer = time.NewTimer(gs.FirstMoveTimeout)
		} else {
			gs.firstMoveTimer.Reset(gs.FirstMoveTimeout)
		}
	}

	if gs.Chess.Position.Ply == 2 {
		// stop black first move timer
		if gs.firstMoveTimer != nil {
			gs.firstMoveTimer.Stop()
		}

		gs.firstMoveTimer = nil
	}

	if gs.Chess.Position.Ply >= 2 {
		remaining := gs.WhiteRemainingGameTime
		if gs.Chess.Position.Turn.IsBlack() {
			remaining = gs.BlackRemainingGameTime
		}

		if gs.activeGameTimer == nil {
			gs.activeGameTimer = time.NewTimer(remaining)
		} else {
			gs.activeGameTimer.Reset(remaining)
		}
	}

	playMoveUciEvent := PlayMoveUCIEvent{
		GameID:                 gs.GameID,
		UserID:                 c.UserID,
		Players:                gs.Players,
		Uci:                    uci,
		San:                    san,
		Lan:                    lan,
		WhiteRemainingGameTime: gs.WhiteRemainingGameTime,
		BlackRemainingGameTime: gs.BlackRemainingGameTime,
		Position:               gs.Chess.Position.Copy(),
		LastMove:               gs.LastMove,
		StartTime:              gs.StartTime,
		Repetitions:            gs.Chess.Repetitions,
		LegalMoves:             legalMoves,
		Version:                gs.Version,
		PlayedAt:               playedAt,
	}

	events = append(events, playMoveUciEvent)

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
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_CHECKMATE

			if player.Color == pb.Color_COLOR_WHITE {
				gs.GameResult = pb.GameResult_GAME_RESULT_WHITE_WON
			} else {
				gs.GameResult = pb.GameResult_GAME_RESULT_BLACK_WON
			}
		case engine.StatusStalemate:
			gs.GameState = pb.GameState_GAME_STATE_FINISHED
			gs.GameResult = pb.GameResult_GAME_RESULT_DRAW
			gs.GameResultStatus = pb.GameResultStatus_GAME_RESULT_STATUS_STALEMATE
		}

		events = append(events, GameFinishedEvent{
			GameID:           gs.GameID,
			GameResult:       gs.GameResult,
			GameResultStatus: gs.GameResultStatus,
			GameState:        gs.GameState,
			EndTime:          time.Now(),
		})
	}

	debug_print_gamestate_info(gs)

	return events, nil
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

func (gs *GameState) whiteReconnectTimeoutExpired() bool {
	if gs.WhiteDisconnectedAt == nil {
		return false
	}

	elapsed := time.Since(*gs.WhiteDisconnectedAt)

	return elapsed >= gs.ReconnectTimeout
}

func (gs *GameState) blackReconnectTimeoutExpired() bool {
	if gs.blackReconnectTimer == nil {
		return false
	}

	elapsed := time.Since(*gs.BlackDisconnectedAt)

	return elapsed >= gs.ReconnectTimeout
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

func debug_print_gamestate_info(gs *GameState) {
	fmt.Printf("game_id: %d\n", gs.GameID)
	fmt.Printf("rated: %v\n", gs.Rated)
	fmt.Printf("white: %s\n", gs.White.Username)
	fmt.Printf("black: %s\n", gs.Black.Username)
	fmt.Printf("variant: %s\n", gs.GameVariant.String())
	fmt.Printf("time_category: %s\n", gs.GameTimeCategory.String())
	fmt.Printf("time_kind: %s\n", gs.GameTimeKind.String())
	fmt.Printf("clock_ms: %d\n", gs.GameTimeControl.GetClockMs())
	fmt.Printf("increment_ms: %d\n", gs.GameTimeControl.GetIncrementMs())
	fmt.Printf("state: %s\n", gs.GameState.String())
	fmt.Printf("result: %s\n", gs.GameResult.String())
	fmt.Printf("result_status: %s\n", gs.GameResultStatus.String())
	fmt.Printf("start_time: %s\n", gs.StartTime)
	fmt.Printf("last_move: %s\n", gs.LastMove)
	fmt.Printf("game_moves: %v\n", gs.GameMoves)
	fmt.Printf("repetitions: %v\n", gs.Chess.Repetitions)
	fmt.Printf("history_hashes: %v\n", gs.Chess.HistoryHashes)
	fmt.Printf("version: %d\n", gs.Version)
	fmt.Printf("white_remaining: %v\n", gs.WhiteRemainingGameTime)
	fmt.Printf("black_remaining: %v\n", gs.BlackRemainingGameTime)
	fmt.Println(gs.Chess.Position.PrintBoard())

	legals := []string{}
	for _, x := range gs.Chess.LegalMoves {
		legals = append(legals, x.String())
	}

	godump.DumpJSON("legal moves", legals)
}
