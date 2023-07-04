package engine

import (
	"regexp"
	"strings"
)

type PieceSymbol string

const (
	Pawn   PieceSymbol = "p"
	Knight PieceSymbol = "n"
	Bishop PieceSymbol = "b"
	Rook   PieceSymbol = "r"
	Queen  PieceSymbol = "q"
	King   PieceSymbol = "k"
)

type PieceFENSymbol string

const (
	WhitePawn   PieceFENSymbol = "P"
	WhiteKnight PieceFENSymbol = "N"
	WhiteBishop PieceFENSymbol = "B"
	WhiteRook   PieceFENSymbol = "R"
	WhiteQueen  PieceFENSymbol = "Q"
	WhiteKing   PieceFENSymbol = "K"
	BlackPawn   PieceFENSymbol = "p"
	BlackKnight PieceFENSymbol = "n"
	BlackBishop PieceFENSymbol = "b"
	BlackRook   PieceFENSymbol = "r"
	BlackQueen  PieceFENSymbol = "q"
	BlackKing   PieceFENSymbol = "k"
)

type PromotionPieceSymbol string

const (
	KnightPromotion PromotionPieceSymbol = "n"
	BishopPromotion PromotionPieceSymbol = "b"
	RookPromotion   PromotionPieceSymbol = "r"
	QueenPromotion  PromotionPieceSymbol = "q"
)

var reBlackFENPieces = regexp.MustCompile("^[prnbqk]$")
var rePromotionPieces = regexp.MustCompile("^[rnbq]$")
var reWhiteFENPieces = regexp.MustCompile("^[PRNBQK]$")
var reWhiteAndBlackFENPieces = regexp.MustCompile("^[prnbqkPRNBQK]$")

func (ps PieceSymbol) String() string {
	return string(ps)
}

func (ps PieceSymbol) Name() string {
	switch ps {
	case Pawn:
		return "Pawn"
	case Rook:
		return "Rook"
	case Bishop:
		return "Bishop"
	case Knight:
		return "Knight"
	case King:
		return "King"
	case Queen:
		return "Queen"
	default:
		return "Unknown"
	}
}

func (ps PieceSymbol) ToPieceFENSymbol(color Color) PieceFENSymbol {
	if color == White {
		return PieceFENSymbol(strings.ToUpper(ps.String()))
	}

	return PieceFENSymbol(ps.String())
}

func (ps PieceSymbol) IsPawn() bool   { return ps == Pawn }
func (ps PieceSymbol) IsQueen() bool  { return ps == Queen }
func (ps PieceSymbol) IsKing() bool   { return ps == King }
func (ps PieceSymbol) IsRook() bool   { return ps == Rook }
func (ps PieceSymbol) IsKnight() bool { return ps == Knight }
func (ps PieceSymbol) IsBishop() bool { return ps == Bishop }

func (pfs PieceFENSymbol) String() string {
	return string(pfs)
}

func (pfs PieceFENSymbol) ToPieceSymbol() PieceSymbol {
	return PieceSymbol(strings.ToLower(pfs.String()))
}

func (pps PromotionPieceSymbol) String() string {
	return string(pps)
}

func (pps PromotionPieceSymbol) Valid() bool {
	return rePromotionPieces.Match([]byte(pps.String()))
}

func (pps PromotionPieceSymbol) ToPieceSymbol() PieceSymbol {
	return PieceSymbol(pps.String())
}

func (ps PieceSymbol) Valid() bool {
	return reBlackFENPieces.Match([]byte(ps.String()))
}

func (pfs PieceFENSymbol) Valid() bool {
	return reWhiteAndBlackFENPieces.Match([]byte(pfs.String()))
}

func (pfs PieceFENSymbol) GetColorAndSymbolPair() (Color, PieceSymbol) {
	var color Color

	if reBlackFENPieces.Match([]byte(pfs.String())) {
		color = Black
	} else if reWhiteFENPieces.Match([]byte(pfs.String())) {
		color = White
	}

	return color, pfs.ToPieceSymbol()
}

type Piece struct {
	symbol PieceSymbol
	color  Color
	alive  bool
}

func NewPiece(symbol PieceSymbol, color Color, alive bool) *Piece {
	return &Piece{
		symbol: symbol,
		color:  color,
		alive:  alive,
	}
}

func NewPieceFromFENSymbol(fenSymbol PieceFENSymbol, alive bool) *Piece {
	color, symbol := fenSymbol.GetColorAndSymbolPair()

	return &Piece{
		symbol: symbol,
		color:  color,
		alive:  alive,
	}
}

func (p *Piece) Symbol() PieceSymbol {
	return p.symbol
}

func (p *Piece) Color() Color {
	return p.color
}

func (p *Piece) Alive() bool {
	return p.alive
}

func (p *Piece) IsPawn() bool {
	return p.symbol.IsPawn()
}

func (p *Piece) IsKing() bool {
	return p.symbol.IsKing()
}

func (p *Piece) IsQueen() bool {
	return p.symbol.IsQueen()
}

func (p *Piece) IsKnight() bool {
	return p.symbol.IsKnight()
}

func (p *Piece) IsBishop() bool {
	return p.symbol.IsBishop()
}

func (p *Piece) IsRook() bool {
	return p.symbol.IsRook()
}

func (p *Piece) IsWhite() bool {
	return p.color.IsWhite()
}

func (p *Piece) IsBlack() bool {
	return p.color.IsBlack()
}

func (p *Piece) String() string {
	return p.ToFENSymbol().String()
}

func (p *Piece) Name() string {
	return p.symbol.Name()
}

func (p *Piece) IsFriendly(curentTurn Color) bool {
	return p.color.Equals(curentTurn)
}

func (p *Piece) IsEnemy(curentTurn Color) bool {
	return !p.color.Equals(curentTurn)
}

func (p *Piece) SymbolEquals(piece Piece) bool {
	return p.symbol.String() == piece.symbol.String()
}

func (p *Piece) Equals(piece Piece) bool {
	return p.symbol.String() == piece.symbol.String() && p.color.Equals(piece.color) && p.alive == piece.alive
}

func (p *Piece) ToFENSymbol() PieceFENSymbol {
	if p.IsWhite() {
		return PieceFENSymbol(strings.ToUpper(p.symbol.String()))
	}

	return PieceFENSymbol(p.symbol.String())
}
