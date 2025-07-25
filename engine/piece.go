package engine

import (
	"fmt"
)

// PieceKind represents the kind/type of the piece
type PieceKind int8

const (
	PieceKindNone PieceKind = iota - 1
	King
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

// String returns a symbol for a piece kind
func (ps PieceKind) String() string {
	switch ps {
	case King:
		return "k"
	case Queen:
		return "q"
	case Rook:
		return "r"
	case Bishop:
		return "b"
	case Knight:
		return "n"
	case Pawn:
		return "p"
	}

	return ""
}

// Piece is the piece kind with its color
type Piece int8

const (
	PieceNone Piece = iota - 1
	WhiteKing
	WhiteQueen
	WhiteRook
	WhiteBishop
	WhiteKnight
	WhitePawn
	BlackKing
	BlackQueen
	BlackRook
	BlackBishop
	BlackKnight
	BlackPawn
)

var (
	whitePieceSymbols = [6]string{"P", "R", "N", "B", "Q", "K"}
	blackPieceSymbols = [6]string{"p", "r", "n", "b", "q", "k"}
	pieceSymbols      = [12]string{"P", "R", "N", "B", "Q", "K", "p", "r", "n", "b", "q", "k"}

	whitePieceUnicodeSymbols = [6]string{"♔", "♕", "♖", "♗", "♘", "♙"}
	blackPieceUnicodeSymbols = [6]string{"♚", "♛", "♜", "♝", "♞", "♟"}
	pieceUnicodeSymbols      = [12]string{"♔", "♕", "♖", "♗", "♘", "♙", "♚", "♛", "♜", "♝", "♞", "♟"}

	pieceKinds = [6]PieceKind{King, Queen, Rook, Bishop, Knight, Pawn}

	whitePieces = [6]Piece{WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn}
	blackPieces = [6]Piece{BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn}
	pieces      = [12]Piece{WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn, BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn}

	fenPieces = map[string]Piece{
		"K": WhiteKing,
		"Q": WhiteQueen,
		"R": WhiteRook,
		"B": WhiteBishop,
		"N": WhiteKnight,
		"P": WhitePawn,
		"k": BlackKing,
		"q": BlackQueen,
		"r": BlackRook,
		"b": BlackBishop,
		"n": BlackKnight,
		"p": BlackPawn,
	}

	fenColors = map[string]Color{
		"w": White,
		"b": Black,
	}
)

// Kind returns the kind of the piece
func (p Piece) Kind() PieceKind {
	switch p {
	case WhiteKing, BlackKing:
		return King
	case WhiteQueen, BlackQueen:
		return Queen
	case WhiteRook, BlackRook:
		return Rook
	case WhiteBishop, BlackBishop:
		return Bishop
	case WhiteKnight, BlackKnight:
		return Knight
	case WhitePawn, BlackPawn:
		return Pawn
	}

	return PieceKindNone
}

// FENSymbol returns the FEN piece symbol
func (p Piece) FENSymbol() string {
	for symbol, piece := range fenPieces {
		if piece == p {
			return symbol
		}
	}

	return "."
}

// UnicodeSymbol returns the unicode piece symbol
func (p Piece) UnicodeSymbol() string {
	for idx, piece := range pieces {
		if piece == p {
			return pieceUnicodeSymbols[idx]
		}
	}

	return "."
}

// Color returns the color of the piece
func (p Piece) Color() Color {
	switch p {
	case WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn:
		return White
	case BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn:
		return Black
	}

	return ColorNone
}

func (p Piece) String() string {
	return p.FENSymbol()
}

func (p Piece) IsKing() bool {
	return p.Kind() == King
}

func (p Piece) IsQueen() bool {
	return p.Kind() == Queen
}

func (p Piece) IsRook() bool {
	return p.Kind() == Rook
}

func (p Piece) IsBishop() bool {
	return p.Kind() == Bishop
}

func (p Piece) IsKnight() bool {
	return p.Kind() == Knight
}

func (p Piece) IsPawn() bool {
	return p.Kind() == Pawn
}

func (p Piece) IsWhite() bool {
	return p.Color().IsWhite()
}

func (p Piece) IsBlack() bool {
	return p.Color().IsBlack()
}

func (p Piece) IsMajor() bool {
	return p.IsRook() || p.IsQueen()
}

func (p Piece) IsMinor() bool {
	return p.IsBishop() || p.IsKnight()
}

func (p Piece) IsPromotable() bool {
	return p.IsMajor() || p.IsMinor()
}

// NewPiece returns the piece given the kind and the color
func NewPiece(kind PieceKind, color Color) Piece {
	for _, p := range pieces {
		if p.Color() == color && p.Kind() == kind {
			return p
		}
	}

	return PieceNone
}

// NewPieceFromFenSymbol returns the piece given the fen piece symbol
func NewPieceFromFenSymbol(symbol string) (Piece, error) {
	piece, exists := fenPieces[symbol]
	if !exists {
		return PieceNone, fmt.Errorf("invalid piece symbol")
	}

	return piece, nil
}
