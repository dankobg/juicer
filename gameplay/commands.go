package gameplay

import "github.com/google/uuid"

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
}

func (OfferDrawCmd) isGameCommand() {}

type AcceptDrawCmd struct {
	GameID int64
	UserID uuid.UUID
}

func (AcceptDrawCmd) isGameCommand() {}

type DeclineDrawCmd struct {
	GameID int64
	UserID uuid.UUID
}

func (DeclineDrawCmd) isGameCommand() {}
