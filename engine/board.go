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
}

func (b Board) whitePiecesOccupancy() bitboard {
	return b.whiteKingOccupancy | b.whiteQueensOccupancy | b.whiteRooksOccupancy | b.whiteBishopsOccupancy | b.whiteKnightsOccupancy | b.whitePawnsOccupancy
}

func (b Board) blackPiecesOccupancy() bitboard {
	return b.blackKingOccupancy | b.blackQueensOccupancy | b.blackRooksOccupancy | b.blackBishopsOccupancy | b.blackKnightsOccupancy | b.blackPawnsOccupancy
}

func (b Board) allPiecesOccupancy() bitboard {
	return b.whitePiecesOccupancy() | b.blackPiecesOccupancy()
}

// rotate90clockwise rotates the board 90 deg clockwise
func (b Board) rotate90clockwise() Board {
	return Board{
		whiteKingOccupancy:    b.whiteKingOccupancy.rotate90clockwise(),
		whiteQueensOccupancy:  b.whiteQueensOccupancy.rotate90clockwise(),
		whiteRooksOccupancy:   b.whiteRooksOccupancy.rotate90clockwise(),
		whiteBishopsOccupancy: b.whiteBishopsOccupancy.rotate90clockwise(),
		whiteKnightsOccupancy: b.whiteKnightsOccupancy.rotate90clockwise(),
		whitePawnsOccupancy:   b.whitePawnsOccupancy.rotate90clockwise(),
		blackKingOccupancy:    b.blackKingOccupancy.rotate90clockwise(),
		blackQueensOccupancy:  b.blackQueensOccupancy.rotate90clockwise(),
		blackRooksOccupancy:   b.blackRooksOccupancy.rotate90clockwise(),
		blackBishopsOccupancy: b.blackBishopsOccupancy.rotate90clockwise(),
		blackKnightsOccupancy: b.blackKnightsOccupancy.rotate90clockwise(),
		blackPawnsOccupancy:   b.blackPawnsOccupancy.rotate90clockwise(),
	}
}

// rotate90counterClockwise rotates the board 90 deg counter-clockwise
func (b Board) rotate90counterClockwise() Board {
	return Board{
		whiteKingOccupancy:    b.whiteKingOccupancy.rotate90counterClockwise(),
		whiteQueensOccupancy:  b.whiteQueensOccupancy.rotate90counterClockwise(),
		whiteRooksOccupancy:   b.whiteRooksOccupancy.rotate90counterClockwise(),
		whiteBishopsOccupancy: b.whiteBishopsOccupancy.rotate90counterClockwise(),
		whiteKnightsOccupancy: b.whiteKnightsOccupancy.rotate90counterClockwise(),
		whitePawnsOccupancy:   b.whitePawnsOccupancy.rotate90counterClockwise(),
		blackKingOccupancy:    b.blackKingOccupancy.rotate90counterClockwise(),
		blackQueensOccupancy:  b.blackQueensOccupancy.rotate90counterClockwise(),
		blackRooksOccupancy:   b.blackRooksOccupancy.rotate90counterClockwise(),
		blackBishopsOccupancy: b.blackBishopsOccupancy.rotate90counterClockwise(),
		blackKnightsOccupancy: b.blackKnightsOccupancy.rotate90counterClockwise(),
		blackPawnsOccupancy:   b.blackPawnsOccupancy.rotate90counterClockwise(),
	}
}

// rotate180 rotates the board 180 deg
func (b Board) rotate180() Board {
	return Board{
		whiteKingOccupancy:    b.whiteKingOccupancy.rotate180(),
		whiteQueensOccupancy:  b.whiteQueensOccupancy.rotate180(),
		whiteRooksOccupancy:   b.whiteRooksOccupancy.rotate180(),
		whiteBishopsOccupancy: b.whiteBishopsOccupancy.rotate180(),
		whiteKnightsOccupancy: b.whiteKnightsOccupancy.rotate180(),
		whitePawnsOccupancy:   b.whitePawnsOccupancy.rotate180(),
		blackKingOccupancy:    b.blackKingOccupancy.rotate180(),
		blackQueensOccupancy:  b.blackQueensOccupancy.rotate180(),
		blackRooksOccupancy:   b.blackRooksOccupancy.rotate180(),
		blackBishopsOccupancy: b.blackBishopsOccupancy.rotate180(),
		blackKnightsOccupancy: b.blackKnightsOccupancy.rotate180(),
		blackPawnsOccupancy:   b.blackPawnsOccupancy.rotate180(),
	}
}

// Draw prints the board in 8x8 grid in ascii style
// it prints the piece fen symbol or `.` when there is no piece on a square
func (b Board) Draw(options *DrawOptions) string {
	return printBoard(options, func(sq Square) string {
		return b.pieceAt(sq).String()
	})
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

func (b Board) bitboardForPiece(piece Piece) bitboard {
	switch piece {
	case WhiteKing:
		return b.whiteKingOccupancy
	case WhiteQueen:
		return b.whiteQueensOccupancy
	case WhiteRook:
		return b.whiteRooksOccupancy
	case WhiteBishop:
		return b.whiteBishopsOccupancy
	case WhiteKnight:
		return b.whiteKnightsOccupancy
	case WhitePawn:
		return b.whitePawnsOccupancy
	case BlackKing:
		return b.blackKingOccupancy
	case BlackQueen:
		return b.blackQueensOccupancy
	case BlackRook:
		return b.blackRooksOccupancy
	case BlackBishop:
		return b.blackBishopsOccupancy
	case BlackKnight:
		return b.blackKnightsOccupancy
	case BlackPawn:
		return b.blackPawnsOccupancy
	}

	return bitboardEmpty
}

// FenPositionPart returns the fen position part without the metadata (turn, enpSq, castle, half/full move clock)
// e.g. fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" -> fenPositionPart: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
func (b Board) FenPositionPart() string {
	var sb strings.Builder

	for r := boardSize - 1; r >= 0; r-- {
		emptySquares := 0

		for f := 0; f < 8; f++ {
			sq := Square(r*8 + f)

			piece := b.pieceAt(sq)

			if piece != PieceNone {
				if emptySquares > 0 {
					sb.WriteString(fmt.Sprint(emptySquares))
					emptySquares = 0
				}

				sb.WriteString(piece.FENSymbol())
			} else {
				emptySquares++
			}
		}

		if emptySquares > 0 {
			sb.WriteString(fmt.Sprint(emptySquares))
		}

		if r > 0 {
			sb.WriteString(fenPositionSeparator)
		}
	}

	return sb.String()
}
