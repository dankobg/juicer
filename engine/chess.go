package juicer

import (
	"fmt"
)

const (
	FENEmptyPosition    = "8/8/8/8/8/8/8/8 w KQkq - 0 1"
	FENStartingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

type Outcome uint8

const (
	OutcomeUnknown Outcome = iota
	OutcomeWhiteWon
	OutcomeBlackWon
	OutcomeDraw
)

type Status uint8

const (
	StatusUnknown Status = iota
	StatusCheckmate
	StatusStalemate
	StatusResignation
	StatusDrawOfferAccepted
	StatusThreeFoldRepetition
	StatusFiveFoldRepetition
	StatusFiftyMoveRule
	StatusSeventyFiveMoveRule
	InsufficientMaterial
)

type History struct {
	move Move
	pos  Position
}

type Chess struct {
	position             *Position
	repetitions          uint16
	history              []History
	historyHashes        []uint64
	legalMoves           []Move
	disableAutoThreefold bool
}

func InitPrecalculatedTables() {
	initAllAttackMasksTables()
	initZobrist()
}

func NewChess(fen string) (*Chess, error) {
	p := &Position{}

	if err := p.LoadFromFEN(fen); err != nil {
		return nil, fmt.Errorf("failed to start new game: %w", err)
	}

	c := &Chess{
		position:      p,
		historyHashes: []uint64{p.hash},
	}

	c.calcLegalMoves()

	return c, nil
}

func (c *Chess) calcLegalMoves() {
	pseudo := c.position.generateAllPseudoLegalMoves()
	legal := c.position.generateAllLegalMoves(pseudo)
	c.legalMoves = legal
}

func (c *Chess) AppendHistoryEntry(m Move, pos Position) {
	h := History{move: m, pos: pos}
	c.history = append(c.history, h)
}

func (c *Chess) AppendHistoryHash(hash uint64) {
	c.historyHashes = append(c.historyHashes, hash)
}

func (c *Chess) MakeMove(m Move) {
	c.position.MakeMove(m)

	pos := c.position.Copy()

	c.AppendHistoryEntry(m, *pos)
	c.AppendHistoryHash(pos.hash)
	c.calcRepetitions()
	c.calcLegalMoves()
}

func (c *Chess) calcRepetitions() {
	var reps uint16
	depth := max(0, c.position.ply-uint16(c.position.halfMoveClock))

	for i := len(c.historyHashes) - 1; i >= int(depth); i -= 2 {
		if c.historyHashes[i] == c.position.hash {
			reps++
		}
	}

	c.repetitions = reps
}

func (c *Chess) IsInsufficientMaterial() bool {
	return c.position.IsInsufficientMaterial()
}

func (c *Chess) IsThreefoldRepetition() bool {
	return !c.disableAutoThreefold && c.repetitions >= 3
}

func (c *Chess) IsFivefoldRepetition() bool {
	return !c.disableAutoThreefold && c.repetitions >= 5
}

func (c *Chess) IsDrawBy50MoveRule() bool {
	return c.position.halfMoveClock >= 100
}

func (c *Chess) IsDrawBy75MoveRule() bool {
	return c.position.halfMoveClock >= 150
}

func (c *Chess) IsDraw() bool {
	return c.IsDrawBy50MoveRule() || c.IsThreefoldRepetition() || c.IsStalemate() || c.IsInsufficientMaterial()
}

func (c *Chess) IsCheckmate() bool {
	return len(c.legalMoves) == 0 && c.position.board.IsInCheck(c.position.turn)
}

func (c *Chess) IsStalemate() bool {
	return len(c.legalMoves) == 0 && !c.position.board.IsInCheck(c.position.turn)
}

func (c *Chess) IsTerminated() bool {
	return c.IsDraw() || c.IsCheckmate()
}
