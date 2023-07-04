package engine

type History struct {
	move           Move
	castlingRights CastlingRights
	enpSquare      *Square
	promotion      *PromotionPieceSymbol
	halfMoves      int
	fullMoves      int
	kingsPosition  map[Color]Coordinate
}

func (h *History) String() string {
	return h.move.color.String() + ": " + h.move.String()
}

func (h *History) Move() Move {
	return h.move
}

func (h *History) CastlingRights() CastlingRights {
	return h.castlingRights
}

func (h *History) EnPassantSquare() *Square {
	return h.enpSquare
}

func (h *History) Promotion() *PromotionPieceSymbol {
	return h.promotion
}

func (h *History) HalfMoves() int {
	return h.halfMoves
}

func (h *History) FullMoves() int {
	return h.fullMoves
}

func (h *History) KingsPosition() map[Color]Coordinate {
	return h.kingsPosition
}

func NewHistory(move Move, cr CastlingRights, enpSquare *Square, promotion *PromotionPieceSymbol, halfMoves, fullMoves int, kingsPos map[Color]Coordinate) *History {
	return &History{
		move:           move,
		castlingRights: cr,
		enpSquare:      enpSquare,
		promotion:      promotion,
		halfMoves:      halfMoves,
		fullMoves:      fullMoves,
		kingsPosition:  kingsPos,
	}
}
