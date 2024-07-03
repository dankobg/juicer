package engine

import (
	"fmt"
	"slices"
)

type Position struct {
	Board          *Board
	Turn           Color
	EpSquare       Square
	CastleRights   CastleRights
	HalfMoveClock  uint8
	FullMoveClock  uint16
	Ply            uint16
	Check          bool
	Comments       []string
	Ceaders        []string
	CapturedPieces []Piece
	Hash           uint64
}

func (p *Position) PrintBoard() string {
	return p.Board.Draw(nil)
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

	p.Board = board
	p.Turn = meta.turnColor
	p.EpSquare = meta.enpSquare
	p.CastleRights = meta.castleRights
	p.HalfMoveClock = meta.halfMoveClock
	p.FullMoveClock = meta.fullMoveClock
	p.Ply = 2*(p.FullMoveClock-1) + uint16(p.Turn)
	p.Check = p.Board.IsInCheck(p.Turn)

	p.InitHash()

	return nil
}

// FenMetaPart returns the fen meta part without the position and it includes the empty string ` ` at start
// e.g. fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" -> fenMetaPart: " w KQkq - 0 1"
func (p *Position) FenMetaPart() string {
	castleToken := p.CastleRights.ToFEN()

	enpSqToken := fenNoneSymbol
	if p.EpSquare != SquareNone {
		enpSqToken = p.EpSquare.Coordinate()
	}

	fenMetaPart := fmt.Sprintf(" %s %s %s %d %d", p.Turn, castleToken, enpSqToken, p.HalfMoveClock, p.FullMoveClock)
	return fenMetaPart
}

// Fen returns the full fen string
func (p *Position) Fen() string {
	return p.Board.FenPositionPart() + p.FenMetaPart()
}

func (p *Position) whiteHasKingSideCastleRights() bool {
	return p.CastleRights.whiteHasKingSideCastleRights()
}

func (p *Position) whiteHasQueenSideCastleRights() bool {
	return p.CastleRights.whiteHasQueenSideCastleRights()
}

func (p *Position) whiteHasCastleRights() bool {
	return p.CastleRights.whiteHasCastleRights()
}

func (p *Position) blackHasKingSideCastleRights() bool {
	return p.CastleRights.blackHasKingSideCastleRights()
}

func (p *Position) blackHasQueenSideCastleRights() bool {
	return p.CastleRights.blackHasQueenSideCastleRights()
}

func (p *Position) blackHasCastleRights() bool {
	return p.CastleRights.blackHasCastleRights()
}

func (p *Position) hasKingSideCastleRights() bool {
	if p.Turn.IsWhite() {
		return p.whiteHasKingSideCastleRights()
	}
	if p.Turn.IsBlack() {
		return p.blackHasKingSideCastleRights()
	}
	return false
}

func (p *Position) hasQueenSideCastleRights() bool {
	if p.Turn.IsWhite() {
		return p.whiteHasQueenSideCastleRights()
	}
	if p.Turn.IsBlack() {
		return p.blackHasQueenSideCastleRights()
	}
	return false
}

func (p *Position) hasCastleRights() bool {
	if p.Turn.IsWhite() {
		return p.whiteHasCastleRights()
	}
	if p.Turn.IsBlack() {
		return p.blackHasCastleRights()
	}
	return false
}

func (p *Position) SwitchTurn() {
	p.Turn = p.Turn.Opposite()
}

// InsufficentMaterial checks if there is insufficient material on the board which leads to a draw
// theoretically possible checkmates are not counted as a draw because they can be achieved with the help of self mate
func (p *Position) IsInsufficientMaterial() bool {
	if p.Board.IsOnlyKingLeft() {
		return true
	}

	if p.Board.pieceOccupancies[White][Queen] != 0 ||
		p.Board.pieceOccupancies[Black][Queen] != 0 ||
		p.Board.pieceOccupancies[White][Rook] != 0 ||
		p.Board.pieceOccupancies[Black][Rook] != 0 ||
		p.Board.pieceOccupancies[White][Pawn] != 0 ||
		p.Board.pieceOccupancies[Black][Pawn] != 0 ||
		p.Board.sideOccupancies[Both].populationCount() > 4 {
		return false
	}

	wn, wb := p.Board.pieceOccupancies[White][Knight].populationCount(), p.Board.pieceOccupancies[White][Bishop].populationCount()
	bn, bb := p.Board.pieceOccupancies[Black][Knight].populationCount(), p.Board.pieceOccupancies[Black][Bishop].populationCount()
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
		return Square(p.Board.pieceOccupancies[White][Bishop].LS1B()).Color() == Square(p.Board.pieceOccupancies[Black][Bishop].LS1B()).Color()
	}

	return false
}

func (p *Position) generatePseudoLegalQueenMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard
	var quiets, captures bitboard
	enemies := p.Board.sideOccupancies[p.Turn.Opposite()]

	piece := NewPiece(Queen, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getQueenAttacks(src, p.Board.sideOccupancies[Both]) & ^p.Board.sideOccupancies[p.Turn]
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
	enemies := p.Board.sideOccupancies[p.Turn.Opposite()]

	piece := NewPiece(Rook, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getRookAttacks(src, p.Board.sideOccupancies[Both]) & ^p.Board.sideOccupancies[p.Turn]
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
	enemies := p.Board.sideOccupancies[p.Turn.Opposite()]

	piece := NewPiece(Bishop, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = getBishopAttacks(src, p.Board.sideOccupancies[Both]) & ^p.Board.sideOccupancies[p.Turn]
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
	enemies := p.Board.sideOccupancies[p.Turn.Opposite()]

	piece := NewPiece(Knight, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	for occupancy > 0 {
		src = Square(occupancy.PopLS1B())
		attacks = knightsAttacksMask[src] & ^p.Board.sideOccupancies[p.Turn]
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
	enemies := p.Board.sideOccupancies[p.Turn.Opposite()]

	piece := NewPiece(King, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	src = Square(occupancy.PopLS1B())
	attacks = kingAttacksMask[src] & ^p.Board.sideOccupancies[p.Turn]
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

	if !p.Board.IsInCheck(p.Turn) {
		if p.Turn.IsWhite() {
			attackedSquares := p.Board.GetAttackedSquares(p.Turn.Opposite(), F1G1|B1D1|C1D1, p.Board.sideOccupancies[Both] & ^occupancy)

			if p.whiteHasKingSideCastleRights() && (p.Board.sideOccupancies[Both]|attackedSquares)&F1G1 == 0 {
				moves = append(moves, newCastleMove(E1, G1, piece))
			}
			if p.whiteHasQueenSideCastleRights() && p.Board.sideOccupancies[Both]&(B1D1|C1D1) == 0 && attackedSquares&C1D1 == 0 {
				moves = append(moves, newCastleMove(E1, C1, piece))
			}
		}

		if p.Turn.IsBlack() {
			attackedSquares := p.Board.GetAttackedSquares(p.Turn.Opposite(), F8G8|B8D8|C8D8, p.Board.sideOccupancies[Both] & ^occupancy)

			if p.blackHasKingSideCastleRights() && (p.Board.sideOccupancies[Both]|attackedSquares)&F8G8 == 0 {
				moves = append(moves, newCastleMove(E8, G8, piece))
			}
			if p.blackHasQueenSideCastleRights() && p.Board.sideOccupancies[Both]&(B8D8|C8D8) == 0 && attackedSquares&C8D8 == 0 {
				moves = append(moves, newCastleMove(E8, C8, piece))
			}
		}
	}

	return moves
}

func (p *Position) generatePseudoLegalPawnMoves() []Move {
	var src, dest Square
	var occupancy, attacks bitboard

	piece := NewPiece(Pawn, p.Turn)
	occupancy = p.Board.pieceOccupancies[p.Turn][piece.Kind()]

	moves := make([]Move, 0)

	if p.Turn.IsWhite() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[White][src] & p.Board.sideOccupancies[Black]

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank7 {
					moves = append(moves, newPossiblePromotionCaptureMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newCaptureMove(src, dest, piece))
				}
			}

			dest = src + 8
			if dest <= 63 && p.Board.sideOccupancies[Both]&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				if src.Rank() == Rank7 {
					moves = append(moves, newPossiblePromotionMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newQuietMove(src, dest, piece))
				}
			}

			dest = src + 16
			if src.Rank() == Rank2 && p.Board.sideOccupancies[Both]&(dest.occupancyMask()|Square(src+8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, newDoublePawnMove(src, dest, piece))
			}

			if p.EpSquare != SquareNone && pawnAttacksMask[White][src]&p.EpSquare.occupancyMask() != 0 {
				moves = append(moves, newEnpCaptureMove(src, p.EpSquare, piece))
			}
		}
	}

	if p.Turn.IsBlack() {
		for occupancy > 0 {
			src = Square(occupancy.PopLS1B())
			attacks = pawnAttacksMask[Black][src] & p.Board.sideOccupancies[White]

			for attacks > 0 {
				dest = Square(attacks.PopLS1B())

				if src.Rank() == Rank2 {
					moves = append(moves, newPossiblePromotionCaptureMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newCaptureMove(src, dest, piece))
				}
			}

			dest = src - 8
			if dest >= 0 && p.Board.sideOccupancies[Both]&dest.occupancyMask() == 0 && dest.occupancyMask() != 0 {
				if src.Rank() == Rank2 {
					moves = append(moves, newPossiblePromotionMoves(src, dest, piece)...)
				} else {
					moves = append(moves, newQuietMove(src, dest, piece))
				}
			}

			dest = src - 16
			if src.Rank() == Rank7 && p.Board.sideOccupancies[Both]&(dest.occupancyMask()|Square(src-8).occupancyMask()) == 0 && dest.occupancyMask() != 0 {
				moves = append(moves, newDoublePawnMove(src, dest, piece))
			}

			if p.EpSquare != SquareNone && pawnAttacksMask[Black][src]&p.EpSquare.occupancyMask() != 0 {
				moves = append(moves, newEnpCaptureMove(src, p.EpSquare, piece))
			}
		}
	}

	return moves
}

func (p *Position) generateAllLegalMoves(pseudoMoves []Move) []Move {
	legalMoves := make([]Move, 0)

	for _, m := range pseudoMoves {
		unmakeMove := p.MakeMove(m)

		if !p.Board.IsInCheck(p.Turn.Opposite()) {
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
	boardCopy := p.Board.Copy()

	positionCopy := Position{
		Board:          &boardCopy,
		Turn:           p.Turn,
		EpSquare:       p.EpSquare,
		CastleRights:   p.CastleRights,
		HalfMoveClock:  p.HalfMoveClock,
		FullMoveClock:  p.FullMoveClock,
		Ply:            p.Ply,
		Check:          p.Check,
		Comments:       slices.Clone(p.Comments),
		Ceaders:        slices.Clone(p.Ceaders),
		CapturedPieces: slices.Clone(p.CapturedPieces),
		Hash:           p.Hash,
	}

	return &positionCopy
}

// UnmakeMove undos the previous move and is returned from MakeMove func
func (p *Position) UnmakeMove() func() {
	pcopy := p.Copy()

	bcopy := pcopy.Board.Copy()

	return func() {
		p.Board = &bcopy
		p.Turn = pcopy.Turn
		p.EpSquare = pcopy.EpSquare
		p.CastleRights = pcopy.CastleRights
		p.HalfMoveClock = pcopy.HalfMoveClock
		p.FullMoveClock = pcopy.FullMoveClock
		p.Ply = pcopy.Ply
		p.Check = pcopy.Check
		p.Comments = pcopy.Comments
		p.Ceaders = pcopy.Ceaders
		p.CapturedPieces = pcopy.CapturedPieces
		p.Hash = pcopy.Hash
	}
}

// MakeMove makes the move and returns UnmakeMove func that undos the move
func (p *Position) MakeMove(m Move) func() {
	unmakeMove := p.UnmakeMove()

	p.Ply++

	if m.IsCapture() || m.Piece().IsPawn() {
		p.HalfMoveClock = 0
	} else {
		p.HalfMoveClock++
	}

	if p.EpSquare != SquareNone {
		p.ZobristEnpSquare(p.EpSquare)
	}

	if m.IsEnPassant() {
		p.ZobristEnpCapture(m)
		p.EpSquare = SquareNone

		direction := 8
		if p.Turn.IsBlack() {
			direction = -8
		}
		p.RemoveCapturedPiece(m.Dest() - Square(direction))
	} else if m.IsCapture() {
		p.EpSquare = SquareNone
		p.ZobristCapture(m)
		p.RemoveCapturedPiece(m.Dest())
	} else if m.IsCastle() {
		if m.Piece().IsKing() {
			p.EpSquare = SquareNone
			p.ZobristMove(m)
			p.CompleteCastling(m)
		}
	} else if m.IsDoublePawn() {
		p.ZobristMove(m)
		p.EpSquare = (m.Dest() + m.Src()) / 2
		p.ZobristEnpSquare(p.EpSquare)
	} else {
		p.EpSquare = SquareNone
		p.ZobristMove(m)
	}

	piecOcc := p.Board.bitboardForPiece(m.Piece())
	piecOcc.clearBit(m.Src())
	piecOcc.setBit(m.Dest())

	p.Promote(m)

	p.Board.calcSideOccupancies()

	if p.Turn.IsBlack() {
		p.FullMoveClock++
	}

	p.updateCastlingRights(m)

	p.ZobristTurn()
	p.SwitchTurn()
	p.Check = p.Board.IsInCheck(p.Turn)

	return unmakeMove
}

func (p *Position) MakeNullMove() func() {
	type unmakeNullMove struct {
		enp Square
	}

	unmakeNull := unmakeNullMove{enp: p.EpSquare}

	p.EpSquare = SquareNone
	p.HalfMoveClock++
	p.Ply++
	p.ZobristEnpSquare(p.EpSquare)
	p.ZobristTurn()
	p.SwitchTurn()

	return func() {
		p.HalfMoveClock--
		p.Ply--
		p.EpSquare = unmakeNull.enp
		p.ZobristTurn()
		p.SwitchTurn()
	}
}

func (p *Position) RemoveCapturedPiece(sq Square) {
	capturedPiece := p.Board.pieceAt(sq)
	p.CapturedPieces = append(p.CapturedPieces, capturedPiece)

	if capturedPiece.IsRook() {
		if sq == A1 {
			p.CastleRights.disableWhiteQueenSideCastleRight()
		}
		if sq == H1 {
			p.CastleRights.disableWhiteKingSideCastleRight()
		}
		if sq == A8 {
			p.CastleRights.disableBlackQueenSideCastleRight()
		}
		if sq == H8 {
			p.CastleRights.disableBlackKingSideCastleRight()
		}
	}

	p.Board.sideOccupancies[p.Turn.Opposite()].clearBit(sq)

	for i := 0; i < len(p.Board.pieceOccupancies[p.Turn.Opposite()]); i++ {
		p.Board.pieceOccupancies[p.Turn.Opposite()][i] &= p.Board.sideOccupancies[p.Turn.Opposite()]
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
	p.Board.pieceOccupancies[p.Turn][Rook].setBit(rookMove.Dest())
	p.Board.pieceOccupancies[p.Turn][Rook].clearBit(rookMove.Src())
}

// Promote promotes (replaces) a pawn on the 8th/1st rank with the promoted piece
func (p *Position) Promote(m Move) {
	if !m.Promotion().IsPromotion() {
		return
	}

	p.Board.pieceOccupancies[p.Turn][Pawn].clearBit(m.Dest())
	p.Board.pieceOccupancies[p.Turn][m.Promotion().PieceKind()].setBit(m.Dest())

	p.ZobristPromotion(m)
}

func (p *Position) updateCastlingRights(m Move) {
	if p.CastleRights == 0 {
		return
	}

	if m.Piece().IsKing() {
		if p.Turn.IsWhite() {
			if p.whiteHasKingSideCastleRights() {
				p.ZobristCastleRights(WhiteKingSideCastle)
			}
			if p.whiteHasQueenSideCastleRights() {
				p.ZobristCastleRights(WhiteQueenSideCastle)
			}

			p.CastleRights.disableWhiteCastleRights()
		}
		if p.Turn.IsBlack() {
			if p.blackHasKingSideCastleRights() {
				p.ZobristCastleRights(BlackKingSideCastle)
			}
			if p.blackHasQueenSideCastleRights() {
				p.ZobristCastleRights(BlackQueenSideCastle)
			}

			p.CastleRights.disableBlackCastleRights()
		}
	}

	if m.Piece().IsRook() {
		if p.Turn.IsWhite() {
			if m.Src() == H1 {
				p.CastleRights.disableWhiteKingSideCastleRight()
				p.ZobristCastleRights(WhiteKingSideCastle)
			}
			if m.Src() == A1 {
				p.CastleRights.disableWhiteQueenSideCastleRight()
				p.ZobristCastleRights(WhiteQueenSideCastle)
			}
		}

		if p.Turn.IsBlack() {
			if m.Src() == H8 {
				p.CastleRights.disableBlackKingSideCastleRight()
				p.ZobristCastleRights(BlackKingSideCastle)
			}
			if m.Src() == A8 {
				p.CastleRights.disableBlackQueenSideCastleRight()
				p.ZobristCastleRights(BlackQueenSideCastle)
			}
		}
	}
}

func (p *Position) InitHash() {
	p.Hash = defaultZobrist.seed

	for color, occupancies := range p.Board.pieceOccupancies {
		for pk, occ := range occupancies {
			copy := occ

			for copy > 0 {
				sq := copy.PopLS1B()
				p.Hash ^= defaultZobrist.occupanciesKeys[color][pk][sq]
			}
		}
	}

	for ct, cr := range defaultZobrist.castleKeys {
		if p.CastleRights&ct != 0 {
			p.Hash ^= cr
		}
	}

	if p.EpSquare != SquareNone {
		p.Hash ^= defaultZobrist.enpKeys[p.EpSquare]
	}

	if p.Turn.IsBlack() {
		p.Hash ^= defaultZobrist.turnKey
	}
}

func (p *Position) ZobristMove(m Move) {
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][m.Piece().Kind()][m.Dest()]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][m.Piece().Kind()][m.Src()]
}

func (p *Position) ZobristCapture(m Move) {
	capturedPiece := p.Board.pieceAt(m.Dest())

	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][m.Piece().Kind()][m.Src()]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][m.Piece().Kind()][m.Dest()]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn.Opposite()][capturedPiece.Kind()][m.Dest()]
}

func (p *Position) ZobristEnpCapture(m Move) {
	direction := 8
	if p.Turn.IsBlack() {
		direction = -8
	}

	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn.Opposite()][Pawn][m.Dest()-Square(direction)]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][Pawn][m.Src()]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][Pawn][m.Dest()]
}

func (p *Position) ZobristTurn() {
	p.Hash ^= defaultZobrist.turnKey
}

func (p *Position) ZobristCastleRights(cr CastleRights) {
	p.Hash ^= defaultZobrist.castleKeys[cr]
}

func (p *Position) ZobristPromotion(m Move) {
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][m.Promotion().PieceKind()][m.Dest()]
	p.Hash ^= defaultZobrist.occupanciesKeys[p.Turn][Pawn][m.Dest()]
}

func (p *Position) ZobristEnpSquare(sq Square) {
	if sq != SquareNone {
		p.Hash ^= defaultZobrist.enpKeys[sq]
	}
}
