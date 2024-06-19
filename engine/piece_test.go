package engine

import "testing"

func TestPieceKinds(t *testing.T) {
	testCases := map[string]struct {
		pk     PieceKind
		symbol string
	}{
		"King is k":   {pk: King, symbol: "k"},
		"Queen is q":  {pk: Queen, symbol: "q"},
		"Rook is r":   {pk: Rook, symbol: "r"},
		"Bishop is b": {pk: Bishop, symbol: "b"},
		"Knight is n": {pk: Knight, symbol: "n"},
		"Pawn is p":   {pk: Pawn, symbol: "p"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.pk.String() != tc.symbol {
				t.Fatalf("invalid piece kind: want %s, got %s", tc.symbol, tc.pk)
			}
		})
	}
}

func TestPieces(t *testing.T) {
	testCases := map[string]struct {
		piece Piece
		fen   string
		kind  PieceKind
		color Color
	}{
		"WhiteKing is K":   {piece: WhiteKing, fen: "K", kind: King, color: White},
		"WhiteQueen is Q":  {piece: WhiteQueen, fen: "Q", kind: Queen, color: White},
		"WhiteRook is R":   {piece: WhiteRook, fen: "R", kind: Rook, color: White},
		"WhiteBishop is B": {piece: WhiteBishop, fen: "B", kind: Bishop, color: White},
		"WhiteKnight is N": {piece: WhiteKnight, fen: "N", kind: Knight, color: White},
		"WhitePawn is P":   {piece: WhitePawn, fen: "P", kind: Pawn, color: White},
		"BlackKing is k":   {piece: BlackKing, fen: "k", kind: King, color: Black},
		"BlackQueen is q":  {piece: BlackQueen, fen: "q", kind: Queen, color: Black},
		"BlackRook is r":   {piece: BlackRook, fen: "r", kind: Rook, color: Black},
		"BlackBishop is b": {piece: BlackBishop, fen: "b", kind: Bishop, color: Black},
		"BlackKnight is n": {piece: BlackKnight, fen: "n", kind: Knight, color: Black},
		"BlackPawn is p":   {piece: BlackPawn, fen: "p", kind: Pawn, color: Black},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.piece.FENSymbol() != tc.fen {
				t.Fatalf("invalid piece fen symbol: want %s, got %s", tc.fen, tc.piece.FENSymbol())
			}
			if tc.piece.Kind() != tc.kind {
				t.Fatalf("invalid piece kind: want %s, got %s", tc.kind, tc.piece.Kind())
			}
			if tc.piece.Color() != tc.color {
				t.Fatalf("invalid piece color: want %s, got %s", tc.color, tc.piece.Color())
			}
		})
	}
}

func TestNewPieceFromFenSymbol(t *testing.T) {
	testCases := map[string]struct {
		in      string
		want    Piece
		wantErr bool
	}{
		"creates white king":       {in: "K", wantErr: false, want: WhiteKing},
		"creates white queen":      {in: "Q", wantErr: false, want: WhiteQueen},
		"creates white rook":       {in: "R", wantErr: false, want: WhiteRook},
		"creates white bishop":     {in: "B", wantErr: false, want: WhiteBishop},
		"creates white knight":     {in: "N", wantErr: false, want: WhiteKnight},
		"creates white pawn":       {in: "P", wantErr: false, want: WhitePawn},
		"creates black king":       {in: "k", wantErr: false, want: BlackKing},
		"creates black queen":      {in: "q", wantErr: false, want: BlackQueen},
		"creates black rook":       {in: "r", wantErr: false, want: BlackRook},
		"creates black bishop":     {in: "b", wantErr: false, want: BlackBishop},
		"creates black knight":     {in: "n", wantErr: false, want: BlackKnight},
		"creates black pawn":       {in: "p", wantErr: false, want: BlackPawn},
		"fails with invalid empty": {in: "", wantErr: true},
		"fails with invalid fen":   {in: "j", wantErr: true},
		"fails with invalid white": {in: "wr", wantErr: true},
		"fails with invalid black": {in: "br", wantErr: true},
		"fails with invalid len w": {in: "PP", wantErr: true},
		"fails with invalid len b": {in: "pp", wantErr: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			piece, err := NewPieceFromFenSymbol(tc.in)

			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid piece, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if !tc.wantErr && piece != tc.want {
				t.Fatalf("invalid piece: want %s, got %s", tc.want, piece)
			}
		})
	}
}
