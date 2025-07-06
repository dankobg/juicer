package engine

import (
	"fmt"
	"strings"
)

type MoveType uint8

const (
	MoveTypeQuiet MoveType = iota
	MoveTypeCapture
	MoveTypeDoublePawn
	MoveTypeEnpCapture
	MoveTypePromotion
	MoveTypePromotionCapture
	MoveTypeCastle
)

func newQuietMove(src, dest Square, piece Piece) Move {
	return newMove(src, dest, piece, PromotionNone, false, false, false, false)
}

func newCaptureMove(src, dest Square, piece Piece) Move {
	return newMove(src, dest, piece, PromotionNone, true, false, false, false)
}

func newDoublePawnMove(src, dest Square, piece Piece) Move {
	return newMove(src, dest, piece, PromotionNone, false, true, false, false)
}

func newEnpCaptureMove(src, dest Square, piece Piece) Move {
	return newMove(src, dest, piece, PromotionNone, true, false, true, false)
}

func newPromotionMove(src, dest Square, piece Piece, promo Promotion) Move {
	return newMove(src, dest, piece, promo, false, false, false, false)
}

func newPromotionCaptureMove(src, dest Square, piece Piece, promo Promotion) Move {
	return newMove(src, dest, piece, promo, true, false, false, false)
}

func newCastleMove(src, dest Square, piece Piece) Move {
	return newMove(src, dest, piece, PromotionNone, false, false, false, true)
}

func newPossiblePromotionMoves(src, dest Square, piece Piece) []Move {
	return []Move{
		newMove(src, dest, piece, PromotionQueen, false, false, false, false),
		newMove(src, dest, piece, PromotionRook, false, false, false, false),
		newMove(src, dest, piece, PromotionBishop, false, false, false, false),
		newMove(src, dest, piece, PromotionKnight, false, false, false, false),
	}
}

func newPossiblePromotionCaptureMoves(src, dest Square, piece Piece) []Move {
	return []Move{
		newMove(src, dest, piece, PromotionQueen, true, false, false, false),
		newMove(src, dest, piece, PromotionRook, true, false, false, false),
		newMove(src, dest, piece, PromotionBishop, true, false, false, false),
		newMove(src, dest, piece, PromotionKnight, true, false, false, false),
	}
}

type Promotion uint8

const (
	PromotionNone Promotion = iota
	PromotionQueen
	PromotionRook
	PromotionBishop
	PromotionKnight
)

func (prm Promotion) IsPromotion() bool {
	return prm != PromotionNone
}

func (prm Promotion) PieceKind() PieceKind {
	switch prm {
	case PromotionQueen:
		return Queen
	case PromotionRook:
		return Rook
	case PromotionBishop:
		return Bishop
	case PromotionKnight:
		return Knight
	}

	return PieceKindNone
}

const (
	srcShift        = 0
	destShift       = 6
	pieceShift      = 12
	promotionShift  = 16
	captureShift    = 20
	doublePawnShift = 21
	enPassantShift  = 22
	castleShift     = 23

	squareMask    = 0x3F
	pieceMask     = 0xF
	promotionMask = 0x7
)

// Move is mapped from: [0..5] src, [6..11] dest, [12..15] piece, [16..18] promotion, 19 capture, 20 doublePawn, 21 en-passant, 22 castle
type Move int32

func newMove(src, dest Square, piece Piece, promotion Promotion, capture, doublePawn, enPassant, castle bool) Move {
	var m Move

	m |= Move(src) << srcShift
	m |= Move(dest) << destShift
	m |= Move(piece) << pieceShift
	m |= Move(promotion) << promotionShift

	if capture {
		m |= 1 << captureShift
	}
	if doublePawn {
		m |= 1 << doublePawnShift
	}
	if enPassant {
		m |= 1 << enPassantShift
	}
	if castle {
		m |= 1 << castleShift
	}

	return m
}

// Src is the source square
func (m Move) Src() Square {
	return Square((m >> srcShift) & squareMask)
}

// Dest is the destination square
func (m Move) Dest() Square {
	return Square((m >> destShift) & squareMask)
}

// SrcDest is the source, destination tuple
func (m Move) SrcDest() (Square, Square) {
	return m.Src(), m.Dest()
}

// Piece is the piece moved
func (m Move) Piece() Piece {
	return Piece((m >> pieceShift) & pieceMask)
}

// Promotion is the promotion
func (m Move) Promotion() Promotion {
	return Promotion((m >> promotionShift) & promotionMask)
}

// IsCapture checks whether the move is capture
func (m Move) IsCapture() bool {
	return ((m >> captureShift) & 1) == 1
}

// IsDoublePawn checks whether the move is double pawn push
func (m Move) IsDoublePawn() bool {
	return ((m >> doublePawnShift) & 1) == 1
}

// IsEnPassant checks whether the move is en-passant
func (m Move) IsEnPassant() bool {
	return ((m >> enPassantShift) & 1) == 1
}

// IsCastle checks whether the move is castle
func (m Move) IsCastle() bool {
	return ((m >> castleShift) & 1) == 1
}

// IsKingSideCastle checks whether the move is castle king side
func (m Move) IsKingSideCastle() bool {
	return m.IsCastle() && m.Piece().Kind() == King && m.Dest().File() > m.Src().File()
}

// IsQueenSideCastle checks whether the move is castle queen side
func (m Move) IsQueenSideCastle() bool {
	return m.IsCastle() && m.Piece().Kind() == King && m.Dest().File() < m.Src().File()
}

func (m Move) String() string {
	return m.ToUCI()
}

func (m Move) ToUCI() string {
	var promo string

	if m.Promotion().IsPromotion() {
		switch m.Promotion() {
		case PromotionQueen:
			promo = Queen.String()
		case PromotionRook:
			promo = Rook.String()
		case PromotionBishop:
			promo = Bishop.String()
		case PromotionKnight:
			promo = Knight.String()
		}
	}

	return fmt.Sprintf("%s%s%s", m.Src(), m.Dest(), promo)
}

// <LAN move descriptor piece moves> ::= <Piece symbol><from square>['-'|'x']<to square>
// <LAN move descriptor pawn moves>  ::= <from square>['-'|'x']<to square>[<promoted to>]
// <Piece symbol> ::= 'N' | 'B' | 'R' | 'Q' | 'K'
func (m Move) ToLAN(p *Position, isCheck, isCheckmate bool) string {
	var piece string
	if m.Piece().Kind() != Pawn {
		piece = strings.ToUpper(m.Piece().Kind().String())
	}
	separator := "-"
	if m.IsCapture() {
		separator = "x"
	}
	var promo string
	if m.Promotion().IsPromotion() {
		promo = fmt.Sprintf("=%s", strings.ToUpper(m.Promotion().PieceKind().String()))
	}
	var checkOrCheckmate string
	if isCheck {
		checkOrCheckmate = "+"
	}
	if isCheckmate {
		checkOrCheckmate = "#"
	}
	return fmt.Sprintf("%s%s%s%s%s%s", piece, m.Src().String(), separator, m.Dest().String(), promo, checkOrCheckmate)
}

// <SAN move descriptor piece moves>   ::= <Piece symbol>[<from file>|<from rank>|<from square>]['x']<to square>
// <SAN move descriptor pawn captures> ::= <from file>[<from rank>] 'x' <to square>[<promoted to>]
// <SAN move descriptor pawn push>     ::= <to square>[<promoted to>]
func (m Move) ToSAN(p *Position, isCheck, isCheckmate bool, legalMoves []Move) string {
	if m.IsCastle() {
		if m.IsKingSideCastle() {
			if isCheckmate {
				return "O-O#"
			}
			if isCheck {
				return "O-O+"
			}
			return "O-O"
		}
		if m.IsQueenSideCastle() {
			if isCheckmate {
				return "O-O-O#"
			}
			if isCheck {
				return "O-O-O+"
			}
			return "O-O-O"
		}
	}
	var piece string
	if m.Piece().Kind() != Pawn {
		piece = strings.ToUpper(m.Piece().Kind().String())
	}
	var disambiguation string
	if m.Piece().Kind() != Pawn {
		ambiguous := make([]Move, 0)
		for _, lm := range legalMoves {
			if lm != m && lm.Dest() == m.Dest() && lm.Piece() == m.Piece() {
				ambiguous = append(ambiguous, lm)
			}
		}
		if len(ambiguous) > 0 {
			var ambigFile, ambigRank bool
			for _, amb := range ambiguous {
				if amb.Src().File() == m.Src().File() {
					ambigFile = true
				}
				if amb.Src().Rank() == m.Src().Rank() {
					ambigRank = true
				}
			}
			if ambigFile && ambigRank {
				disambiguation = m.Src().String()
			} else if ambigFile {
				disambiguation = m.Src().Rank().String()
			} else {
				disambiguation = m.Src().File().String()
			}
		}
	} else {
		if m.IsCapture() {
			disambiguation = m.Src().File().String()
		}
	}
	var capture string
	if m.IsCapture() {
		capture = "x"
	}
	var promo string
	if m.Promotion().IsPromotion() {
		promo = fmt.Sprintf("=%s", strings.ToUpper(m.Promotion().PieceKind().String()))
	}
	var checkOrCheckmate string
	if isCheck {
		checkOrCheckmate = "+"
	}
	if isCheckmate {
		checkOrCheckmate = "#"
	}
	return fmt.Sprintf("%s%s%s%s%s%s", piece, disambiguation, capture, m.Dest().String(), promo, checkOrCheckmate)
}
