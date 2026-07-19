package gameplay

import (
	"time"

	"github.com/google/uuid"
)

type GameCommand interface {
	isGameCommand()
}

type PlayMoveUCICmd struct {
	GameID int64
	UserID uuid.UUID
	UCI    string
	Ack    int32
}

func (PlayMoveUCICmd) isGameCommand() {}

type AbortGameCmd struct {
	GameID int64
	UserID uuid.UUID
}

func (AbortGameCmd) isGameCommand() {}

type ResignGameCmd struct {
	GameID int64
	UserID uuid.UUID
}

func (ResignGameCmd) isGameCommand() {}

type OfferDrawCmd struct {
	GameID int64
	UserID uuid.UUID
	Ply    int
}

func (OfferDrawCmd) isGameCommand() {}

type AcceptDrawCmd struct {
	GameID int64
	UserID uuid.UUID
	Ply    int
}

func (AcceptDrawCmd) isGameCommand() {}

type DeclineDrawCmd struct {
	GameID int64
	UserID uuid.UUID
}

func (DeclineDrawCmd) isGameCommand() {}

type RejoinedGame struct {
	GameID     int64
	UserID     uuid.UUID
	RejoinedAt time.Time
}

func (RejoinedGame) isGameCommand() {}

type LeftGame struct {
	GameID int64
	UserID uuid.UUID
	LefAt  time.Time
}

func (LeftGame) isGameCommand() {}
