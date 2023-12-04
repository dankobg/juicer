package juicer

import (
	"fmt"
	"slices"
)

type Position struct {
	board                *Board
	turn                 Color
	enpSquare            Square
	castleRights         CastleRights
	halfMoveClock        uint8
	fullMoveClock        uint16
	ply                  uint16
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
	hash                 uint64
}

func (p *Position) PrintBoard() string {
	return p.board.Draw(nil)
}

func (p *Position) LoadFromFEN(fen string) error {
	meta, err := validateFEN(fen, validateFenOps{})
	if err != nil {
		return fmt.Errorf("failed to load position from fen: %w", err)
	}

	var board Board

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

	// ply refers to a single move by one player - a full move consists of two ply e.g. 1.e4 e5
	// half move clock gets reset after each irrevirsible move played (pawn push, promotion, castle)
	// ply is basically total of half moves and doesnt get reset
	p.ply = p.fullMoveClock * 2
	if p.turn.IsBlack() {
		p.ply--
	}

	p.hash = defaultZobrist.hash
	p.InitHash()

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

func (p *Position) SwitchTurn() {
	p.turn = p.turn.Opposite()
}

// InsufficentMaterial checks if there is insufficient material on the board which leads to a draw
// theoretically possible checkmates are not counted as a draw because they can be achieved with the help of self mate
func (p *Position) IsInsufficientMaterial() bool {
	if p.board.IsOnlyKingLeft() {
		return true
	}

	if p.board.occupancies[White][Queen] != 0 ||
		p.board.occupancies[Black][Queen] != 0 ||
		p.board.occupancies[White][Rook] != 0 ||
		p.board.occupancies[Black][Rook] != 0 ||
		p.board.occupancies[White][Pawn] != 0 ||
		p.board.occupancies[Black][Pawn] != 0 ||
		p.board.allPiecesOccupancy().populationCount() > 4 {
		return false
	}

	wn, wb := p.board.occupancies[White][Knight].populationCount(), p.board.occupancies[White][Bishop].populationCount()
	bn, bb := p.board.occupancies[Black][Knight].populationCount(), p.board.occupancies[Black][Bishop].populationCount()
	wm, bm := wn+wb, bn+bb

	// k vs k+b/n (1 minor) is a draw
	if wm+bm <= 1 {
		return true
	}

	// k vs k+n+n, k vs k+b+b, k vs k+b+n is not a draw
	if wm > 1 || bm > 1 {
		return false
	}

	// same color bishops is a draw
	if wb == 1 && bb == 1 {
		return Square(p.board.occupancies[White][Bishop].LS1B()).Color() == Square(p.board.occupancies[Black][Bishop].LS1B()).Color()
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

func (p *Position) generateAllLegalMoves(pseudoMoves []Move) []Move {
	legalMoves := make([]Move, 0)

	for _, m := range pseudoMoves {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			legalMoves = append(legalMoves, m)
		}

		unmakeMove()
	}

	return legalMoves
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

func (p *Position) Copy() *Position {
	boardCopy := p.board.Copy()

	positionCopy := Position{
		board:                &boardCopy,
		turn:                 p.turn,
		enpSquare:            p.enpSquare,
		castleRights:         p.castleRights,
		halfMoveClock:        p.halfMoveClock,
		fullMoveClock:        p.fullMoveClock,
		check:                p.check,
		checkmate:            p.checkmate,
		stalemate:            p.stalemate,
		draw:                 p.draw,
		threeFold:            p.threeFold,
		insufficientMaterial: p.insufficientMaterial,
		terminated:           p.terminated,
		outcome:              p.outcome,
		comments:             slices.Clone(p.comments),
		headers:              slices.Clone(p.headers),
		capturedPieces:       slices.Clone(p.capturedPieces),
	}

	return &positionCopy
}

// UnmakeMove undos the previous move and is returned from MakeMove func
func (p *Position) UnmakeMove() func() {
	pcopy := p.Copy()

	return func() {
		p.board = pcopy.board
		p.turn = pcopy.turn
		p.enpSquare = pcopy.enpSquare
		p.castleRights = pcopy.castleRights
		p.halfMoveClock = pcopy.halfMoveClock
		p.fullMoveClock = pcopy.fullMoveClock
		p.check = pcopy.check
		p.checkmate = pcopy.checkmate
		p.stalemate = pcopy.stalemate
		p.draw = pcopy.draw
		p.threeFold = pcopy.threeFold
		p.insufficientMaterial = pcopy.insufficientMaterial
		p.terminated = pcopy.terminated
		p.outcome = pcopy.outcome
		p.comments = pcopy.comments
		p.headers = pcopy.headers
		p.capturedPieces = pcopy.capturedPieces
		p.ply = pcopy.ply
	}
}

// MakeMove makes the move and returns UnmakeMove func that undos the move
func (p *Position) MakeMove(m Move) func() {
	unmakeMove := p.UnmakeMove()

	if m.IsCapture() || m.Piece().IsPawn() {
		p.halfMoveClock = 0
	} else {
		p.halfMoveClock++
	}

	p.ply++

	if p.enpSquare != SquareNone {
		p.ZobristEnpSquare(p.enpSquare)
	}

	if m.IsEnPassant() {
		p.ZobristEnpCapture(m)
		p.enpSquare = SquareNone

		direction := 8
		if p.turn.IsBlack() {
			direction = -8
		}
		p.RemoveCapturedPiece(m.Dest() - Square(direction))
	} else if m.IsCapture() {
		p.enpSquare = SquareNone
		p.RemoveCapturedPiece(m.Dest())
		p.ZobristCapture(m)
	} else if m.IsCastle() {
		p.enpSquare = SquareNone
		p.ZobristMove(m)
		p.CompleteCastling(m)
	} else if m.IsDoublePawn() {
		p.enpSquare = (m.Dest() + m.Src()) / 2
		p.ZobristMove(m)
		p.ZobristEnpSquare(p.enpSquare)
	} else {
		p.enpSquare = SquareNone
		p.ZobristMove(m)
	}

	piecOcc := p.board.bitboardForPiece(m.Piece())
	piecOcc.clearBit(m.Src())
	piecOcc.setBit(m.Dest())

	p.Promote(m)

	if p.turn.IsBlack() {
		p.fullMoveClock++
	}

	p.updateCastlingRights(m)

	p.ZobristTurn()
	p.SwitchTurn()
	p.check = p.board.IsInCheck(p.turn)

	return unmakeMove
}

func (p *Position) RemoveCapturedPiece(sq Square) {
	enemyOcc := p.board.piecesOccupancyForSide(p.turn.Opposite())
	enemyOcc.clearBit(sq)

	for _, occ := range p.board.occupancies[p.turn.Opposite()] {
		occ &= enemyOcc
	}
}

func (p *Position) CompleteCastling(m Move) {
	occ := p.board.occupancies[p.turn][Rook]

	var rookMove Move

	if m.Src() == E1 && m.Dest() == G1 {
		rookMove = NewMove(H1, F1, WhiteRook, PromotionNone, false, false, false, false)
	}
	if m.Src() == E1 && m.Dest() == C1 {
		rookMove = NewMove(A1, D1, WhiteRook, PromotionNone, false, false, false, false)
	}
	if m.Src() == E8 && m.Dest() == G8 {
		rookMove = NewMove(H8, F8, BlackRook, PromotionNone, false, false, false, false)
	}
	if m.Src() == E8 && m.Dest() == C8 {
		rookMove = NewMove(A8, D8, BlackRook, PromotionNone, false, false, false, false)
	}

	p.ZobristMove(rookMove)
	occ.setBit(rookMove.Dest())
	occ.clearBit(rookMove.Src())
}

// Promote promotes (replaces) a pawn on the 8th/1st rank with the promoted piece
func (p *Position) Promote(m Move) {
	if m.Promotion() == PromotionNone {
		return
	}

	pawnOcc := p.board.occupancies[p.turn][Pawn]
	promoOcc := p.board.occupancies[p.turn][m.Promotion().PieceKind()]

	pawnOcc.clearBit(m.Dest())
	promoOcc.setBit(m.Dest())

	p.ZobristPromotion(m)
}

func (p *Position) updateCastlingRights(m Move) {
	if p.castleRights == 0 {
		return
	}

	if m.Piece().IsKing() {
		if p.turn.IsWhite() {
			if p.whiteCanCastleKingSide() {
				p.ZobristCastleRights(WhiteKingSideCastle)
			}
			if p.whiteCanCastleQueenSide() {
				p.ZobristCastleRights(WhiteQueenSideCastle)
			}

			p.castleRights.preventWhiteFromCastling()
		}
		if p.turn.IsBlack() {
			if p.blackCanCastleKingSide() {
				p.ZobristCastleRights(BlackKingSideCastle)
			}
			if p.blackCanCastleQueenSide() {
				p.ZobristCastleRights(BlackQueenSideCastle)
			}

			p.castleRights.preventBlackFromCastling()
		}
	}

	if m.Piece().IsRook() {
		if p.turn.IsWhite() {
			if m.Src() == H1 {
				p.castleRights.preventWhiteFromCastlingKingSide()
				p.ZobristCastleRights(WhiteKingSideCastle)
			}
			if m.Src() == A1 {
				p.castleRights.preventWhiteFromCastlingQueenSide()
				p.ZobristCastleRights(WhiteQueenSideCastle)
			}
		}

		if p.turn.IsBlack() {
			if m.Src() == H8 {
				p.castleRights.preventBlackFromCastlingKingSide()
				p.ZobristCastleRights(BlackKingSideCastle)
			}
			if m.Src() == A8 {
				p.castleRights.preventBlackFromCastlingQueenSide()
				p.ZobristCastleRights(BlackQueenSideCastle)
			}
		}
	}
}

func (p *Position) InitHash() {
	for color, occupancies := range p.board.occupancies {
		for pk, occ := range occupancies {
			copy := occ

			for occ > 0 {
				sq := copy.PopLS1B()
				p.hash ^= defaultZobrist.occupanciesKeys[color][pk][sq]
			}
		}
	}

	for ct, cr := range defaultZobrist.castleKeys {
		if p.castleRights&ct != 0 {
			p.hash ^= cr
		}
	}

	if p.enpSquare != SquareNone {
		p.hash ^= defaultZobrist.enpKeys[p.enpSquare]
	}

	if p.turn.IsBlack() {
		p.hash ^= defaultZobrist.turnKey
	}
}

func (p *Position) ZobristMove(m Move) {
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece().Kind()][m.Dest()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece().Kind()][m.Src()]
}

func (p *Position) ZobristCapture(m Move) {
	capturedPiece := p.board.pieceAt(m.Dest())

	p.hash ^= defaultZobrist.occupanciesKeys[p.turn.Opposite()][capturedPiece.Kind()][m.Dest()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece()][m.Src()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece()][m.Dest()]
}

func (p *Position) ZobristEnpCapture(m Move) {
	direction := 8
	if p.turn.IsBlack() {
		direction = -8
	}

	p.hash ^= defaultZobrist.occupanciesKeys[p.turn.Opposite()][Pawn][m.Dest()-Square(direction)]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][Pawn][m.Src()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][Pawn][m.Dest()]
}

func (p *Position) ZobristTurn() {
	p.hash ^= defaultZobrist.turnKey
}

func (p *Position) ZobristCastleRights(cr CastleRights) {
	p.hash ^= defaultZobrist.castleKeys[cr]
}

func (p *Position) ZobristPromotion(m Move) {
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Promotion().PieceKind()][m.Dest()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][Pawn][m.Dest()]
}

func (p *Position) ZobristEnpSquare(sq Square) {
	if sq != SquareNone {
		p.hash ^= defaultZobrist.enpKeys[sq]
	}
}

func (p *Position) IsThreeFoldRepetition() bool {
	// historyDepth := max(0, e.Ply-2-int(e.Board.HalfMoveCounter))
	// for ply := e.Ply - 3; ply >= historyDepth; ply -= 2 {
	// 	if e.Board.Hash == e.Plys[ply] {
	// 		return true
	// 	}
	// }

	return false
}
