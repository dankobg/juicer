package engine

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
	StatusInsufficientMaterial
)

type History struct {
	move Move
	pos  Position
}

type Chess struct {
	Position             *Position
	Repetitions          uint16
	History              []History
	HistoryHashes        []uint64
	LegalMoves           []Move
	DisableAutoThreefold bool
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
		Position:      p,
		HistoryHashes: []uint64{p.Hash},
	}

	c.calcLegalMoves()

	return c, nil
}

func (c *Chess) calcLegalMoves() {
	pseudo := c.Position.generateAllPseudoLegalMoves()
	legal := c.Position.generateAllLegalMoves(pseudo)
	c.LegalMoves = legal
}

func (c *Chess) AppendHistoryEntry(m Move, pos Position) {
	h := History{move: m, pos: pos}
	c.History = append(c.History, h)
}

func (c *Chess) AppendHistoryHash(hash uint64) {
	c.HistoryHashes = append(c.HistoryHashes, hash)
}

func (c *Chess) MakeMove(m Move) {
	c.Position.MakeMove(m)

	pos := c.Position.Copy()

	c.AppendHistoryEntry(m, *pos)
	c.AppendHistoryHash(pos.Hash)
	c.calcRepetitions()
	c.calcLegalMoves()
}

func (c *Chess) MakeMoveUCI(uciMove string) (Move, error) {
	for _, m := range c.LegalMoves {
		if m.String() == uciMove {
			c.MakeMove(m)
			return m, nil
		}
	}
	return 0, fmt.Errorf("invalid move: %v", uciMove)
}

func (c *Chess) calcRepetitions() {
	var reps uint16
	depth := max(0, c.Position.Ply-uint16(c.Position.HalfMoveClock))

	for i := len(c.HistoryHashes) - 1; i >= int(depth); i -= 2 {
		if c.HistoryHashes[i] == c.Position.Hash {
			reps++
		}
	}

	c.Repetitions = reps
}

func (c *Chess) IsInsufficientMaterial() bool {
	return c.Position.IsInsufficientMaterial()
}

func (c *Chess) IsThreefoldRepetition() bool {
	return !c.DisableAutoThreefold && c.Repetitions >= 3
}

func (c *Chess) IsFivefoldRepetition() bool {
	return !c.DisableAutoThreefold && c.Repetitions >= 5
}

func (c *Chess) IsDrawBy50MoveRule() bool {
	return c.Position.HalfMoveClock >= 100
}

func (c *Chess) IsDrawBy75MoveRule() bool {
	return c.Position.HalfMoveClock >= 150
}

func (c *Chess) IsDraw() bool {
	return c.IsDrawBy50MoveRule() || c.IsThreefoldRepetition() || c.IsStalemate() || c.IsInsufficientMaterial()
}

func (c *Chess) IsCheckmate() bool {
	return len(c.LegalMoves) == 0 && c.Position.Board.IsInCheck(c.Position.Turn)
}

func (c *Chess) IsStalemate() bool {
	return len(c.LegalMoves) == 0 && !c.Position.Board.IsInCheck(c.Position.Turn)
}

func (c *Chess) IsTerminated() bool {
	return c.IsDraw() || c.IsCheckmate()
}

func (c *Chess) Status() Status {
	if c.IsInsufficientMaterial() {
		return StatusInsufficientMaterial
	}
	if c.IsThreefoldRepetition() {
		return StatusThreeFoldRepetition
	}
	if c.IsFivefoldRepetition() {
		return StatusFiveFoldRepetition
	}
	if c.IsDrawBy50MoveRule() {
		return StatusFiftyMoveRule
	}
	if c.IsDrawBy75MoveRule() {
		return StatusSeventyFiveMoveRule
	}
	if c.IsCheckmate() {
		return StatusCheckmate
	}
	if c.IsStalemate() {
		return StatusStalemate
	}

	return StatusUnknown
}
