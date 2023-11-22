package juicer

import "fmt"

type Position struct {
	board                *Board
	turn                 Color
	enpSquare            Square
	castleRights         CastleRights
	halfMoveClock        uint8
	fullMoveClock        uint16
	check                bool
	checkmate            bool
	stalemate            bool
	draw                 bool
	threeFold            bool
	insufficientMaterial bool
	terminated           bool
	outcome              string
	comments             []string
	headers              []string
	capturedPieces       []Piece
	alivePieces          []Piece
}

func (p *Position) PrintBoard() string {
	return p.board.Draw(nil)
}

func (p *Position) LoadFromFEN(fen string) error {
	meta, err := validateFEN(fen, validateFenOps{})
	if err != nil {
		return fmt.Errorf("failed to load position from fen: %w", err)
	}

	board := Board{}

	for sq, piece := range meta.squares {
		if piece == WhiteKing {
			board.whiteKingOccupancy.setBit(sq)
		} else if piece == WhiteQueen {
			board.whiteKingOccupancy.setBit(sq)
		} else if piece == WhiteRook {
			board.whiteRooksOccupancy.setBit(sq)
		} else if piece == WhiteBishop {
			board.whiteBishopsOccupancy.setBit(sq)
		} else if piece == WhiteKnight {
			board.whiteKnightsOccupancy.setBit(sq)
		} else if piece == WhitePawn {
			board.whitePawnsOccupancy.setBit(sq)
		} else if piece == BlackKing {
			board.blackKingOccupancy.setBit(sq)
		} else if piece == BlackQueen {
			board.blackQueensOccupancy.setBit(sq)
		} else if piece == BlackRook {
			board.blackRooksOccupancy.setBit(sq)
		} else if piece == BlackBishop {
			board.blackBishopsOccupancy.setBit(sq)
		} else if piece == BlackKnight {
			board.blackKnightsOccupancy.setBit(sq)
		} else if piece == BlackPawn {
			board.blackPawnsOccupancy.setBit(sq)
		}
	}

	p.board = &board
	p.turn = meta.turnColor
	p.enpSquare = meta.enpSquare
	p.castleRights = meta.castleRights
	p.halfMoveClock = meta.halfMoveClock
	p.fullMoveClock = meta.fullMoveClock

	// p.check = false
	// p.checkmate = false
	// p.stalemate = false
	// p.draw = false
	// p.threeFold = false
	// p.insufficientMaterial = false
	// p.terminated = false
	// p.outcome = ""
	// p.comments = []string{}
	// p.headers = []string{}
	// p.capturedPieces = []Piece{}
	// p.alivePieces = []Piece{}

	return nil
}

// FenMetaPart returns the fen meta part without the position and it includes the empty string ` ` at start
// e.g. fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" -> fenMetaPart: " w KQkq - 0 1"
func (p *Position) FenMetaPart() string {
	castleToken := p.castleRights.ToFEN()

	enpSqToken := fenNoneSymbol
	if p.enpSquare != SquareNone {
		enpSqToken = p.enpSquare.Coordinate()
	}

	fenMetaPart := fmt.Sprintf(" %s %s %s %d %d", p.turn, castleToken, enpSqToken, p.halfMoveClock, p.fullMoveClock)
	return fenMetaPart
}

// Fen returns the full fen string
func (p *Position) Fen() string {
	return p.board.FenPositionPart() + p.FenMetaPart()
}

func (p *Position) whiteCanCastleKingSide() bool {
	return p.castleRights.whiteCanCastleKingSide()
}

func (p *Position) whiteCanCastleQueenSide() bool {
	return p.castleRights.whiteCanCastleQueenSide()
}

func (p *Position) whiteCanCastle() bool {
	return p.castleRights.whiteCanCastle()
}

func (p *Position) blackCanCastleKingSide() bool {
	return p.castleRights.blackCanCastleKingSide()
}

func (p *Position) blackCanCastleQueenSide() bool {
	return p.castleRights.blackCanCastleQueenSide()
}

func (p *Position) blackCanCastle() bool {
	return p.castleRights.blackCanCastle()
}

func (p *Position) canCastleKingSide() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastleKingSide()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastleKingSide()
	}
	return false
}

func (p *Position) canCastleQueenSide() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastleQueenSide()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastleQueenSide()
	}
	return false
}

func (p *Position) canCastle() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastle()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastle()
	}
	return false
}
