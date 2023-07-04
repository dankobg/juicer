package engine

import (
	"fmt"
	"math"
	"strings"
)

const (
	shortCastleNotation = "O-O"
	longCastleNotation  = "O-O-O"
	captureNotation     = "x"
	checkNotation       = "+"
	checkmateNotation   = "#"
)

type Move struct {
	color         Color
	fromSquare    Square
	toSquare      Square
	pieceMoved    Piece
	pieceCaptured *Piece
	enpSquare     *Square
	// castlingRights castlingRights
	promotion *PromotionPieceSymbol
}

func NewMove(color Color, from Square, to Square, piece Piece, pieceCaptured *Piece, enpSquare *Square, castleMove CastlingRights, promotion *PromotionPieceSymbol) *Move {
	return &Move{
		color:         color,
		fromSquare:    from,
		toSquare:      to,
		pieceMoved:    piece,
		pieceCaptured: pieceCaptured,
		enpSquare:     enpSquare,
		promotion:     promotion,
	}
}

func (m *Move) String() string {
	return m.FromCoordinate().String() + m.ToCoordinate().String()
}

func (m *Move) ID() int {
	return int(m.FromRow())*1000 + int(m.FromColumn())*100 + int(m.ToRow())*10 + int(m.ToColumn())
}

func (m *Move) Color() Color {
	return m.color
}

func (m *Move) FromSquare() Square {
	return m.fromSquare
}

func (m *Move) ToSquare() Square {
	return m.toSquare
}

func (m *Move) PieceMoved() Piece {
	return *m.fromSquare.piece
}

func (m *Move) PieceCaptured() *Piece {
	return m.toSquare.piece
}

func (m *Move) EnpSquare() *Square {
	return m.enpSquare
}

func (m *Move) Promotion() *PromotionPieceSymbol {
	return m.promotion
}

func (m *Move) FromRow() Row {
	return m.fromSquare.row
}

func (m *Move) FromColumn() Column {
	return m.fromSquare.column
}

func (m *Move) FromFile() File {
	return m.fromSquare.File()
}

func (m *Move) FromRank() Rank {
	return m.fromSquare.Rank()
}

func (m *Move) FromCoordinate() Coordinate {
	return m.fromSquare.Coordinate()
}

func (m *Move) ToRow() Row {
	return m.toSquare.row
}

func (m *Move) ToColumn() Column {
	return m.toSquare.column
}

func (m *Move) ToFile() File {
	return m.toSquare.File()
}

func (m *Move) ToRank() Rank {
	return m.toSquare.Rank()
}

func (m *Move) ToCoordinate() Coordinate {
	return m.toSquare.Coordinate()
}

func moveToSAN(move Move, moves []Move, cr CastlingRights) string {
	var san string

	isCastle, castleSide := move.IsCastle(cr)
	if isCastle {
		if castleSide == kingSideCastle {
			san += shortCastleNotation
		}

		if castleSide == queenSideCastle {
			san += longCastleNotation
		}
	}

	if !move.pieceMoved.IsPawn() {
		unique := getUnambiguousMoveNotation(move, moves)
		fmt.Printf("UNIQUE: %+v\n", unique)
		san += strings.ToUpper(move.pieceMoved.symbol.String()) + unique
	}

	if move.pieceCaptured != nil {
		if move.pieceMoved.IsPawn() {
			san += move.ToFile().String() + move.ToRank().String()
		}

		san += captureNotation
	}

	san += move.ToCoordinate().String()

	if move.IsPromotion() {
		san += strings.ToUpper(move.promotion.String())
	}

	// makeMove()

	// if isCheck {
	// 	if mate {
	// 		san += checkmateNotation
	// 	} else {
	// 		san += checkNotation
	// 	}
	// }

	// undomove()

	return san
}

func (m *Move) SAN(moves []Move, cr CastlingRights) string {
	return moveToSAN(*m, moves, cr)
}

func (m *Move) LAN() string {
	return ""
}

func (m *Move) IsWhite() bool {
	return m.color.IsWhite()
}

func (m *Move) IsBlack() bool {
	return m.color.IsBlack()
}

// IsCastle returns true for a castling move and castlingSide ("k" or "q" or "" if non castling move)
func (m *Move) IsCastle(cr CastlingRights) (bool, string) {
	if !m.pieceMoved.IsKing() {
		return false, ""
	}

	if math.Abs(float64(m.FromColumn())-float64(m.ToColumn())) != 2 {
		return false, ""
	}

	if int(m.ToColumn()) == 6 && (((cr & WhiteKingSideCastle) > 0) || ((cr & BlackKingSideCastle) > 0)) {
		return true, kingSideCastle
	}

	if int(m.ToColumn()) == 2 && (((cr & WhiteQueenSideCastle) > 0) || ((cr & BlackQueenSideCastle) > 0)) {
		return true, queenSideCastle
	}

	return false, ""
}

func (m *Move) IsEnPassant(enpSquare *Square) bool {
	return m.pieceMoved.IsPawn() && enpSquare != nil && enpSquare.CoordEquals(m.ToSquare())
}

func (m *Move) IsPromotion() bool {
	onLastRank := (m.IsWhite() && m.ToRow() == 0) || (m.IsBlack() && m.ToRow() == 7)
	return m.pieceMoved.IsPawn() && onLastRank
}

func (m *Move) GetPromotedPiece() (*Piece, PromotionPieceSymbol) {
	promPiece := QueenPromotion

	if m.promotion != nil {
		promPiece = *m.promotion
	}

	if m.IsPromotion() {
		return NewPiece(promPiece.ToPieceSymbol(), m.color, true), promPiece
	}

	return nil, promPiece
}
