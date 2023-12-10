package juicer

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

type Chess struct {
	position *Position
	history  []Position
}

func InitPrecalculatedTables() {
	initAllAttackMasksTables()
	initZobrist()
}

func (c *Chess) MakeMove(m Move) {

}
