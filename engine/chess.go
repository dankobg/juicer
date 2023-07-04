package engine

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

type Color string

const (
	White Color = "w"
	Black Color = "b"
)

func (c Color) String() string {
	return string(c)
}

func (c Color) Name() string {
	if c == White {
		return "White"
	}

	return "Black"
}

func (c Color) Equals(color Color) bool {
	return c.String() == color.String()
}

func (c Color) IsWhite() bool {
	return c.String() == White.String()
}

func (c Color) IsBlack() bool {
	return c.String() == Black.String()
}

type CastlingRights int

const (
	WhiteKingSideCastle CastlingRights = 1 << iota
	WhiteQueenSideCastle
	BlackKingSideCastle
	BlackQueenSideCastle
)

var startingCastlingRights = WhiteKingSideCastle | WhiteQueenSideCastle | BlackKingSideCastle | BlackQueenSideCastle

const (
	kingSideCastle  string = "k"
	queenSideCastle string = "q"
)

const (
	boardSize           = 8
	FENStartingPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	FENEmptyPosition    = "8/8/8/8/8/8/8/8"
)

var (
	startingWhiteKingPos = Coordinate("e1")
	startingBlackKingPos = Coordinate("e8")
	startingKingsPos     = map[Color]Coordinate{"w": startingWhiteKingPos, "b": startingBlackKingPos}
)

type Chess struct {
	board                *Board
	turn                 Color
	halfMoves            int
	fullMoves            int
	legalMoves           []*Move
	history              []*History
	enpSquare            *Square
	kingsPos             map[Color]Coordinate
	castlingRights       CastlingRights
	comments             []string
	headers              []string
	alivePieces          []Piece
	capturedPieces       []Piece
	check                bool
	checkmate            bool
	stalemate            bool
	draw                 bool
	threeFold            bool
	insufficientMaterial bool
	terminated           bool
	outcome              string
}

func (c *Chess) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%v\n", c.DrawBoard()))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("FEN: %v\n", c.FEN()))
	b.WriteString(fmt.Sprintf("Turn: %v\n", c.Turn()))
	b.WriteString(fmt.Sprintf("Castle: %v\n", c.CastlingRights()))
	b.WriteString(fmt.Sprintf("Kings: %v\n", c.KingsPos()))
	b.WriteString(fmt.Sprintf("Full moves: %v\n", c.FullMoves()))
	b.WriteString(fmt.Sprintf("Half moves: %v\n", c.HalfMoves()))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("Check: %v\n", c.Check()))
	b.WriteString(fmt.Sprintf("Stalemate: %v\n", c.Stalemate()))
	b.WriteString(fmt.Sprintf("Checkmate: %v\n", c.Checkmate()))
	b.WriteString(fmt.Sprintf("Threefold: %v\n", c.ThreeFoldRepetition()))
	b.WriteString(fmt.Sprintf("Insufficient: %v\n", c.InsufficientMaterial()))
	b.WriteString(fmt.Sprintf("Draw: %v\n", c.Draw()))
	b.WriteString(fmt.Sprintf("Terminated: %v\n", c.Terminated()))
	b.WriteString(fmt.Sprintf("Outcome: %v\n", c.Outcome()))

	if len(c.history) > 0 {
		b.WriteString("\nHistory: \n")

		for _, h := range c.History() {
			b.WriteString(fmt.Sprintf("%v\n", h.String()))
		}
	}

	b.WriteString(fmt.Sprintf("\nLegal Moves (%d): \n", len(c.LegalMoves())))
	for i, m := range c.LegalMoves() {

		if i == len(c.LegalMoves())-1 {
			b.WriteString(fmt.Sprintf("%v", m.String()))
		} else {
			b.WriteString(fmt.Sprintf("%v\n", m.String()))
		}
	}

	return b.String()
}

func (c *Chess) Board() *Board {
	return c.board
}

func (c *Chess) Turn() Color {
	return c.turn
}

func (c *Chess) HalfMoves() int {
	return c.halfMoves
}

func (c *Chess) FullMoves() int {
	return c.fullMoves
}

func (c *Chess) LegalMoves() []*Move {
	return c.legalMoves
}

func (c *Chess) History() []*History {
	return c.history
}

func (c *Chess) EnpSquare() *Square {
	return c.enpSquare
}

func (c *Chess) KingsPos() map[Color]Coordinate {
	return c.kingsPos
}

func (c *Chess) CurrentKingPos() Coordinate {
	return c.kingsPos[c.turn]
}

func (c *Chess) CastlingRights() CastlingRights {
	return c.castlingRights
}

func (c *Chess) Comments() []string {
	return c.comments
}

func (c *Chess) Headers() []string {
	return c.headers
}

func (c *Chess) AlivePieces() []Piece {
	return c.alivePieces
}

func (c *Chess) CapturedPieces() []Piece {
	return c.capturedPieces
}

func (c *Chess) Check() bool {
	return c.check
}

func (c *Chess) Checkmate() bool {
	return c.checkmate
}

func (c *Chess) Stalemate() bool {
	return c.stalemate
}

func (c *Chess) ThreeFoldRepetition() bool {
	return c.threeFold
}

func (c *Chess) InsufficientMaterial() bool {
	return c.insufficientMaterial
}

func (c *Chess) Draw() bool {
	return c.draw
}

func (c *Chess) Terminated() bool {
	return c.terminated
}

func (c *Chess) Outcome() string {
	return c.outcome
}

func (c *Chess) isInsufficientMaterial() bool {
	// k vs k
	if len(c.alivePieces) == 2 {
		return true
	} else if len(c.alivePieces) == 3 {
		// k vs k+n or k vs k+b
		if ok := slices.ContainsFunc(c.alivePieces, func(p Piece) bool {
			if p.IsBishop() || p.IsKnight() {
				return true
			}
			return false
		}); ok {
			return true
		}
	} else if len(c.alivePieces) == 4 {
		pieces := map[PieceFENSymbol]int{}

		for _, p := range c.alivePieces {
			fenSymb := p.ToFENSymbol()
			c, ok := pieces[fenSymb]
			if ok {
				pieces[fenSymb] = c + 1
			} else {
				pieces[fenSymb] = 1
			}
		}

		wn, bn, wb, bb := pieces[WhiteKnight], pieces[BlackKnight], pieces[WhiteBishop], pieces[BlackBishop]

		// k vs k+n+n
		if wn == 2 || bn == 2 {
			return true
		}

		// k+n vs k+n
		if wn == 1 && bn == 1 {
			return true
		}

		// k+b vs k+b
		if wb == 1 && bb == 1 {
			return true
		}
	}

	return false
}

func (c *Chess) IsThreeFoldRepetition() bool {
	positions := make(map[string]int)
	moves := make([]Move, 0)

	var threeFold bool

	for {
		if len(c.history) == 0 {
			break
		}

		lastMove := c.history[len(c.history)-1].move
		moves = append(moves, lastMove)
	}

	for {
		// last 2 fields are not needed to check for three fold repetition
		fen := c.FEN()[:2]

		if count, ok := positions[fen]; ok {
			positions[fen] = count + 1
		} else {
			positions[fen] = 1
		}

		if count, ok := positions[fen]; ok && count >= 3 {
			threeFold = true
		}

		if len(moves) == 0 {
			break
		}

		m := moves[len(moves)-1]
		moves = moves[:len(moves)-1]
		c.makeMoveInternal(m)
	}

	return threeFold
}

func NewChess(fen string) (*Chess, error) {
	c := &Chess{
		castlingRights: startingCastlingRights,
		kingsPos:       startingKingsPos,
		fullMoves:      1,
		turn:           White,
	}

	fenStr := FENStartingPosition
	if fen != "" {
		fenStr = fen
	}

	meta, err := ParseMetadataFromFEN(fenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata from FEN: %w", err)
	}

	board, err := NewBoard(fenStr)
	if err != nil {
		return nil, err
	}

	c.board = board
	c.turn = meta.turn
	c.castlingRights = meta.castlingRights
	c.enpSquare = meta.enpSquare
	c.halfMoves = int(meta.halfMoves)
	c.fullMoves = int(meta.fullMoves)

	c.generateAllStrictlyLegalMoves()
	c.updateGameTerminationState()
	c.getAllAlivePieces()

	return c, nil
}

func (c *Chess) FEN() string {
	meta := Metadata{
		board:          c.board.squares,
		turn:           c.turn,
		castlingRights: c.castlingRights,
		enpSquare:      c.enpSquare,
		halfMoves:      uint8(c.halfMoves),
		fullMoves:      uint8(c.fullMoves),
	}

	return c.board.FEN(meta)
}

func (c *Chess) DrawBoard() string {
	return c.board.Draw()
}

func (c *Chess) getAllAlivePieces() {
	for i := 0; i < len(c.board.squares); i++ {
		for j := 0; j < len(c.board.squares[i]); j++ {
			if c.board.squares[i][j].HasPiece() && c.board.squares[i][j].piece.alive {
				c.alivePieces = append(c.alivePieces, *c.board.squares[i][j].piece)
			}
		}
	}
}

func (c *Chess) updateCastlingRights(move Move) {
	if move.IsWhite() && move.FromRow() == 7 {
		if move.pieceMoved.IsKing() {
			c.castlingRights &= ^(WhiteKingSideCastle | WhiteQueenSideCastle)
		}

		if move.pieceMoved.IsRook() {
			if move.FromColumn() == 0 {
				c.castlingRights &= ^WhiteQueenSideCastle
			}
			if move.FromColumn() == 7 {
				c.castlingRights &= ^WhiteKingSideCastle
			}
		}
	}

	if move.IsBlack() && move.fromSquare.row == 0 {
		if move.pieceMoved.IsKing() {
			c.castlingRights &= ^(BlackKingSideCastle | BlackQueenSideCastle)
		}

		if move.pieceMoved.IsRook() {
			if move.FromColumn() == 0 {
				c.castlingRights &= ^BlackQueenSideCastle
			}
			if move.FromColumn() == 7 {
				c.castlingRights &= ^BlackKingSideCastle
			}
		}
	}
}

func (c *Chess) enPassantPossible() bool {
	return c.enpSquare != nil
}

func (c *Chess) whiteToMove() bool {
	return c.turn.IsWhite()
}

func (c *Chess) blackToMove() bool {
	return c.turn.IsBlack()
}

func (c *Chess) switchTurn() {
	c.turn = swapColor(c.turn)
}

func (c *Chess) whiteCanCastleKingSide() bool {
	return (c.castlingRights & WhiteKingSideCastle) > 0
}

func (c *Chess) whiteCanCastleQueenSide() bool {
	return (c.castlingRights & WhiteQueenSideCastle) > 0
}

func (c *Chess) blackCanCastleKingSide() bool {
	return (c.castlingRights & BlackKingSideCastle) > 0
}

func (c *Chess) blackCanCastleQueenSide() bool {
	return (c.castlingRights & BlackQueenSideCastle) > 0
}

func (c *Chess) WhiteCanCastle() bool {
	return c.whiteCanCastleKingSide() || c.whiteCanCastleQueenSide()
}

func (c *Chess) BlackCanCastle() bool {
	return c.blackCanCastleKingSide() || c.blackCanCastleQueenSide()
}

func (c *Chess) resetKingsPos() {
	c.kingsPos = startingKingsPos
}

func (c *Chess) resetMovesCount() {
	c.halfMoves = 0
	c.fullMoves = 1
}

func (c *Chess) resetActiveTurn() {
	c.turn = White
}

func (c *Chess) resetLegalMoves() {
	// @TODO: maybe generate again??
	c.legalMoves = make([]*Move, 0)
}

func (c *Chess) resetHistory() {
	c.history = make([]*History, 0)
}

func (c *Chess) resetCastlingRights() {
	c.castlingRights = startingCastlingRights
}

func (c *Chess) Reset() error {
	if _, err := c.board.LoadFromFEN(""); err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	c.resetCastlingRights()
	c.resetKingsPos()
	c.resetActiveTurn()
	c.resetMovesCount()
	c.resetHistory()
	c.resetLegalMoves()

	return nil
}

func (c *Chess) updateRookPositionAfterCastling(move Move) {
	isCastle, castleSide := move.IsCastle(c.castlingRights)

	if !isCastle {
		return
	}

	if castleSide == kingSideCastle {
		tmpRook := c.board.squares[move.ToRow()][move.ToColumn()+1].piece
		c.board.squares[move.ToRow()][move.ToColumn()+1].piece = nil
		c.board.squares[move.ToRow()][move.ToColumn()-1].piece = tmpRook
	}

	if castleSide == queenSideCastle {
		tmpRook := c.board.squares[move.ToRow()][move.ToColumn()-2].piece
		c.board.squares[move.ToRow()][move.ToColumn()-2].piece = nil
		c.board.squares[move.ToRow()][move.ToColumn()+1].piece = tmpRook
	}
}

func (c *Chess) updateEnPassantSquare(move Move) {
	if move.pieceMoved.IsPawn() && math.Abs(float64(move.FromRow())-float64(move.ToRow())) == 2 {
		enpRow := (move.FromRow() + move.ToRow()) / 2
		enpCol := move.FromColumn()
		c.enpSquare = c.board.squares[enpRow][enpCol]
	} else {
		c.enpSquare = nil
	}
}

func (c *Chess) generateCastlingMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	kr, kc := c.CurrentKingPos().RowCol()

	// white castling moves
	if c.whiteToMove() {
		// white king side castle
		if c.whiteCanCastleKingSide() {
			if c.board.squares[kr][kc+1].IsEmpty() && c.board.squares[kr][kc+2].IsEmpty() {
				toSquare := c.board.squares[kr][kc+2]
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, WhiteKingSideCastle, nil)
				moves = append(moves, move)
			}
		}

		// white queen side castle
		if c.whiteCanCastleQueenSide() {
			if c.board.squares[kr][kc-1].IsEmpty() && c.board.squares[kr][kc-2].IsEmpty() && c.board.squares[kr][kc-3].IsEmpty() {
				toSquare := c.board.squares[kr][kc-2]
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, WhiteQueenSideCastle, nil)
				moves = append(moves, move)
			}
		}
	}

	// black castling moves
	if c.blackToMove() {
		// black king side castle
		if c.blackCanCastleKingSide() {
			if c.board.squares[kr][kc+1].IsEmpty() && c.board.squares[kr][kc+2].IsEmpty() {
				toSquare := c.board.squares[kr][kc+2]
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, BlackKingSideCastle, nil)
				moves = append(moves, move)
			}
		}

		// black queen side castle
		if c.blackCanCastleQueenSide() {
			if c.board.squares[kr][kc-1].IsEmpty() && c.board.squares[kr][kc-2].IsEmpty() && c.board.squares[kr][kc-3].IsEmpty() {
				toSquare := c.board.squares[kr][kc-2]
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, BlackQueenSideCastle, nil)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (c *Chess) generatePawnPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	row, col := fromSquare.RowCol()

	// white pawn moves
	if c.whiteToMove() {
		if nr := row - 1; validRowCol(nr) && c.board.squares[nr][col].IsEmpty() {
			// white pawn 1 square forward move
			toSquare := NewSquare(nr, col, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)

			if nr := row - 2; validRowCol(nr) && row == 6 && c.board.squares[nr][col].IsEmpty() {
				toSquare := NewSquare(nr, col, nil)
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
				moves = append(moves, move)
			}
		}

		// white pawn capture to the left
		if nr, nc := row-1, col-1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].HasEnemyPiece(c.turn) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// white pawn capture to the right
		if nr, nc := row-1, col+1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].HasEnemyPiece(c.turn) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// white en-passant to the left
		if nr, nc := row-1, col-1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].IsEmpty() && c.enPassantPossible() && c.board.squares[nr][nc].CoordEquals(*c.enpSquare) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// white en-passant to the right
		if nr, nc := row-1, col+1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].IsEmpty() && c.enPassantPossible() && c.board.squares[nr][nc].CoordEquals(*c.enpSquare) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}
	}

	// black pawn moves
	if c.blackToMove() {
		if nr := row + 1; validRowCol(uint8(nr)) && c.board.squares[nr][col].IsEmpty() {
			// black pawn 1 square forward move
			toSquare := NewSquare(nr, col, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)

			// black pawn 2 squares forward move
			if nr := row + 2; validRowCol(nr) && row == 1 && c.board.squares[nr][col].IsEmpty() {
				toSquare := NewSquare(nr, col, nil)
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
				moves = append(moves, move)
			}
		}

		// black pawn capture to the left
		if nr, nc := row+1, col-1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].HasEnemyPiece(c.turn) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// black pawn capture to the right
		if nr, nc := row+1, col+1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].HasEnemyPiece(c.turn) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// black en-passant to the left
		if nr, nc := row+1, col-1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].IsEmpty() && c.enPassantPossible() && c.board.squares[nr][nc].CoordEquals(*c.enpSquare) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}

		// black en-passant to the right
		if nr, nc := row+1, col+1; validRowCol(uint8(nr), uint8(nc)) && c.board.squares[nr][nc].IsEmpty() && c.enPassantPossible() && c.board.squares[nr][nc].CoordEquals(*c.enpSquare) {
			toSquare := NewSquare(nr, nc, nil)
			move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
			moves = append(moves, move)
		}
	}

	return moves
}

func (c *Chess) generateRookPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	row, col := fromSquare.RowCol()

	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for _, dir := range dirs {
		for i := 1; i < boardSize; i++ {
			toRow := int(row) + int(dir[0])*i
			toCol := int(col) + int(dir[1])*i

			if validRowCol(toRow, toCol) {
				toSquare := c.board.squares[toRow][toCol]

				if toSquare.IsEmpty() {
					move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
					moves = append(moves, move)
				} else if toSquare.HasEnemyPiece(c.turn) {
					move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
					moves = append(moves, move)
					break
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return moves
}

func (c *Chess) generateKnightPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	row, col := fromSquare.RowCol()

	dirs := [][]int{{-1, 2}, {-1, -2}, {1, 2}, {1, -2}, {-2, 1}, {-2, -1}, {2, 1}, {2, -1}}

	for _, dir := range dirs {
		toRow := int(row) + int(dir[0])
		toCol := int(col) + int(dir[1])

		if validRowCol(toRow, toCol) {
			toSquare := c.board.squares[toRow][toCol]

			if toSquare.isEmptyOrHasEnemyPiece(c.turn) {
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (c *Chess) generateBishopPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	row, col := fromSquare.RowCol()

	dirs := [][]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}}

	for _, dir := range dirs {
		for i := 1; i < boardSize; i++ {
			toRow := int(row) + int(dir[0])*i
			toCol := int(col) + int(dir[1])*i

			if validRowCol(toRow, toCol) {
				toSquare := c.board.squares[toRow][toCol]

				if toSquare.IsEmpty() {
					move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
					moves = append(moves, move)
				} else if toSquare.HasEnemyPiece(c.turn) {
					move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
					moves = append(moves, move)
					break
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return moves
}

func (c *Chess) generateQueenPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	orthogonal := c.generateRookPseudoLegalMoves(fromSquare)
	diagonal := c.generateBishopPseudoLegalMoves(fromSquare)

	moves = append(moves, orthogonal...)
	moves = append(moves, diagonal...)

	return moves
}

func (c *Chess) generateKingPseudoLegalMoves(fromSquare Square) []*Move {
	moves := make([]*Move, 0)

	row, col := fromSquare.RowCol()

	dirs := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {1, 0}}

	for _, dir := range dirs {
		toRow := int(row) + int(dir[0])
		toCol := int(col) + int(dir[1])

		if validRowCol(toRow, toCol) {
			toSquare := c.board.squares[toRow][toCol]

			if toSquare.isEmptyOrHasEnemyPiece(c.turn) {
				move := NewMove(c.turn, fromSquare, *toSquare, *fromSquare.piece, toSquare.piece, c.enpSquare, 0, nil)
				moves = append(moves, move)
			}
		}
	}

	castlingMoves := c.generateCastlingMoves(fromSquare)
	moves = append(moves, castlingMoves...)

	return moves
}

func (c *Chess) generateAllPseudoLegalMovesForPiece(p Piece, fromSquare Square) []*Move {
	pseudoMoves := make([]*Move, 0)

	switch p.symbol {
	case Pawn:
		moves := c.generatePawnPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	case Rook:
		moves := c.generateRookPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	case Bishop:
		moves := c.generateBishopPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	case Knight:
		moves := c.generateKnightPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	case Queen:
		moves := c.generateQueenPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	case King:
		moves := c.generateKingPseudoLegalMoves(fromSquare)
		pseudoMoves = append(pseudoMoves, moves...)
	}

	return pseudoMoves
}

func (c *Chess) generateAllPseudoLegalMoves() []*Move {
	pseudoMoves := make([]*Move, 0)

	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			sq := c.board.squares[row][col]

			if sq.piece != nil {
				if (c.whiteToMove() && sq.piece.IsWhite()) || (c.blackToMove() && sq.piece.IsBlack()) {
					moves := c.generateAllPseudoLegalMovesForPiece(*sq.piece, *sq)
					pseudoMoves = append(pseudoMoves, moves...)
				}
			}
		}
	}

	return pseudoMoves
}

func (c *Chess) generateAllStrictlyLegalMoves() {
	pseudo := c.generateAllPseudoLegalMoves()

	legal := make([]*Move, 0)

	for _, pm := range pseudo {
		isCastle, castleSide := pm.IsCastle(c.castlingRights)

		if isCastle {
			if castleSide == kingSideCastle {
				tosq := c.board.squares[pm.FromColumn()][pm.ToColumn()+1]
				m := NewMove(c.turn, pm.fromSquare, *tosq, *pm.fromSquare.piece, nil, c.enpSquare, c.castlingRights, nil)

				inCheckWhileCastling := c.kingInCheckAfterPlaying(*m)
				if inCheckAfterCastling := c.kingInCheckAfterPlaying(*pm); !inCheckWhileCastling && !inCheckAfterCastling {
					legal = append(legal, pm)
				}
			}

			if castleSide == queenSideCastle {
				tosq := c.board.squares[pm.FromColumn()][pm.ToColumn()-1]
				m := NewMove(c.turn, pm.fromSquare, *tosq, *pm.fromSquare.piece, nil, c.enpSquare, c.castlingRights, nil)

				inCheckWhileCastling := c.kingInCheckAfterPlaying(*m)
				if inCheckAfterCastling := c.kingInCheckAfterPlaying(*pm); !inCheckWhileCastling && !inCheckAfterCastling {
					legal = append(legal, pm)
				}
			}
		} else {
			if inCheck := c.kingInCheckAfterPlaying(*pm); !inCheck {
				legal = append(legal, pm)
			}
		}
	}

	c.legalMoves = legal
}

func (c *Chess) updateGameTerminationState() {
	if c.kingInCheck() {
		c.check = true
	} else {
		c.check = false
	}

	if c.check && len(c.legalMoves) == 0 {
		c.checkmate = true
	} else {
		c.checkmate = false
	}

	if !c.check && len(c.legalMoves) == 0 {
		c.stalemate = true
	} else {
		c.stalemate = false
	}

	if c.isInsufficientMaterial() {
		c.insufficientMaterial = true
	} else {
		c.insufficientMaterial = false
	}

	// if c.isThreeFoldRepetition() {
	// 	c.threeFold = true
	// } else {
	// 	c.threeFold = false
	// }

	if c.halfMoves >= 100 || c.stalemate || c.insufficientMaterial || c.threeFold {
		c.draw = true
	} else {
		c.draw = false
	}

	if c.checkmate || c.draw {
		c.terminated = true
	} else {
		c.terminated = false
	}

	if c.terminated {
		if c.stalemate {
			c.outcome = "Draw by stalemate."
		} else if c.threeFold {
			c.outcome = "Draw by threefold repetition."
		} else if c.insufficientMaterial {
			c.outcome = "Draw by insufficient material."
		} else if c.checkmate {
			c.outcome = swapColor(c.turn).Name() + " wins by checkmate."
		}
	} else {
		c.outcome = "Match is in progress..."
	}
}

func (c *Chess) updateMovesClock(move Move) {
	if move.pieceMoved.IsPawn() {
		c.halfMoves = 0
	} else if move.pieceCaptured != nil {
		c.halfMoves = 0
	} else {
		c.halfMoves++
	}

	if move.IsBlack() {
		c.fullMoves++
	}
}

func (c *Chess) makeMoveInternal(move Move) {
	if move.fromSquare.piece.IsKing() {
		c.kingsPos[move.color] = move.toSquare.Coordinate()
	}

	c.board.squares[move.FromRow()][move.FromColumn()].piece = nil

	c.updateRookPositionAfterCastling(move)

	var prmSymb *PromotionPieceSymbol
	if move.IsPromotion() {
		prmPiece, prmPieceSymb := move.GetPromotedPiece()
		c.board.squares[move.ToRow()][move.ToColumn()].piece = prmPiece
		prmSymb = &prmPieceSymb
	} else {
		c.board.squares[move.ToRow()][move.ToColumn()].piece = &move.pieceMoved
	}

	if move.pieceCaptured != nil {
		c.capturedPieces = append(c.capturedPieces, *move.pieceCaptured)

		idx := slices.IndexFunc(c.alivePieces, func(p Piece) bool {
			return p.Equals(*move.pieceCaptured)
		})

		if idx != -1 {
			c.alivePieces = slices.Delete(c.alivePieces, idx, idx+1)
		}
	}

	history := NewHistory(move, c.castlingRights, c.enpSquare, prmSymb, 0, 1, c.kingsPos)
	c.history = append(c.history, history)

	// if move.pieceMoved.IsPawn() {
	// 	c.halfMoves = 0
	// } else if move.pieceCaptured != nil {
	// 	c.halfMoves = 0
	// } else {
	// 	c.halfMoves++
	// }

	// if move.IsBlack() {
	// 	c.fullMoves++
	// }

	c.updateEnPassantSquare(move)
	c.updateCastlingRights(move)
	c.switchTurn()
}

func (c *Chess) undoMoveInternal() {
	if len(c.history) == 0 {
		fmt.Printf("EMPTY HISTORY: unable to undo the move\n")
		return
	}

	old := c.history[len(c.history)-1]
	c.history = c.history[:len(c.history)-1]

	c.kingsPos = old.kingsPosition
	c.turn = old.move.color
	c.castlingRights = old.castlingRights
	// c.halfMoves = old.halfMoves
	// c.fullMoves = old.fullMoves

	c.board.squares[old.move.FromRow()][old.move.FromColumn()].piece = &old.move.pieceMoved
	c.board.squares[old.move.ToRow()][old.move.ToColumn()].piece = old.move.pieceCaptured

	if old.move.pieceMoved.IsKing() {
		c.kingsPos[old.move.color] = old.move.FromCoordinate()
	}

	if old.move.IsEnPassant(old.enpSquare) {
		c.board.squares[old.move.ToRow()][old.move.ToColumn()].piece = nil
		c.board.squares[old.move.FromRow()][old.move.ToColumn()].piece = old.move.pieceCaptured
	}

	isCastle, castleSide := old.move.IsCastle(old.castlingRights)

	if isCastle {
		if castleSide == kingSideCastle {
			tmpRook := c.board.squares[old.move.ToRow()][old.move.ToColumn()-1].piece
			c.board.squares[old.move.ToRow()][old.move.ToColumn()-1].piece = nil
			c.board.squares[old.move.ToRow()][old.move.ToColumn()+1].piece = tmpRook
		}

		if castleSide == queenSideCastle {
			tmpRook := c.board.squares[old.move.ToRow()][old.move.ToColumn()+1].piece
			c.board.squares[old.move.ToRow()][old.move.ToColumn()+1].piece = nil
			c.board.squares[old.move.ToRow()][old.move.ToColumn()-2].piece = tmpRook
		}
	}
}

func (c *Chess) MakeMove(move Move) {
	if ok := slices.ContainsFunc(c.legalMoves, func(m *Move) bool {
		return m.ID() == move.ID()
	}); !ok {
		fmt.Printf("ILLEGAL MOVE ATTEMPT: %v\n", move.String())
		return
	}

	c.makeMoveInternal(move)
	c.generateAllStrictlyLegalMoves()
	c.updateGameTerminationState()
	c.updateMovesClock(move)

	fmt.Printf("PLAYED MOVE: %v\n", move.String())
}

func (c *Chess) UndoMove() {
	c.undoMoveInternal()
	c.generateAllStrictlyLegalMoves()

	fmt.Printf("UNDONE MOVE:\n")
}

func (c *Chess) kingInCheckAfterPlaying(pseudoMove Move) bool {
	return c.kingAttackedByQueenOrBishopAfterPlaying(pseudoMove) ||
		c.kingAttackedByQueenOrRookAfterPlaying(pseudoMove) ||
		c.kingAttackedByKnightAfterPlaying(pseudoMove) ||
		c.kingAttackedByPawnAfterPlaying(pseudoMove) ||
		c.kingAttackedByKingAfterPlaying(pseudoMove)
}

func (c *Chess) kingInCheck() bool {
	return c.kingAttackedByQueenOrBishop() ||
		c.kingAttackedByQueenOrRook() ||
		c.kingAttackedByKnight() ||
		c.kingAttackedByPawn() ||
		c.kingAttackedByKing()
}

func (c *Chess) checkIfAttacked(pseudoMove Move, attackFn func() bool) bool {
	c.makeMoveInternal(pseudoMove)
	c.switchTurn()

	attacked := attackFn()

	c.switchTurn()
	c.undoMoveInternal()

	return attacked
}

func (c *Chess) kingAttackedByQueenOrRookAfterPlaying(pseudoMove Move) bool {
	return c.checkIfAttacked(pseudoMove, c.kingAttackedByQueenOrRook)
}

func (c *Chess) kingAttackedByQueenOrBishopAfterPlaying(pseudoMove Move) bool {
	return c.checkIfAttacked(pseudoMove, c.kingAttackedByQueenOrBishop)
}

func (c *Chess) kingAttackedByKnightAfterPlaying(pseudoMove Move) bool {
	return c.checkIfAttacked(pseudoMove, c.kingAttackedByKnight)
}

func (c *Chess) kingAttackedByPawnAfterPlaying(pseudoMove Move) bool {
	return c.checkIfAttacked(pseudoMove, c.kingAttackedByPawn)
}

func (c *Chess) kingAttackedByKingAfterPlaying(pseudoMove Move) bool {
	return c.checkIfAttacked(pseudoMove, c.kingAttackedByKing)
}

func (c *Chess) kingAttackedByQueenOrRook() bool {
	row, col := c.CurrentKingPos().RowCol()

	dirs := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	var attacked bool

loop:
	for _, dir := range dirs {
		for i := 1; i < boardSize; i++ {
			toRow := int(row) + int(dir[0])*i
			toCol := int(col) + int(dir[1])*i

			if validRowCol(toRow, toCol) {
				toSquare := c.board.squares[toRow][toCol]

				if toSquare.HasEnemyPiece(c.turn) && (toSquare.piece.IsRook() || toSquare.piece.IsQueen()) {
					attacked = true
					break loop
				}
			}
		}
	}

	return attacked
}

func (c *Chess) kingAttackedByQueenOrBishop() bool {
	row, col := c.CurrentKingPos().RowCol()

	dirs := [][]int{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}}

	var attacked bool

loop:
	for _, dir := range dirs {
		for i := 1; i < boardSize; i++ {
			toRow := int(row) + int(dir[0])*i
			toCol := int(col) + int(dir[1])*i

			if validRowCol(toRow, toCol) {
				toSquare := c.board.squares[toRow][toCol]

				if toSquare.HasEnemyPiece(c.turn) && (toSquare.piece.IsBishop() || toSquare.piece.IsQueen()) {
					attacked = true
					break loop
				}
			}
		}
	}

	return attacked
}

func (c *Chess) kingAttackedByKnight() bool {
	row, col := c.CurrentKingPos().RowCol()

	dirs := [][]int{{-1, 2}, {-1, -2}, {1, 2}, {1, -2}, {-2, 1}, {-2, -1}, {2, 1}, {2, -1}}

	var attacked bool

	for _, dir := range dirs {
		toRow := int(row) + int(dir[0])
		toCol := int(col) + int(dir[1])

		if validRowCol(toRow, toCol) {
			toSquare := c.board.squares[toRow][toCol]

			if toSquare.HasEnemyPiece(c.turn) && toSquare.piece.IsKnight() {
				attacked = true
				break
			}
		}
	}

	return attacked
}

func (c *Chess) kingAttackedByPawn() bool {
	row, col := c.CurrentKingPos().RowCol()

	var attacked bool

	checkForAttack := func(newRow, newCol uint8) {
		toSquare := c.board.squares[newRow][newCol]
		if toSquare.HasEnemyPiece(c.turn) && toSquare.piece.IsPawn() {
			attacked = true
		}
	}

	if c.whiteToMove() {
		if nr, nc := row-1, col-1; validRowCol(uint8(nr), uint8(nc)) {
			checkForAttack(uint8(nr), uint8(nc))
		}

		if nr, nc := row-1, col+1; validRowCol(uint8(nr), uint8(nc)) {
			checkForAttack(uint8(nr), uint8(nc))
		}
	}

	if c.blackToMove() {
		if nr, nc := row+1, col-1; validRowCol(uint8(nr), uint8(nc)) {
			checkForAttack(uint8(nr), uint8(nc))
		}

		if nr, nc := row+1, col+1; validRowCol(uint8(nr), uint8(nc)) {
			checkForAttack(uint8(nr), uint8(nc))
		}
	}

	return attacked
}

func (c *Chess) kingAttackedByKing() bool {
	row, col := c.CurrentKingPos().RowCol()

	dirs := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 1}, {1, 0}}

	var attacked bool

	for _, dir := range dirs {
		toRow := int(row) + int(dir[0])
		toCol := int(col) + int(dir[1])

		if validRowCol(toRow, toCol) {
			toSquare := c.board.squares[toRow][toCol]

			if toSquare.HasEnemyPiece(c.turn) && toSquare.piece.IsKing() {
				attacked = true
				break
			}
		}
	}

	return attacked
}
