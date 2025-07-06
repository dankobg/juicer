package engine

import (
	"fmt"
	"strings"
)

// Board represents the chess board and internally uses 12 bitboards for all pieces
type Board struct {
	pieceOccupancies [2][6]bitboard
	sideOccupancies  [3]bitboard
}

func (b *Board) calcSideOccupancies() {
	b.sideOccupancies[White] = b.pieceOccupancies[White][King] | b.pieceOccupancies[White][Queen] | b.pieceOccupancies[White][Rook] | b.pieceOccupancies[White][Bishop] | b.pieceOccupancies[White][Knight] | b.pieceOccupancies[White][Pawn]
	b.sideOccupancies[Black] = b.pieceOccupancies[Black][King] | b.pieceOccupancies[Black][Queen] | b.pieceOccupancies[Black][Rook] | b.pieceOccupancies[Black][Bishop] | b.pieceOccupancies[Black][Knight] | b.pieceOccupancies[Black][Pawn]
	b.sideOccupancies[Both] = b.sideOccupancies[White] | b.sideOccupancies[Black]
}

func (b *Board) bitboardForPiece(piece Piece) *bitboard {
	return &b.pieceOccupancies[piece.Color()][piece.Kind()]
}

// rotate90clockwise rotates the board 90 deg clockwise
func (b Board) rotate90clockwise() Board {
	var board Board

	for color, occupancies := range b.pieceOccupancies {
		for pk, occ := range occupancies {
			board.pieceOccupancies[color][pk] = occ.rotate90clockwise()
		}
	}

	board.calcSideOccupancies()

	return board
}

// rotate90counterClockwise rotates the board 90 deg counter-clockwise
func (b Board) rotate90counterClockwise() Board {
	var board Board

	for color, occupancies := range b.pieceOccupancies {
		for pk, occ := range occupancies {
			board.pieceOccupancies[color][pk] = occ.rotate90counterClockwise()
		}
	}

	board.calcSideOccupancies()

	return board
}

// rotate180 rotates the board 180 deg
func (b Board) rotate180() Board {
	var board Board

	for color, occupancies := range b.pieceOccupancies {
		for pk, occ := range occupancies {
			board.pieceOccupancies[color][pk] = occ.rotate180()
		}
	}

	board.calcSideOccupancies()

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
			if b.pieceOccupancies[color][pk].bitIsSet(sq) {
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

		for f := range 8 {
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

// isSquareAttacked checks if the square is attacked by the provided side
func (b Board) isSquareAttacked(sq Square, side Color, occupancy bitboard) bool {
	return pawnAttacksMask[side.Opposite()][sq]&b.pieceOccupancies[side][Pawn] != 0 ||
		kingAttacksMask[sq]&b.pieceOccupancies[side][King] != 0 ||
		knightsAttacksMask[sq]&b.pieceOccupancies[side][Knight] != 0 ||
		getBishopAttacks(sq, occupancy)&(b.pieceOccupancies[side][Bishop]|b.pieceOccupancies[side][Queen]) != 0 ||
		getRookAttacks(sq, occupancy)&(b.pieceOccupancies[side][Rook]|b.pieceOccupancies[side][Queen]) != 0
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
	return b.isSquareAttacked(Square(b.pieceOccupancies[side][King].LS1B()), side.Opposite(), b.sideOccupancies[Both])
}

// IsOnlyKingAndPawnLeft checks if only king and pawns are left
func (b Board) IsOnlyKingAndPawnLeft() bool {
	return b.pieceOccupancies[White][Pawn]|b.pieceOccupancies[White][King]|b.pieceOccupancies[Black][Pawn]|b.pieceOccupancies[Black][King] == b.sideOccupancies[Both]
}

// IsOnlyKingLeft checks if only kings are left
func (b Board) IsOnlyKingLeft() bool {
	return b.pieceOccupancies[White][King]|b.pieceOccupancies[Black][King] == b.sideOccupancies[Both]
}

// AlivePieces gets the total alive pieces count
func (b Board) AlivePieces() uint8 {
	return b.sideOccupancies[Both].populationCount()
}

// AlivePiecesForSide gets the total alive pieces count for given side
func (b Board) AlivePiecesForSide(turn Color) uint8 {
	return b.sideOccupancies[turn].populationCount()
}

// WhiteAlivePieces gets the total alive pieces count for white side
func (b Board) WhiteAlivePieces() uint8 {
	return b.sideOccupancies[White].populationCount()
}

// BlackAlivePieces gets the total alive pieces count for black side
func (b Board) BlackAlivePieces() uint8 {
	return b.sideOccupancies[White].populationCount()
}

// Copy returns the copy of a board
func (b Board) Copy() Board {
	return Board{
		pieceOccupancies: b.pieceOccupancies,
		sideOccupancies:  b.sideOccupancies,
	}
}
