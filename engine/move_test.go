package engine

import (
	"fmt"
	"testing"
)

func TestMoves(t *testing.T) {
	testCases := map[string]struct {
		m                                    Move
		label                                string
		src, dest                            Square
		piece                                Piece
		promo                                Promotion
		isCapture, isDouble, isEnp, isCastle bool
	}{
		"P quiet":         {m: newQuietMove(A2, A3, WhitePawn), label: "a2a3", src: A2, dest: A3, piece: WhitePawn, promo: PromotionNone, isCapture: false, isDouble: false, isEnp: false, isCastle: false},
		"p double":        {m: newDoublePawnMove(E7, E5, BlackPawn), label: "e7e5", src: E7, dest: E5, piece: BlackPawn, promo: PromotionNone, isCapture: false, isDouble: true, isEnp: false, isCastle: false},
		"P enp":           {m: newEnpCaptureMove(E2, D3, WhitePawn), label: "e2d3", src: E2, dest: D3, piece: WhitePawn, promo: PromotionNone, isCapture: true, isDouble: false, isEnp: true, isCastle: false},
		"Q cap":           {m: newCaptureMove(H4, H8, WhiteQueen), label: "h4h8", src: H4, dest: H8, piece: WhiteQueen, promo: PromotionNone, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
		"promo Q":         {m: newPromotionMove(G7, G8, WhitePawn, PromotionQueen), label: "g7g8q", src: G7, dest: G8, piece: WhitePawn, promo: PromotionQueen, isCapture: false, isDouble: false, isEnp: false, isCastle: false},
		"promo R":         {m: newPromotionMove(G7, G8, WhitePawn, PromotionRook), label: "g7g8r", src: G7, dest: G8, piece: WhitePawn, promo: PromotionRook, isCapture: false, isDouble: false, isEnp: false, isCastle: false},
		"promo B":         {m: newPromotionMove(G7, G8, WhitePawn, PromotionBishop), label: "g7g8b", src: G7, dest: G8, piece: WhitePawn, promo: PromotionBishop, isCapture: false, isDouble: false, isEnp: false, isCastle: false},
		"promo N":         {m: newPromotionMove(G7, G8, WhitePawn, PromotionKnight), label: "g7g8n", src: G7, dest: G8, piece: WhitePawn, promo: PromotionKnight, isCapture: false, isDouble: false, isEnp: false, isCastle: false},
		"capture promo Q": {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionQueen), label: "g7h8q", src: G7, dest: H8, piece: WhitePawn, promo: PromotionQueen, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
		"capture promo R": {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionRook), label: "g7h8r", src: G7, dest: H8, piece: WhitePawn, promo: PromotionRook, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
		"capture promo B": {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionBishop), label: "g7h8b", src: G7, dest: H8, piece: WhitePawn, promo: PromotionBishop, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
		"capture promo N": {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionKnight), label: "g7h8n", src: G7, dest: H8, piece: WhitePawn, promo: PromotionKnight, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.m.String() != tc.label {
				t.Fatalf("invalid move label: want %s, got %s", tc.label, tc.m.String())
			}
			if tc.m.Src() != tc.src {
				t.Fatalf("invalid move src square: want %s, got %s", tc.src, tc.m.Src())
			}
			if tc.m.Dest() != tc.dest {
				t.Fatalf("invalid move dest square: want %s, got %s", tc.dest, tc.m.Dest())
			}
			if tc.m.Piece() != tc.piece {
				t.Fatalf("invalid move piece: want %s, got %s", tc.piece, tc.m.Piece())
			}
			if tc.m.Promotion() != tc.promo {
				t.Fatalf("invalid move promotion: want %v, got %v", tc.promo, tc.m.Promotion())
			}
			if (tc.m.IsCapture() != tc.isCapture) || (tc.m.IsDoublePawn() != tc.isDouble) || (tc.m.IsEnPassant() != tc.isEnp) || (tc.m.IsCastle() != tc.isCastle) {
				t.Fatalf("invalid move flags (cap, dbl, enp, castle): want (%v %v %v %v), got (%v %v %v %v)", tc.isCapture, tc.isDouble, tc.isEnp, tc.isCastle, tc.m.IsCapture(), tc.m.IsDoublePawn(), tc.m.IsEnPassant(), tc.m.IsCastle())
			}
		})
	}

	src, dest, piece := B7, B8, WhitePawn
	promoPieces := []Promotion{PromotionQueen, PromotionRook, PromotionBishop, PromotionKnight}

	promos := newPossiblePromotionMoves(src, dest, piece)
	promoLabels := []string{"b7b8q", "b7b8r", "b7b8b", "b7b8n"}

	t.Run("possible promotion moves", func(t *testing.T) {
		for i, pm := range promos {
			t.Run(fmt.Sprintf("promo: %v", promos[i]), func(t *testing.T) {
				if pm.String() != promoLabels[i] {
					t.Fatalf("invalid move label: want %s, got %s", promoLabels[i], pm.String())
				}
				if pm.Src() != src {
					t.Fatalf("invalid move src square: want %s, got %s", src, pm.Src())
				}
				if pm.Dest() != dest {
					t.Fatalf("invalid move dest square: want %s, got %s", dest, pm.Dest())
				}
				if pm.Piece() != piece {
					t.Fatalf("invalid move piece: want %s, got %s", piece, pm.Piece())
				}
				if pm.Promotion() != promoPieces[i] {
					t.Fatalf("invalid move promotion: want %v, got %v", promoPieces[i], pm.Promotion())
				}
				if pm.IsCapture() || pm.IsDoublePawn() || pm.IsEnPassant() || pm.IsCastle() {
					t.Fatalf("invalid move flags (cap, dbl, enp, castle): want (%v %v %v %v), got (%v %v %v %v)", false, false, false, false, pm.IsCapture(), pm.IsDoublePawn(), pm.IsEnPassant(), pm.IsCastle())
				}
			})
		}
	})

	promoCaptures := newPossiblePromotionCaptureMoves(src, dest, piece)
	promoCapturesLabels := []string{"b7b8q", "b7b8r", "b7b8b", "b7b8n"}

	t.Run("possible promotion captures moves", func(t *testing.T) {
		for i, pm := range promoCaptures {
			t.Run(fmt.Sprintf("promo: %v", i+1), func(t *testing.T) {
				if pm.String() != promoCapturesLabels[i] {
					t.Fatalf("invalid move label: want %s, got %s", promoCapturesLabels[i], pm.String())
				}
				if pm.Src() != src {
					t.Fatalf("invalid move src square: want %s, got %s", src, pm.Src())
				}
				if pm.Dest() != dest {
					t.Fatalf("invalid move dest square: want %s, got %s", dest, pm.Dest())
				}
				if pm.Piece() != piece {
					t.Fatalf("invalid move piece: want %s, got %s", piece, pm.Piece())
				}
				if pm.Promotion() != promoPieces[i] {
					t.Fatalf("invalid move promotion: want %v, got %v", promoPieces[i], pm.Promotion())
				}
				if !pm.IsCapture() || pm.IsDoublePawn() || pm.IsEnPassant() || pm.IsCastle() {
					t.Fatalf("invalid move flags (cap, dbl, enp, castle): want (%v %v %v %v), got (%v %v %v %v)", true, false, false, false, pm.IsCapture(), pm.IsDoublePawn(), pm.IsEnPassant(), pm.IsCastle())
				}
			})
		}
	})
}
