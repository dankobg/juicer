package engine

import "fmt"

type Board struct {
	squares [][]*Square
}

func NewBoard(fen string) (*Board, error) {
	b := &Board{}

	b.initializeSquares()

	if _, err := b.LoadFromFEN(fen); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Board) String() string {
	return b.Draw()
}

func (b *Board) Squares() [][]*Square {
	return b.squares
}

func (b *Board) FEN(meta Metadata) string {
	return convertMetadataToFEN(meta)
}

func (b *Board) Draw() string {
	return drawBoardASCII(b.squares)
}

func (b *Board) initializeSquares() {
	b.squares = initBoardSquares()
}

func (b *Board) LoadFromFEN(FEN string) (*Metadata, error) {
	fenStr := FENStartingPosition
	if FEN != "" {
		fenStr = FEN
	}

	meta, err := ParseMetadataFromFEN(fenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load board from FEN: %w", err)
	}

	for r := 0; r < len(meta.board); r++ {
		for c := 0; c < len(meta.board[r]); c++ {
			if meta.board[r][c].HasPiece() {
				b.squares[r][c].piece = meta.board[r][c].piece
			}
		}
	}

	return meta, nil
}
