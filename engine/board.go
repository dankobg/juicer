package engine

// Board represents the chess board and internally uses 12 bitboards for all pieces
type Board struct {
	whiteKingOccupancy    bitboard
	whiteQueensOccupancy  bitboard
	whiteRooskOccupancy   bitboard
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

// var whitePiecesOccupancy bitboard = whiteKingOccupancy |whiteQueensOccupancy |whiteRooskOccupancy |whiteBishopsOccupancy |whiteKnightsOccupancy |whitePawnsOccupancy
// var blackPiecesOccupancy bitboard = blackKingOccupancy |blackQueensOccupancy |blackRooskOccupancy |blackBishopsOccupancy |blackKnightsOccupancy |blackPawnsOccupancy
// var allPiecesOccupancy = whitePiecesOccupancy | blackPiecesOccupancy
