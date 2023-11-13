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

// Rotate90clockwise rotates the board 90 deg clockwise
func (b *Board) Rotate90clockwise() {
	b.whitePawnsOccupancy.rotate90clockwise()
}

// Rotate90counterClockwise rotates the board 90 deg counter-clockwise
func (b *Board) Rotate90counterClockwise() {
	b.whitePawnsOccupancy.rotate90counterClockwise()
}

// Rotate180 rotates the board 180 deg
func (b *Board) Rotate180() {
	b.whiteKingOccupancy.rotate180()
	b.whiteQueensOccupancy.rotate180()
	b.whiteRooksOccupancy.rotate180()
	b.whiteBishopsOccupancy.rotate180()
	b.whiteKnightsOccupancy.rotate180()
	b.whitePawnsOccupancy.rotate180()
	b.blackKingOccupancy.rotate180()
	b.blackQueensOccupancy.rotate180()
	b.blackRooksOccupancy.rotate180()
	b.blackBishopsOccupancy.rotate180()
	b.blackKnightsOccupancy.rotate180()
	b.blackPawnsOccupancy.rotate180()
}

// Draw prints the board in 8x8 grid with ascii style
func (b *Board) Draw(options *drawOptions) string {
	var sb strings.Builder

	opts := drawOptions{side: White}
	if options != nil {
		opts = *options
	}

	if !opts.compact {
		sb.WriteString("   +------------------------+\n")
	}

	for r := boardSize - 1; r >= 0; r-- {
		s1 := ""
		if !opts.compact {
			s1 = "|"
		}

		sb.WriteString(fmt.Sprintf(" %d %s", r+1, s1))

		for f := 0; f < 8; f++ {
			s2 := ""
			if !opts.compact {
				s2 = " "
			}

			square := Square(r*8 + f)
			piece := b.pieceAt(square)

			sb.WriteString(fmt.Sprintf(" %s%s", piece, s2))
		}

		s3 := ""
		if !opts.compact {
			s3 = "|"
		}

		sb.WriteString(fmt.Sprintf("%s \n", s3))
	}

	if !opts.compact {
		sb.WriteString("   +------------------------+\n")
	}

	sb.WriteString("     a  b  c  d  e  f  g  h")

	return sb.String()
}

func (b *Board) pieceAt(sq Square) Piece {
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
