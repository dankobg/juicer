package engine

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

func (p *Position) FEN() string {
	return FENStartingPosition
}

func (p *Position) LoadFromFEN(fen string) {
	fmt.Println("fen", fen)
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
