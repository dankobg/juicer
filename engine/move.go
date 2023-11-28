package juicer

import "fmt"

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

func NewMove(src, dest Square, piece Piece, promotion Promotion, capture, doublePawn, enPassant, castle bool) Move {
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

func (m Move) String() string {
	promo := ""

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

	return fmt.Sprintf("%s%s%s", m.Src(), m.Dest(), promo)
}
