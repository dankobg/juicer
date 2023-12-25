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
	insufficientMaterial bool
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

	board := &Board{}

	for sq, piece := range meta.squares {
		occ := board.pieceOccupancies[piece.Color()][piece.Kind()]
		occ.setBit(sq)
		board.pieceOccupancies[piece.Color()][piece.Kind()] = occ
	}

	board.calcSideOccupancies()

	p.board = board
	p.turn = meta.turnColor
	p.enpSquare = meta.enpSquare
	p.castleRights = meta.castleRights
	p.halfMoveClock = meta.halfMoveClock
	p.fullMoveClock = meta.fullMoveClock
	p.ply = 2*(p.fullMoveClock-1) + uint16(p.turn)
	p.check = p.board.IsInCheck(p.turn)

	p.InitHash()

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

	if p.board.pieceOccupancies[White][Queen] != 0 ||
		p.board.pieceOccupancies[Black][Queen] != 0 ||
		p.board.pieceOccupancies[White][Rook] != 0 ||
		p.board.pieceOccupancies[Black][Rook] != 0 ||
		p.board.pieceOccupancies[White][Pawn] != 0 ||
		p.board.pieceOccupancies[Black][Pawn] != 0 ||
		p.board.sideOccupancies[Both].populationCount() > 4 {
		return false
	}

	wn, wb := p.board.pieceOccupancies[White][Knight].populationCount(), p.board.pieceOccupancies[White][Bishop].populationCount()
	bn, bb := p.board.pieceOccupancies[Black][Knight].populationCount(), p.board.pieceOccupancies[Black][Bishop].populationCount()
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
		return Square(p.board.pieceOccupancies[White][Bishop].LS1B()).Color() == Square(p.board.pieceOccupancies[Black][Bishop].LS1B()).Color()
	}

	return false
}

func (p *Position) generatePseudoLegalQueenMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.sideOccupancies[p.turn.Opposite()]

	piece := NewPiece(Queen, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getQueenAttacks(src, p.board.sideOccupancies[Both]) & ^p.board.sideOccupancies[p.turn]
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, newQuietMove(src, dest, piece))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, newCaptureMove(src, dest, piece))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalRookMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.sideOccupancies[p.turn.Opposite()]

	piece := NewPiece(Rook, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getRookAttacks(src, p.board.sideOccupancies[Both]) & ^p.board.sideOccupancies[p.turn]
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, newQuietMove(src, dest, piece))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, newCaptureMove(src, dest, piece))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalBishopMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.sideOccupancies[p.turn.Opposite()]

	piece := NewPiece(Bishop, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getBishopAttacks(src, p.board.sideOccupancies[Both]) & ^p.board.sideOccupancies[p.turn]
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, newQuietMove(src, dest, piece))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, newCaptureMove(src, dest, piece))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalKnightMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.sideOccupancies[p.turn.Opposite()]

	piece := NewPiece(Knight, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = knightsAttacksMask[src] & ^p.board.sideOccupancies[p.turn]
		captures = attacks & enemies
		quiets = attacks & ^enemies

		for quiets > 0 {
			dest = Square(quiets.PopLS1B())
			moves = append(moves, newQuietMove(src, dest, piece))
		}

		for captures > 0 {
			dest = Square(captures.PopLS1B())
			moves = append(moves, newCaptureMove(src, dest, piece))
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalKingMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.board.sideOccupancies[p.turn.Opposite()]

	piece := NewPiece(King, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	src = Square(occupancy.PopLS1B())
	attacks = kingAttacksMask[src] & ^p.board.sideOccupancies[p.turn]
	captures = attacks & enemies
	quiets = attacks & ^enemies

	for quiets > 0 {
		dest = Square(quiets.PopLS1B())
		moves = append(moves, newQuietMove(src, dest, piece))
	}

	for captures > 0 {
		dest = Square(captures.PopLS1B())
		moves = append(moves, newCaptureMove(src, dest, piece))
	}

	if !p.board.IsInCheck(p.turn) {
		if p.turn.IsWhite() {
			attackedSquares := p.board.GetAttackedSquares(p.turn.Opposite(), F1G1|B1D1|C1D1, p.board.sideOccupancies[Both] & ^occupancy)

			if p.whiteCanCastleKingSide() && (p.board.sideOccupancies[Both]|attackedSquares)&F1G1 == 0 {
				moves = append(moves, newCastleMove(E1, G1, piece))
			}
			if p.whiteCanCastleQueenSide() && p.board.sideOccupancies[Both]&(B1D1|C1D1) == 0 && attackedSquares&C1D1 == 0 {
				moves = append(moves, newCastleMove(E1, C1, piece))
			}
		}

		if p.turn.IsBlack() {
			attackedSquares := p.board.GetAttackedSquares(p.turn.Opposite(), F8G8|B8D8|C8D8, p.board.sideOccupancies[Both] & ^occupancy)

			if p.blackCanCastleKingSide() && (p.board.sideOccupancies[Both]|attackedSquares)&F8G8 == 0 {
				moves = append(moves, newCastleMove(E8, G8, piece))
			}
			if p.blackCanCastleQueenSide() && p.board.sideOccupancies[Both]&(B8D8|C8D8) == 0 && attackedSquares&C8D8 == 0 {
				moves = append(moves, newCastleMove(E8, C8, piece))
			}
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalPawnMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard

	piece := NewPiece(Pawn, p.turn)
	occupancy = p.board.pieceOccupancies[p.turn][piece.Kind()]

	moves := make([]Move, 0)

	if p.turn.IsWhite() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[White][src] & p.board.sideOccupancies[Black]

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank7 {
					moves = append(moves, newPossiblePromotionCaptureMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newCaptureMove(src, dest, piece))
				}
			}

			dest = src + 8
			if dest <= 63 && p.board.sideOccupancies[Both]&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				if src.Rank() == Rank7 {
					moves = append(moves, newPossiblePromotionMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newQuietMove(src, dest, piece))
				}
			}

			dest = src + 16
			if src.Rank() == Rank2 && p.board.sideOccupancies[Both]&(dest.occupancyMask()|Square(src+8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, newDoublePawnMove(src, dest, piece))
			}

			if p.enpSquare != SquareNone && pawnAttacksMask[White][src]&p.enpSquare.occupancyMask() != 0 {
				moves = append(moves, newEnpCaptureMove(src, p.enpSquare, piece))
			}
		}
	}

	if p.turn.IsBlack() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[Black][src] & p.board.sideOccupancies[White]

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank2 {
					moves = append(moves, newPossiblePromotionCaptureMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newCaptureMove(src, dest, piece))
				}
			}

			dest = src - 8
			if dest >= 0 && p.board.sideOccupancies[Both]&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				if src.Rank() == Rank2 {
					moves = append(moves, newPossiblePromotionMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newQuietMove(src, dest, piece))
				}
			}

			dest = src - 16
			if src.Rank() == Rank7 && p.board.sideOccupancies[Both]&(dest.occupancyMask()|Square(src-8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, newDoublePawnMove(src, dest, piece))
			}

			if p.enpSquare != SquareNone && pawnAttacksMask[Black][src]&p.enpSquare.occupancyMask() != 0 {
				moves = append(moves, newEnpCaptureMove(src, p.enpSquare, piece))
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
		ply:                  p.ply,
		check:                p.check,
		insufficientMaterial: p.insufficientMaterial,
		comments:             slices.Clone(p.comments),
		headers:              slices.Clone(p.headers),
		capturedPieces:       slices.Clone(p.capturedPieces),
		hash:                 p.hash,
	}

	return &positionCopy
}

// UnmakeMove undos the previous move and is returned from MakeMove func
func (p *Position) UnmakeMove() func() {
	pcopy := p.Copy()

	bcopy := pcopy.board.Copy()

	return func() {
		p.board = &bcopy
		p.turn = pcopy.turn
		p.enpSquare = pcopy.enpSquare
		p.castleRights = pcopy.castleRights
		p.halfMoveClock = pcopy.halfMoveClock
		p.fullMoveClock = pcopy.fullMoveClock
		p.ply = pcopy.ply
		p.check = pcopy.check
		p.insufficientMaterial = pcopy.insufficientMaterial
		p.comments = pcopy.comments
		p.headers = pcopy.headers
		p.capturedPieces = pcopy.capturedPieces
		p.hash = pcopy.hash
	}
}

// MakeMove makes the move and returns UnmakeMove func that undos the move
func (p *Position) MakeMove(m Move) func() {
	unmakeMove := p.UnmakeMove()

	p.ply++

	if m.IsCapture() || m.Piece().IsPawn() {
		p.halfMoveClock = 0
	} else {
		p.halfMoveClock++
	}

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
		p.ZobristCapture(m)
		p.RemoveCapturedPiece(m.Dest())
	} else if m.IsCastle() {
		if m.Piece().IsKing() {
			p.enpSquare = SquareNone
			p.ZobristMove(m)
			p.CompleteCastling(m)
		}
	} else if m.IsDoublePawn() {
		p.ZobristMove(m)
		p.enpSquare = (m.Dest() + m.Src()) / 2
		p.ZobristEnpSquare(p.enpSquare)
	} else {
		p.enpSquare = SquareNone
		p.ZobristMove(m)
	}

	piecOcc := p.board.bitboardForPiece(m.Piece())
	piecOcc.clearBit(m.Src())
	piecOcc.setBit(m.Dest())

	p.Promote(m)

	p.board.calcSideOccupancies()

	if p.turn.IsBlack() {
		p.fullMoveClock++
	}

	p.updateCastlingRights(m)

	p.ZobristTurn()
	p.SwitchTurn()
	p.check = p.board.IsInCheck(p.turn)

	return unmakeMove
}

func (p *Position) MakeNullMove() func() {
	type unmakeNullMove struct {
		enp Square
	}

	unmakeNull := unmakeNullMove{enp: p.enpSquare}

	p.enpSquare = SquareNone
	p.halfMoveClock++
	p.ply++
	p.ZobristEnpSquare(p.enpSquare)
	p.ZobristTurn()
	p.SwitchTurn()

	return func() {
		p.halfMoveClock--
		p.ply--
		p.enpSquare = unmakeNull.enp
		p.ZobristTurn()
		p.SwitchTurn()
	}
}

func (p *Position) RemoveCapturedPiece(sq Square) {
	capturedPiece := p.board.pieceAt(sq)
	if capturedPiece.IsRook() {
		if sq == A1 {
			p.castleRights.preventWhiteFromCastlingQueenSide()
		}
		if sq == H1 {
			p.castleRights.preventWhiteFromCastlingKingSide()
		}
		if sq == A8 {
			p.castleRights.preventBlackFromCastlingQueenSide()
		}
		if sq == H8 {
			p.castleRights.preventBlackFromCastlingKingSide()
		}
	}

	p.board.sideOccupancies[p.turn.Opposite()].clearBit(sq)

	for i := 0; i < len(p.board.pieceOccupancies[p.turn.Opposite()]); i++ {
		p.board.pieceOccupancies[p.turn.Opposite()][i] &= p.board.sideOccupancies[p.turn.Opposite()]
	}
}

func (p *Position) CompleteCastling(m Move) {
	var rookMove Move

	if m.Src() == E1 && m.Dest() == G1 {
		rookMove = newCastleMove(H1, F1, WhiteRook)
	}
	if m.Src() == E1 && m.Dest() == C1 {
		rookMove = newCastleMove(A1, D1, WhiteRook)
	}
	if m.Src() == E8 && m.Dest() == G8 {
		rookMove = newCastleMove(H8, F8, WhiteRook)
	}
	if m.Src() == E8 && m.Dest() == C8 {
		rookMove = newCastleMove(A8, D8, WhiteRook)
	}

	p.ZobristMove(rookMove)
	p.board.pieceOccupancies[p.turn][Rook].setBit(rookMove.Dest())
	p.board.pieceOccupancies[p.turn][Rook].clearBit(rookMove.Src())
}

// Promote promotes (replaces) a pawn on the 8th/1st rank with the promoted piece
func (p *Position) Promote(m Move) {
	if !m.Promotion().IsPromotion() {
		return
	}

	p.board.pieceOccupancies[p.turn][Pawn].clearBit(m.Dest())
	p.board.pieceOccupancies[p.turn][m.Promotion().PieceKind()].setBit(m.Dest())

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
	p.hash = defaultZobrist.seed

	for color, occupancies := range p.board.pieceOccupancies {
		for pk, occ := range occupancies {
			copy := occ

			for copy > 0 {
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

	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece().Kind()][m.Src()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn][m.Piece().Kind()][m.Dest()]
	p.hash ^= defaultZobrist.occupanciesKeys[p.turn.Opposite()][capturedPiece.Kind()][m.Dest()]
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
