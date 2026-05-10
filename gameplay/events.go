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
	GameID                int64
	UserID                uuid.UUID
	Players               map[uuid.UUID]*Player
	Uci                   string
	San                   string
	Lan                   string
	WhiteClockRemainingMs int64
	BlackClockRemainingMs int64
	Position              *engine.Position
	GameResult            pb.GameResult
	GameResultStatus      pb.GameResultStatus
	GameState             pb.GameState
	LastMove              *time.Time
	StartTime             *time.Time
	EndTime               *time.Time
	Repetitions           uint16
	LegalMoves            []string
	Version               int
	// History []engine.History
	// Hashes []engine.HistoryHash
}

func (PlayMoveUCIEvent) isGameEvent() {}

type PlayMoveUCIErrorEvent struct {
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
	UserID uuid.UUID
	Err    error
}

func (ResignErrorEvent) isGameEvent() {}

type OfferDrawEvent struct{}

func (OfferDrawEvent) isGameEvent() {}

type OfferDrawErrorEvent struct {
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
	UserID uuid.UUID
	Err    error
}

func (AcceptDrawErrorEvent) isGameEvent() {}

type DeclineDrawEvent struct{}

func (DeclineDrawEvent) isGameEvent() {}

type DeclineDrawErrorEvent struct {
	UserID uuid.UUID
	Err    error
}

func (DeclineDrawErrorEvent) isGameEvent() {}
