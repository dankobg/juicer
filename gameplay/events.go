package gameplay

import (
	"time"

	"github.com/dankobg/juicer/engine"
	pb "github.com/dankobg/juicer/pb/proto/juicer"
	"github.com/google/uuid"
)

type GameEvent interface {
	isGameEvent()
}

type PlayMoveUCIEvent struct {
	GameID                 int64
	UserID                 uuid.UUID
	Players                map[uuid.UUID]*Player
	Uci                    string
	San                    string
	Lan                    string
	WhiteRemainingGameTime time.Duration
	BlackRemainingGameTime time.Duration
	Position               *engine.Position
	LastMove               *time.Time
	StartTime              *time.Time
	Repetitions            uint16
	LegalMoves             []string
	Version                int
	PlayedAt               time.Time
	// History []engine.History
	// Hashes []engine.HistoryHash
}

func (PlayMoveUCIEvent) isGameEvent() {}

type PlayMoveUCIErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (PlayMoveUCIErrorEvent) isGameEvent() {}

type AbortEvent struct {
	GameID           int64
	GameResult       pb.GameResult
	GameResultStatus pb.GameResultStatus
	GameState        pb.GameState
	EndTime          time.Time
}

func (AbortEvent) isGameEvent() {}

type AbortErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (AbortErrorEvent) isGameEvent() {}

type ResignEvent struct {
	GameID           int64
	GameResult       pb.GameResult
	GameResultStatus pb.GameResultStatus
	GameState        pb.GameState
	EndTime          time.Time
}

func (ResignEvent) isGameEvent() {}

type ResignErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (ResignErrorEvent) isGameEvent() {}

type OfferDrawEvent struct {
	GameID int64
	UserID uuid.UUID
}

func (OfferDrawEvent) isGameEvent() {}

type OfferDrawErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (OfferDrawErrorEvent) isGameEvent() {}

type AcceptDrawEvent struct {
	GameID           int64
	GameResult       pb.GameResult
	GameResultStatus pb.GameResultStatus
	GameState        pb.GameState
	EndTime          time.Time
}

func (AcceptDrawEvent) isGameEvent() {}

type AcceptDrawErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (AcceptDrawErrorEvent) isGameEvent() {}

type DeclineDrawEvent struct {
	GameID int64
	UserID uuid.UUID
}

func (DeclineDrawEvent) isGameEvent() {}

type DeclineDrawErrorEvent struct {
	GameID int64
	UserID uuid.UUID
	Err    error
}

func (DeclineDrawErrorEvent) isGameEvent() {}

type GameFinishedEvent struct {
	GameID           int64
	GameResult       pb.GameResult
	GameResultStatus pb.GameResultStatus
	GameState        pb.GameState
	EndTime          time.Time
}

func (GameFinishedEvent) isGameEvent() {}
