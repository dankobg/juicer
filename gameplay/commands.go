package gameplay

import "github.com/google/uuid"

type GameCommand interface {
	isGameCommand()
}

type PlayMoveUCICmd struct {
	UserID uuid.UUID
	UCI    string
}

func (PlayMoveUCICmd) isGameCommand() {}

type AbortGameCmd struct {
	UserID uuid.UUID
}

func (AbortGameCmd) isGameCommand() {}

type ResignGameCmd struct {
	UserID uuid.UUID
}

func (ResignGameCmd) isGameCommand() {}

type OfferDrawCmd struct {
	UserID uuid.UUID
}

func (OfferDrawCmd) isGameCommand() {}

type AcceptDrawCmd struct {
	UserID uuid.UUID
}

func (AcceptDrawCmd) isGameCommand() {}

type DeclineDrawCmd struct {
	UserID uuid.UUID
}

func (DeclineDrawCmd) isGameCommand() {}
