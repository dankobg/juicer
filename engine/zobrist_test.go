package engine

import (
	"testing"
)

func TestZobristSeed(t *testing.T) {
	InitPrecalculatedTables()

	p1, p2 := &Position{}, &Position{}

	if err := p1.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos1: %v", err)
	}
	if err := p2.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos2: %v", err)
	}

	if p1.Hash != p2.Hash {
		t.Fatalf("zobrist seeds are not equal: (%v, %v)", p1.Hash, p2.Hash)
	}
}

func TestZobristTransposition(t *testing.T) {
	InitPrecalculatedTables()

	p1, p2 := &Position{}, &Position{}

	if err := p1.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos1: %v", err)
	}
	if err := p2.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos2: %v", err)
	}

	p1.MakeMove(newDoublePawnMove(E2, E4, WhitePawn))
	p1.MakeMove(newDoublePawnMove(E7, E5, BlackPawn))
	p1.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	p1.MakeMove(newQuietMove(B8, C6, BlackKnight))

	p2.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	p2.MakeMove(newDoublePawnMove(E7, E5, BlackPawn))
	p2.MakeMove(newDoublePawnMove(E2, E4, WhitePawn))
	p2.MakeMove(newQuietMove(B8, C6, BlackKnight))

	if p1.Hash != p2.Hash {
		t.Fatalf("zobrist transposition hashes are not equal: (%v, %v)", p1.Hash, p2.Hash)
	}
}

func TestZobristEnpTranspositionDiff(t *testing.T) {
	InitPrecalculatedTables()

	p1, p2 := &Position{}, &Position{}

	if err := p1.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos1: %v", err)
	}
	if err := p2.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos2: %v", err)
	}

	p1.MakeMove(newDoublePawnMove(E2, E4, WhitePawn))
	p1.MakeMove(newDoublePawnMove(E7, E5, BlackPawn))
	p1.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	p1.MakeMove(newQuietMove(B8, C6, BlackKnight))

	p2.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	p2.MakeMove(newQuietMove(B8, C6, BlackKnight))
	p2.MakeMove(newDoublePawnMove(E2, E4, WhitePawn))
	p2.MakeMove(newDoublePawnMove(E7, E5, BlackPawn))

	if p1.Hash == p2.Hash {
		t.Fatalf("zobrist transposition hashes are equal: (%v, %v)", p1.Hash, p2.Hash)
	}
}

func TestZobristTurnDiff(t *testing.T) {
	InitPrecalculatedTables()

	p1, p2 := &Position{}, &Position{}

	if err := p1.LoadFromFEN("rnbqkbnr/pppp1ppp/4p3/8/8/4PP2/PPPP2PP/RNBQKBNR b - - 0 2"); err != nil {
		t.Fatalf("failed to load pos1: %v", err)
	}
	if err := p2.LoadFromFEN(FENStartingPosition); err != nil {
		t.Fatalf("failed to load pos2: %v", err)
	}

	p2.MakeMove(newQuietMove(E8, E7, BlackKing))
	p2.MakeMove(newQuietMove(E1, F2, WhiteKing))
	p2.MakeMove(newQuietMove(E7, E8, BlackKing))
	p2.MakeMove(newQuietMove(F2, E2, WhiteKing))
	p2.MakeMove(newQuietMove(E8, E7, BlackKing))
	p2.MakeMove(newQuietMove(E2, E1, WhiteKing))
	p2.MakeMove(newQuietMove(E7, E8, BlackKing))

	if p1.Hash == p2.Hash {
		t.Fatalf("zobrist transposition hashes are equal: (%v, %v)", p1.Hash, p2.Hash)
	}
}
