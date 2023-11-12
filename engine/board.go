package juicer

import (
	"fmt"
	"strings"
)

// Board represents the chess board and internally uses 12 bitboards for all pieces
type Board struct {
	whiteKingOccupancy    bitboard
	whiteQueensOccupancy  bitboard
	whiteRooksOccupancy   bitboard
	whiteBishopsOccupancy bitboard
	whiteKnightsOccupancy bitboard
	whitePawnsOccupancy   bitboard

	blackKingOccupancy    bitboard
	blackQueensOccupancy  bitboard
	blackRooksOccupancy   bitboard
	blackBishopsOccupancy bitboard
	blackKnightsOccupancy bitboard
	blackPawnsOccupancy   bitboard

	whitePiecesOccupancy bitboard
	blackPiecesOccupancy bitboard
	allPiecesOccupancy   bitboard
}

// DrawCompact prints the bitboard for debugging as 8x8 grid of 0s and 1s in a compact way
func (b Board) DrawCompact() string {
	var sb strings.Builder

	for r := boardSize - 1; r >= 0; r-- {
		sb.WriteString(fmt.Sprintf(" %d ", r+1))

		for f := 0; f < boardSize; f++ {
			square := Square(r*8 + f)
			piece := b.pieceAt(square)
			sb.WriteString(fmt.Sprintf(" %s", piece.String()))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("\n    a b c d e f g h")

	return sb.String()
}

// DrawPretty prints the board and pieces visually in 8x8 grid
func (b Board) DrawPretty() string {
	var sb strings.Builder
	sb.WriteString("   +------------------------+\n")

	for r := boardSize - 1; r >= 0; r-- {
		sb.WriteString(fmt.Sprintf(" %d |", r+1))

		for f := 0; f < 8; f++ {
			square := Square(r*8 + f)
			piece := b.pieceAt(square)
			sb.WriteString(fmt.Sprintf(" %s ", piece.String()))
		}

		sb.WriteString("| \n")
	}

	sb.WriteString("   +------------------------+\n")
	sb.WriteString("     a  b  c  d  e  f  g  h")

	return sb.String()
}

func (b Board) pieceAt(sq Square) Piece {
	if b.whiteKingOccupancy.bitIsSet(sq) {
		return WhiteKing
	} else if b.whiteQueensOccupancy.bitIsSet(sq) {
		return WhiteQueen
	} else if b.whiteRooksOccupancy.bitIsSet(sq) {
		return WhiteRook
	} else if b.whiteBishopsOccupancy.bitIsSet(sq) {
		return WhiteBishop
	} else if b.whiteKnightsOccupancy.bitIsSet(sq) {
		return WhiteKnight
	} else if b.whitePawnsOccupancy.bitIsSet(sq) {
		return WhitePawn
	} else if b.blackKingOccupancy.bitIsSet(sq) {
		return BlackKing
	} else if b.blackQueensOccupancy.bitIsSet(sq) {
		return BlackQueen
	} else if b.blackRooksOccupancy.bitIsSet(sq) {
		return BlackRook
	} else if b.blackBishopsOccupancy.bitIsSet(sq) {
		return BlackBishop
	} else if b.blackKnightsOccupancy.bitIsSet(sq) {
		return BlackKnight
	} else if b.blackPawnsOccupancy.bitIsSet(sq) {
		return BlackPawn
	} else {
		return PieceNone
	}
}

// var whitePiecesOccupancy bitboard = whiteKingOccupancy |whiteQueensOccupancy |whiteRooskOccupancy |whiteBishopsOccupancy |whiteKnightsOccupancy |whitePawnsOccupancy
// var blackPiecesOccupancy bitboard = blackKingOccupancy |blackQueensOccupancy |blackRooskOccupancy |blackBishopsOccupancy |blackKnightsOccupancy |blackPawnsOccupancy
// var allPiecesOccupancy = whitePiecesOccupancy | blackPiecesOccupancy
