package juicer

import (
	"fmt"
)

const (
	FENEmptyPosition    = "8/8/8/8/8/8/8/8 w KQkq - 0 1"
	FENStartingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	FenTestPosition1 = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	FenTestPosition3 = "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"
	FenTestPosition4 = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	FenTestPosition5 = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
	FenTestPosition6 = "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
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
	position      *Position
	history       []History
	repetitions   uint16
	autoThreefold bool
	legalMoves    []Move
	historyHashes []uint64
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
		autoThreefold: true,
	}

	pseudo := c.position.generateAllPseudoLegalMoves()
	legal := c.position.generateAllLegalMoves(pseudo)
	c.legalMoves = legal

	return c, nil
}

func (c *Chess) MakeMove(m Move) {
	c.position.MakeMove(m)
	pos := c.position.Copy()

	h := History{move: m, pos: *pos}
	c.history = append(c.history, h)
	c.historyHashes = append(c.historyHashes, pos.hash)

	c.calcRepetitions()
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

func (c *Chess) IsThreefoldRepetition() bool {
	return c.autoThreefold && c.repetitions >= 3
}

func (c *Chess) IsFivefoldRepetition() bool {
	return c.autoThreefold && c.repetitions >= 5
}
