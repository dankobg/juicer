package juicer

import "fmt"

type Position struct {
	board                *Board
	turn                 Color
	enpSquare            Square
	castleRights         CastleRights
	halfMoveClock        uint8
	fullMoveClock        uint16
	check                bool
	checkmate            bool
	stalemate            bool
	draw                 bool
	threeFold            bool
	insufficientMaterial bool
	terminated           bool
	outcome              string
	comments             []string
	headers              []string
	capturedPieces       []Piece
	alivePieces          []Piece
}

func (p *Position) PrintBoard() string {
	return p.board.Draw(nil)
}

func (p *Position) LoadFromFEN(fen string) error {
	meta, err := validateFEN(fen, validateFenOps{})
	if err != nil {
		return fmt.Errorf("failed to load position from fen: %w", err)
	}

	board := newEmptyBoard()

	for sq, piece := range meta.squares {
		occ := board.occupancies[piece.Color()][piece.Kind()]
		occ.setBit(sq)
		board.occupancies[piece.Color()][piece.Kind()] = occ
	}

	p.board = &board
	p.turn = meta.turnColor
	p.enpSquare = meta.enpSquare
	p.castleRights = meta.castleRights
	p.halfMoveClock = meta.halfMoveClock
	p.fullMoveClock = meta.fullMoveClock

	// p.check = false
	// p.checkmate = false
	// p.stalemate = false
	// p.draw = false
	// p.threeFold = false
	// p.insufficientMaterial = false
	// p.terminated = false
	// p.outcome = ""
	// p.comments = []string{}
	// p.headers = []string{}
	// p.capturedPieces = []Piece{}
	// p.alivePieces = []Piece{}

	return nil
}

// FenMetaPart returns the fen meta part without the position and it includes the empty string ` ` at start
// e.g. fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" -> fenMetaPart: " w KQkq - 0 1"
func (p *Position) FenMetaPart() string {
	castleToken := p.castleRights.ToFEN()

	enpSqToken := fenNoneSymbol
	if p.enpSquare != SquareNone {
		enpSqToken = p.enpSquare.Coordinate()
	}

	fenMetaPart := fmt.Sprintf(" %s %s %s %d %d", p.turn, castleToken, enpSqToken, p.halfMoveClock, p.fullMoveClock)
	return fenMetaPart
}

// Fen returns the full fen string
func (p *Position) Fen() string {
	return p.board.FenPositionPart() + p.FenMetaPart()
}

func (p *Position) whiteCanCastleKingSide() bool {
	return p.castleRights.whiteCanCastleKingSide()
}

func (p *Position) whiteCanCastleQueenSide() bool {
	return p.castleRights.whiteCanCastleQueenSide()
}

func (p *Position) whiteCanCastle() bool {
	return p.castleRights.whiteCanCastle()
}

func (p *Position) blackCanCastleKingSide() bool {
	return p.castleRights.blackCanCastleKingSide()
}

func (p *Position) blackCanCastleQueenSide() bool {
	return p.castleRights.blackCanCastleQueenSide()
}

func (p *Position) blackCanCastle() bool {
	return p.castleRights.blackCanCastle()
}

func (p *Position) canCastleKingSide() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastleKingSide()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastleKingSide()
	}
	return false
}

func (p *Position) canCastleQueenSide() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastleQueenSide()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastleQueenSide()
	}
	return false
}

func (p *Position) canCastle() bool {
	if p.turn.IsWhite() {
		return p.whiteCanCastle()
	}
	if p.turn.IsBlack() {
		return p.blackCanCastle()
	}
	return false
}

func (p *Position) generatePseudoLegalQueenMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.piecesOccupancyForSide(p.turn.Opposite())

	piece := NewPiece(Queen, p.turn)
	occupancy = p.board.occupancies[p.turn][Queen]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getQueenAttacks(src, p.board.allPiecesOccupancy()) & ^p.board.piecesOccupancyForSide(p.turn)
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalRookMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.piecesOccupancyForSide(p.turn.Opposite())

	piece := NewPiece(Rook, p.turn)
	occupancy = p.board.occupancies[p.turn][Rook]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getRookAttacks(src, p.board.allPiecesOccupancy()) & ^p.board.piecesOccupancyForSide(p.turn)
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalBishopMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.piecesOccupancyForSide(p.turn.Opposite())

	piece := NewPiece(Bishop, p.turn)
	occupancy = p.board.occupancies[p.turn][Bishop]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getBishopAttacks(src, p.board.allPiecesOccupancy()) & ^p.board.piecesOccupancyForSide(p.turn)
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalKnightMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.piecesOccupancyForSide(p.turn.Opposite())

	piece := NewPiece(Knight, p.turn)
	occupancy = p.board.occupancies[p.turn][Knight]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = knightsAttacksMask[src] & ^p.board.piecesOccupancyForSide(p.turn)
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalKingMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.piecesOccupancyForSide(p.turn.Opposite())

	piece := NewPiece(King, p.turn)
	occupancy = p.board.occupancies[p.turn][King]

	moves := make([]Move, 0)

	src = Square(occupancy.PopLS1B())
	attacks = kingAttacksMask[src] & ^p.board.piecesOccupancyForSide(p.turn)
	captures = attacks & enemies
	quiets = attacks & ^enemies

	for quiets > 0 {
		dest = Square(quiets.PopLS1B())
		moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
	}

	for captures > 0 {
		dest = Square(captures.PopLS1B())
		moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
	}

	if !p.board.IsInCheck(p.turn) {
		if p.turn.IsWhite() {
			attackedSquares := p.board.GetAttackedSquares(p.turn.Opposite(), F1G1|B1D1|C1D1, p.board.allPiecesOccupancy() & ^occupancy)

			if p.whiteCanCastleKingSide() && (p.board.allPiecesOccupancy()|attackedSquares)&F1G1 == 0 {
				moves = append(moves, NewMove(E1, G1, piece, PromotionNone, false, false, false, true))
			}
			if p.whiteCanCastleQueenSide() && p.board.allPiecesOccupancy()&B1D1 == 0 && attackedSquares&C1D1 == 0 {
				moves = append(moves, NewMove(E1, C1, piece, PromotionNone, false, false, false, true))
			}
		}

		if p.turn.IsBlack() {
			attackedSquares := p.board.GetAttackedSquares(p.turn.Opposite(), F8G8|B8D8|C8D8, p.board.allPiecesOccupancy() & ^occupancy)

			if p.blackCanCastleKingSide() && (p.board.allPiecesOccupancy()|attackedSquares)&F8G8 == 0 {
				moves = append(moves, NewMove(E8, G8, piece, PromotionNone, false, false, false, true))
			}
			if p.blackCanCastleQueenSide() && p.board.allPiecesOccupancy()&B8D8 == 0 && attackedSquares&C8D8 == 0 {
				moves = append(moves, NewMove(E8, C8, piece, PromotionNone, false, false, false, true))
			}
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalPawnMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard

	piece := NewPiece(Pawn, p.turn)
	occupancy = p.board.occupancies[p.turn][Pawn]

	moves := make([]Move, 0)

	if p.turn.IsWhite() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[White][src] & p.board.blackPiecesOccupancy()

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank7 {
					moves = append(moves,
						NewMove(src, dest, piece, PromotionQueen, true, false, false, false),
						NewMove(src, dest, piece, PromotionRook, true, false, false, false),
						NewMove(src, dest, piece, PromotionBishop, true, false, false, false),
						NewMove(src, dest, piece, PromotionKnight, true, false, false, false),
					)
				} else {
					moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
				}
			}

			dest = src + 8
			if src.Rank() == Rank7 && p.board.allPiecesOccupancy()&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				moves = append(moves,
					NewMove(src, dest, piece, PromotionQueen, false, false, false, false),
					NewMove(src, dest, piece, PromotionRook, false, false, false, false),
					NewMove(src, dest, piece, PromotionBishop, false, false, false, false),
					NewMove(src, dest, piece, PromotionKnight, false, false, false, false),
				)
			} else {
				moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
			}

			dest = src + 16
			if src.Rank() == Rank2 && p.board.allPiecesOccupancy()&(dest.occupancyMask()|Square(src+8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, true, false, false))
			}

			if p.enpSquare != SquareNone && pawnAttacksMask[White][src]&p.enpSquare.occupancyMask() != 0 {
				moves = append(moves, NewMove(src, p.enpSquare, piece, PromotionNone, true, false, true, false))
			}
		}
	}

	if p.turn.IsBlack() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[Black][src] & p.board.whitePiecesOccupancy()

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank2 {
					moves = append(moves,
						NewMove(src, dest, piece, PromotionQueen, true, false, false, false),
						NewMove(src, dest, piece, PromotionRook, true, false, false, false),
						NewMove(src, dest, piece, PromotionBishop, true, false, false, false),
						NewMove(src, dest, piece, PromotionKnight, true, false, false, false),
					)
				} else {
					moves = append(moves, NewMove(src, dest, piece, PromotionNone, true, false, false, false))
				}
			}

			dest = src - 8
			if dest >= 0 && p.board.allPiecesOccupancy()&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				if src.Rank() == Rank2 {
					moves = append(moves,
						NewMove(src, dest, piece, PromotionQueen, false, false, false, false),
						NewMove(src, dest, piece, PromotionRook, false, false, false, false),
						NewMove(src, dest, piece, PromotionBishop, false, false, false, false),
						NewMove(src, dest, piece, PromotionKnight, false, false, false, false),
					)
				} else {
					moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, false, false, false))
				}
			}

			dest = src - 16
			if src.Rank() == Rank7 && p.board.allPiecesOccupancy()&(dest.occupancyMask()|Square(src-8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, NewMove(src, dest, piece, PromotionNone, false, true, false, false))
			}

			if p.enpSquare != SquareNone && pawnAttacksMask[Black][src]&p.enpSquare.occupancyMask() != 0 {
				moves = append(moves, NewMove(src, p.enpSquare, piece, PromotionNone, true, false, true, false))
			}
		}
	}

	return moves
}

func (p *Position) generateAllPseudoLegalMoves() []Move {
	kingMoves := p.generatePseudoLegalKingMoves()
	queenMoves := p.generatePseudoLegalQueenMoves()
	rookMoves := p.generatePseudoLegalRookMoves()
	bishopMoves := p.generatePseudoLegalBishopMoves()
	knightMoves := p.generatePseudoLegalKnightMoves()
	pawnMoves := p.generatePseudoLegalPawnMoves()

	allPseudo := append(kingMoves, queenMoves...)
	allPseudo = append(allPseudo, rookMoves...)
	allPseudo = append(allPseudo, bishopMoves...)
	allPseudo = append(allPseudo, knightMoves...)
	allPseudo = append(allPseudo, pawnMoves...)

	return allPseudo
}
