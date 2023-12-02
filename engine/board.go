package juicer

import (
	"fmt"
	"strings"
)

// Board represents the chess board and internally uses 12 bitboards for all pieces
type Board struct {
	occupancies [2][6]bitboard
}

func (b *Board) bitboardForPiece(piece Piece) *bitboard {
	return &b.occupancies[piece.Color()][piece.Kind()]
}

func (b Board) whitePiecesOccupancy() bitboard {
	return b.occupancies[White][King] | b.occupancies[White][Queen] | b.occupancies[White][Rook] | b.occupancies[White][Bishop] | b.occupancies[White][Knight] | b.occupancies[White][Pawn]
}

func (b Board) blackPiecesOccupancy() bitboard {
	return b.occupancies[Black][King] | b.occupancies[Black][Queen] | b.occupancies[Black][Rook] | b.occupancies[Black][Bishop] | b.occupancies[Black][Knight] | b.occupancies[Black][Pawn]
}

func (b Board) piecesOccupancyForSide(side Color) bitboard {
	if side.IsWhite() {
		return b.whitePiecesOccupancy()
	}
	return b.blackPiecesOccupancy()
}

func (b Board) allPiecesOccupancy() bitboard {
	return b.whitePiecesOccupancy() | b.blackPiecesOccupancy()
}

// rotate90clockwise rotates the board 90 deg clockwise
func (b Board) rotate90clockwise() Board {
	var board Board

	for color, occupancies := range b.occupancies {
		for pk, occ := range occupancies {
			board.occupancies[color][pk] = occ.rotate90clockwise()
		}
	}

	return board
}

// rotate90counterClockwise rotates the board 90 deg counter-clockwise
func (b Board) rotate90counterClockwise() Board {
	var board Board

	for color, occupancies := range b.occupancies {
		for pk, occ := range occupancies {
			board.occupancies[color][pk] = occ.rotate90counterClockwise()
		}
	}

	return board
}

// rotate180 rotates the board 180 deg
func (b Board) rotate180() Board {
	var board Board

	for color, occupancies := range b.occupancies {
		for pk, occ := range occupancies {
			board.occupancies[color][pk] = occ.rotate180()
		}
	}

	return board
}

// Draw prints the board in 8x8 grid in ascii style
// it prints the piece fen symbol or `.` when there is no piece on a square
func (b Board) Draw(options *DrawOptions) string {
	return printBoard(options, func(sq Square) string {
		return b.pieceAt(sq).String()
	})
}

func (b Board) pieceAt(sq Square) Piece {
	for _, color := range colors {
		for _, pk := range pieceKinds {
			if b.occupancies[color][pk].bitIsSet(sq) {
				return NewPiece(pk, color)
			}
		}
	}

	return PieceNone
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

// isSquareAttackedByKing checks if the square is attacked by king by the provided side
func (b Board) isSquareAttackedByKing(sq Square, side Color) bool {
	return kingAttacksMask[sq]&b.occupancies[side][King] != 0
}

// isSquareAttackedByPawn checks if the square is attacked by pawn by the provided side
func (b Board) isSquareAttackedByPawn(sq Square, side Color) bool {
	return pawnAttacksMask[side.Opposite()][sq]&b.occupancies[side][Pawn] != 0
}

// isSquareAttackedByKnight checks if the square is attacked by knight by the provided side
func (b Board) isSquareAttackedByKnight(sq Square, side Color) bool {
	return knightsAttacksMask[sq]&b.occupancies[side][Knight] != 0
}

// isSquareAttackedByRook checks if the square is attacked by rook by the provided side
func (b Board) isSquareAttackedByRook(sq Square, side Color, occupancy bitboard) bool {
	return getRookAttacks(sq, occupancy)&b.occupancies[side][Rook] != 0
}

// isSquareAttackedByBishop checks if the square is attacked by bishop by the provided side
func (b Board) isSquareAttackedByBishop(sq Square, side Color, occupancy bitboard) bool {
	return getBishopAttacks(sq, occupancy)&b.occupancies[side][Bishop] != 0
}

// isSquareAttackedByQueen checks if the square is attacked by queen by the provided side
func (b Board) isSquareAttackedByQueen(sq Square, side Color, occupancy bitboard) bool {
	return getQueenAttacks(sq, occupancy)&b.occupancies[side][Queen] != 0
}

// isSquareAttackedByBishopOrQueen checks if the square is attacked by bishop or queen by the provided side
func (b Board) isSquareAttackedByBishopOrQueen(sq Square, side Color, occupancy bitboard) bool {
	return getBishopAttacks(sq, occupancy)&(b.occupancies[side][Bishop]|b.occupancies[side][Queen]) != 0
}

// isSquareAttackedByRookOrQueen checks if the square is attacked by rook or queen by the provided side
func (b Board) isSquareAttackedByRookOrQueen(sq Square, side Color, occupancy bitboard) bool {
	return getRookAttacks(sq, occupancy)&(b.occupancies[side][Rook]|b.occupancies[side][Queen]) != 0
}

// isSquareAttacked checks if the square is attacked by the provided side
func (b Board) isSquareAttacked(sq Square, side Color, occupancy bitboard) bool {
	return b.isSquareAttackedByPawn(sq, side) ||
		b.isSquareAttackedByKing(sq, side) ||
		b.isSquareAttackedByKnight(sq, side) ||
		b.isSquareAttackedByBishopOrQueen(sq, side, occupancy) ||
		b.isSquareAttackedByRookOrQueen(sq, side, occupancy)
}

// GetAttackedSquares gets the attacked squares by the provided side
func (b Board) GetAttackedSquares(side Color, mask, occupancy bitboard) bitboard {
	var attacked bitboard

	for mask > 0 {
		sq := Square(mask.PopLS1B())

		if b.isSquareAttacked(sq, side, occupancy) {
			attacked |= sq.occupancyMask()
		}
	}

	return attacked
}

// IsChecked checks if the provided side is in check
func (b Board) IsInCheck(side Color) bool {
	return b.isSquareAttacked(Square(b.occupancies[side][King].LS1B()), side.Opposite(), b.allPiecesOccupancy())
}

// IsOnlyKingAndPawnLeft checks if only king and pawns are left
func (b Board) IsOnlyKingAndPawnLeft() bool {
	return b.occupancies[White][Pawn]|b.occupancies[White][King]|b.occupancies[Black][Pawn]|b.occupancies[Black][King] == b.allPiecesOccupancy()
}

// IsOnlyKingLeft checks if only kings are left
func (b Board) IsOnlyKingLeft() bool {
	return b.occupancies[White][King]|b.occupancies[Black][King] == b.allPiecesOccupancy()
}

// AlivePieces gets the total alive pieces count
func (b Board) AlivePieces() uint8 {
	return b.allPiecesOccupancy().populationCount()
}

// AlivePiecesForSide gets the total alive pieces count for given side
func (b Board) AlivePiecesForSide(turn Color) uint8 {
	return b.piecesOccupancyForSide(turn).populationCount()
}

// WhiteAlivePieces gets the total alive pieces count for white side
func (b Board) WhiteAlivePieces() uint8 {
	return b.whitePiecesOccupancy().populationCount()
}

// BlackAlivePieces gets the total alive pieces count for black side
func (b Board) BlackAlivePieces() uint8 {
	return b.blackPiecesOccupancy().populationCount()
}

// Copy returns the copy of a board
func (b Board) Copy() Board {
	occupanciesCopy := b.occupancies

	return Board{
		occupancies: occupanciesCopy,
	}
}
