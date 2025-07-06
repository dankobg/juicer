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
		"P enpassant":     {m: newEnpCaptureMove(E2, D3, WhitePawn), label: "e2d3", src: E2, dest: D3, piece: WhitePawn, promo: PromotionNone, isCapture: true, isDouble: false, isEnp: true, isCastle: false},
		"Q capture":       {m: newCaptureMove(H4, H8, WhiteQueen), label: "h4h8", src: H4, dest: H8, piece: WhiteQueen, promo: PromotionNone, isCapture: true, isDouble: false, isEnp: false, isCastle: false},
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

func TestMoveNotations(t *testing.T) {
	testCases := map[string]struct {
		m           Move
		isCheck     bool
		isCheckmate bool
		uci         string
		san         string
		lan         string
		legalMoves  []Move
	}{
		"P quiet w":                             {m: newQuietMove(A2, A3, WhitePawn), uci: "a2a3", san: "a3", lan: "a2-a3"},
		"P quiet b":                             {m: newQuietMove(G7, G6, BlackPawn), uci: "g7g6", san: "g6", lan: "g7-g6"},
		"p double w":                            {m: newDoublePawnMove(B2, B4, WhitePawn), uci: "b2b4", san: "b4", lan: "b2-b4"},
		"p double b":                            {m: newDoublePawnMove(E7, E5, BlackPawn), uci: "e7e5", san: "e5", lan: "e7-e5"},
		"P enp w":                               {m: newEnpCaptureMove(E5, D6, WhitePawn), uci: "e5d6", san: "exd6", lan: "e5xd6"},
		"P enp b":                               {m: newEnpCaptureMove(D4, E3, BlackPawn), uci: "d4e3", san: "dxe3", lan: "d4xe3"},
		"enp checkmate":                         {m: newEnpCaptureMove(E5, F6, WhitePawn), isCheck: true, isCheckmate: true, uci: "e5f6", san: "exf6#", lan: "e5xf6#"},
		"Q move":                                {m: newQuietMove(D1, H4, WhiteQueen), uci: "d1h4", san: "Qh4", lan: "Qd1-h4"},
		"Q capture":                             {m: newCaptureMove(D1, D8, WhiteQueen), uci: "d1d8", san: "Qxd8", lan: "Qd1xd8"},
		"W kingside castle":                     {m: newCastleMove(E1, G1, WhiteKing), uci: "e1g1", san: "O-O", lan: "Ke1-g1"},
		"W queenside castle":                    {m: newCastleMove(E1, C1, WhiteKing), uci: "e1c1", san: "O-O-O", lan: "Ke1-c1"},
		"B kingside castle":                     {m: newCastleMove(E8, G8, BlackKing), uci: "e8g8", san: "O-O", lan: "Ke8-g8"},
		"B queenside castle":                    {m: newCastleMove(E8, C8, BlackKing), uci: "e8c8", san: "O-O-O", lan: "Ke8-c8"},
		"W castle check":                        {m: newCastleMove(E1, G1, WhiteKing), isCheck: true, uci: "e1g1", san: "O-O+", lan: "Ke1-g1+"},
		"W castle checkmate":                    {m: newCastleMove(E1, G1, WhiteKing), isCheck: true, isCheckmate: true, uci: "e1g1", san: "O-O#", lan: "Ke1-g1#"},
		"B castle check":                        {m: newCastleMove(E8, G8, BlackKing), isCheck: true, uci: "e8g8", san: "O-O+", lan: "Ke8-g8+"},
		"B castle checkmate":                    {m: newCastleMove(E8, G8, BlackKing), isCheck: true, isCheckmate: true, uci: "e8g8", san: "O-O#", lan: "Ke8-g8#"},
		"promo Q":                               {m: newPromotionMove(G7, G8, WhitePawn, PromotionQueen), uci: "g7g8q", san: "g8=Q", lan: "g7-g8=Q"},
		"promo Q capture":                       {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionQueen), uci: "h7g8q", san: "hxg8=Q", lan: "h7xg8=Q"},
		"promo Q check":                         {m: newPromotionMove(G7, G8, WhitePawn, PromotionQueen), isCheck: true, uci: "g7g8q", san: "g8=Q+", lan: "g7-g8=Q+"},
		"promo Q capture check":                 {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionQueen), isCheck: true, uci: "g7h8q", san: "gxh8=Q+", lan: "g7xh8=Q+"},
		"promo Q checkmate":                     {m: newPromotionMove(G7, G8, WhitePawn, PromotionQueen), isCheck: true, isCheckmate: true, uci: "g7g8q", san: "g8=Q#", lan: "g7-g8=Q#"},
		"promo Q capture checkmate":             {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionQueen), isCheck: true, isCheckmate: true, uci: "h7g8q", san: "hxg8=Q#", lan: "h7xg8=Q#"},
		"promo R":                               {m: newPromotionMove(G7, G8, WhitePawn, PromotionRook), uci: "g7g8r", san: "g8=R", lan: "g7-g8=R"},
		"promo R capture":                       {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionRook), uci: "h7g8r", san: "hxg8=R", lan: "h7xg8=R"},
		"promo R check":                         {m: newPromotionMove(G7, G8, WhitePawn, PromotionRook), isCheck: true, uci: "g7g8r", san: "g8=R+", lan: "g7-g8=R+"},
		"promo R capture check":                 {m: newPromotionCaptureMove(H7, H8, WhitePawn, PromotionRook), isCheck: true, uci: "h7h8r", san: "hxh8=R+", lan: "h7xh8=R+"},
		"promo R checkmate":                     {m: newPromotionMove(G7, G8, WhitePawn, PromotionRook), isCheck: true, isCheckmate: true, uci: "g7g8r", san: "g8=R#", lan: "g7-g8=R#"},
		"promo R capture checkmate":             {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionRook), isCheck: true, isCheckmate: true, uci: "h7g8r", san: "hxg8=R#", lan: "h7xg8=R#"},
		"promo B":                               {m: newPromotionMove(G7, G8, WhitePawn, PromotionBishop), uci: "g7g8b", san: "g8=B", lan: "g7-g8=B"},
		"promo B capture":                       {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionBishop), uci: "h7g8b", san: "hxg8=B", lan: "h7xg8=B"},
		"promo B check":                         {m: newPromotionMove(G7, G8, WhitePawn, PromotionBishop), isCheck: true, uci: "g7g8b", san: "g8=B+", lan: "g7-g8=B+"},
		"promo B capture check":                 {m: newPromotionCaptureMove(H7, H8, WhitePawn, PromotionBishop), isCheck: true, uci: "h7h8b", san: "hxh8=B+", lan: "h7xh8=B+"},
		"promo B checkmate":                     {m: newPromotionMove(H7, H8, WhitePawn, PromotionBishop), isCheck: true, isCheckmate: true, uci: "h7h8b", san: "h8=B#", lan: "h7-h8=B#"},
		"promo B capture checkmate":             {m: newPromotionCaptureMove(G7, H8, WhitePawn, PromotionBishop), isCheck: true, isCheckmate: true, uci: "g7h8b", san: "gxh8=B#", lan: "g7xh8=B#"},
		"promo N":                               {m: newPromotionMove(G7, G8, WhitePawn, PromotionKnight), uci: "g7g8n", san: "g8=N", lan: "g7-g8=N"},
		"promo N capture":                       {m: newPromotionCaptureMove(H7, G8, WhitePawn, PromotionKnight), uci: "h7g8n", san: "hxg8=N", lan: "h7xg8=N"},
		"promo N check":                         {m: newPromotionMove(G7, G8, WhitePawn, PromotionKnight), isCheck: true, uci: "g7g8n", san: "g8=N+", lan: "g7-g8=N+"},
		"promo N capture check":                 {m: newPromotionCaptureMove(H7, H8, WhitePawn, PromotionKnight), isCheck: true, uci: "h7h8n", san: "hxh8=N+", lan: "h7xh8=N+"},
		"promo N checkmate":                     {m: newPromotionMove(C7, C8, WhitePawn, PromotionKnight), isCheck: true, isCheckmate: true, uci: "c7c8n", san: "c8=N#", lan: "c7-c8=N#"},
		"promo N capture checkmate":             {m: newPromotionCaptureMove(D7, C8, WhitePawn, PromotionKnight), isCheck: true, isCheckmate: true, uci: "d7c8n", san: "dxc8=N#", lan: "d7xc8=N#"},
		"2 knights ambiguity":                   {m: newQuietMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngf3", lan: "Ng1-f3", legalMoves: []Move{newQuietMove(G1, F3, WhiteKnight), newQuietMove(D2, F3, WhiteKnight)}},
		"2 knights ambiguity capture":           {m: newCaptureMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngxf3", lan: "Ng1xf3", legalMoves: []Move{newCaptureMove(G1, F3, WhiteKnight), newCaptureMove(D2, F3, WhiteKnight)}},
		"2 knights ambiguity same file":         {m: newQuietMove(G1, F3, WhiteKnight), uci: "g1f3", san: "N1f3", lan: "Ng1-f3", legalMoves: []Move{newQuietMove(G1, F3, WhiteKnight), newQuietMove(G5, F3, WhiteKnight)}},
		"2 knights ambiguity same file capture": {m: newCaptureMove(G1, F3, WhiteKnight), uci: "g1f3", san: "N1xf3", lan: "Ng1xf3", legalMoves: []Move{newCaptureMove(G1, F3, WhiteKnight), newCaptureMove(G5, F3, WhiteKnight)}},
		"2 knights ambiguity same rank":         {m: newQuietMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngf3", lan: "Ng1-f3", legalMoves: []Move{newQuietMove(E1, F3, WhiteKnight), newQuietMove(E1, F3, WhiteKnight)}},
		"2 knights ambiguity same rank capture": {m: newCaptureMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngxf3", lan: "Ng1xf3", legalMoves: []Move{newCaptureMove(E1, F3, WhiteKnight), newCaptureMove(E1, F3, WhiteKnight)}},
		"3 knights ambiguity":                   {m: newQuietMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngf3", lan: "Ng1-f3", legalMoves: []Move{newQuietMove(E1, F3, WhiteKnight), newQuietMove(E1, F3, WhiteKnight), newQuietMove(D4, F3, WhiteKnight)}},
		"3 knights ambiguity capture":           {m: newCaptureMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ngxf3", lan: "Ng1xf3", legalMoves: []Move{newCaptureMove(E1, F3, WhiteKnight), newCaptureMove(E1, F3, WhiteKnight), newCaptureMove(D4, F3, WhiteKnight)}},
		"8 knights ambiguity":                   {m: newQuietMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ng1f3", lan: "Ng1-f3", legalMoves: []Move{newQuietMove(G1, F3, WhiteKnight), newQuietMove(E1, F3, WhiteKnight), newQuietMove(D2, F3, WhiteKnight), newQuietMove(D4, F3, WhiteKnight), newQuietMove(E5, F3, WhiteKnight), newQuietMove(G5, F3, WhiteKnight), newQuietMove(H4, F3, WhiteKnight), newQuietMove(H2, F3, WhiteKnight)}},
		"8 knights ambiguity capture":           {m: newCaptureMove(G1, F3, WhiteKnight), uci: "g1f3", san: "Ng1xf3", lan: "Ng1xf3", legalMoves: []Move{newCaptureMove(G1, F3, WhiteKnight), newCaptureMove(E1, F3, WhiteKnight), newCaptureMove(D2, F3, WhiteKnight), newCaptureMove(D4, F3, WhiteKnight), newCaptureMove(E5, F3, WhiteKnight), newCaptureMove(G5, F3, WhiteKnight), newCaptureMove(H4, F3, WhiteKnight), newCaptureMove(H2, F3, WhiteKnight)}},
		"2 rooks ambiguity":                     {m: newQuietMove(H1, H3, WhiteRook), uci: "h1h3", san: "Rhh3", lan: "Rh1-h3", legalMoves: []Move{newQuietMove(H1, H3, WhiteRook), newQuietMove(F3, H3, WhiteRook)}},
		"2 rooks ambiguity capture":             {m: newCaptureMove(H1, H3, WhiteRook), uci: "h1h3", san: "Rhxh3", lan: "Rh1xh3", legalMoves: []Move{newCaptureMove(H1, H3, WhiteRook), newCaptureMove(F3, H3, WhiteRook)}},
		"2 rooks ambiguity same file":           {m: newQuietMove(H1, H3, WhiteRook), uci: "h1h3", san: "R1h3", lan: "Rh1-h3", legalMoves: []Move{newQuietMove(H1, H3, WhiteRook), newQuietMove(H5, H3, WhiteRook)}},
		"2 rooks ambiguity same file capture":   {m: newCaptureMove(H1, H3, WhiteRook), uci: "h1h3", san: "R1xh3", lan: "Rh1xh3", legalMoves: []Move{newCaptureMove(H1, H3, WhiteRook), newCaptureMove(H5, H3, WhiteRook)}},
		"2 rooks ambiguity same rank":           {m: newQuietMove(H1, F1, WhiteRook), uci: "h1f1", san: "Rhf1", lan: "Rh1-f1", legalMoves: []Move{newQuietMove(H1, F1, WhiteRook), newQuietMove(D1, F1, WhiteRook)}},
		"2 rooks ambiguity same rank capture":   {m: newCaptureMove(H1, F1, WhiteRook), uci: "h1f1", san: "Rhxf1", lan: "Rh1xf1", legalMoves: []Move{newCaptureMove(H1, F1, WhiteRook), newCaptureMove(D1, F1, WhiteRook)}},
		"4 rooks ambiguity horiz":               {m: newQuietMove(D3, F3, WhiteRook), uci: "d3f3", san: "Rdf3", lan: "Rd3-f3", legalMoves: []Move{newQuietMove(F1, F3, WhiteRook), newQuietMove(D3, F3, WhiteRook), newQuietMove(H3, F3, WhiteRook), newQuietMove(F5, F3, WhiteRook)}},
		"4 rooks ambiguity horiz capture":       {m: newCaptureMove(D3, F3, WhiteRook), uci: "d3f3", san: "Rdxf3", lan: "Rd3xf3", legalMoves: []Move{newCaptureMove(F1, F3, WhiteRook), newCaptureMove(D3, F3, WhiteRook), newCaptureMove(H3, F3, WhiteRook), newCaptureMove(F5, F3, WhiteRook)}},
		"4 rooks ambiguity vert":                {m: newQuietMove(F1, F3, WhiteRook), uci: "f1f3", san: "R1f3", lan: "Rf1-f3", legalMoves: []Move{newQuietMove(F1, F3, WhiteRook), newQuietMove(D3, F3, WhiteRook), newQuietMove(H3, F3, WhiteRook), newQuietMove(F5, F3, WhiteRook)}},
		"4 rooks ambiguity vert capture":        {m: newCaptureMove(F1, F3, WhiteRook), uci: "f1f3", san: "R1xf3", lan: "Rf1xf3", legalMoves: []Move{newCaptureMove(F1, F3, WhiteRook), newCaptureMove(D3, F3, WhiteRook), newCaptureMove(H3, F3, WhiteRook), newCaptureMove(F5, F3, WhiteRook)}},
		"2 bishops ambiguity":                   {m: newQuietMove(C2, E4, WhiteBishop), uci: "c2e4", san: "Bce4", lan: "Bc2-e4", legalMoves: []Move{newQuietMove(C2, E4, WhiteBishop), newQuietMove(G2, E4, WhiteBishop)}},
		"2 bishops ambiguity capture":           {m: newCaptureMove(C2, E4, WhiteBishop), uci: "c2e4", san: "Bcxe4", lan: "Bc2xe4", legalMoves: []Move{newCaptureMove(C2, E4, WhiteBishop), newCaptureMove(G2, E4, WhiteBishop)}},
		"2 bishops same file":                   {m: newQuietMove(G2, E4, WhiteBishop), uci: "g2e4", san: "B2e4", lan: "Bg2-e4", legalMoves: []Move{newQuietMove(G6, E4, WhiteBishop), newQuietMove(G2, E4, WhiteBishop)}},
		"2 bishops same file capture":           {m: newCaptureMove(G2, E4, WhiteBishop), uci: "g2e4", san: "B2xe4", lan: "Bg2xe4", legalMoves: []Move{newCaptureMove(G6, E4, WhiteBishop), newCaptureMove(G2, E4, WhiteBishop)}},
		"2 bishops same rank":                   {m: newQuietMove(G2, E4, WhiteBishop), uci: "g2e4", san: "B2e4", lan: "Bg2-e4", legalMoves: []Move{newQuietMove(G6, E4, WhiteBishop), newQuietMove(G2, E4, WhiteBishop)}},
		"2 bishops same rank capture":           {m: newCaptureMove(G2, E4, WhiteBishop), uci: "g2e4", san: "B2xe4", lan: "Bg2xe4", legalMoves: []Move{newCaptureMove(G6, E4, WhiteBishop), newCaptureMove(G2, E4, WhiteBishop)}},
		"4 bishops":                             {m: newQuietMove(G2, E4, WhiteBishop), uci: "g2e4", san: "Bg2e4", lan: "Bg2-e4", legalMoves: []Move{newQuietMove(G6, E4, WhiteBishop), newQuietMove(G2, E4, WhiteBishop), newQuietMove(C6, E4, WhiteBishop), newQuietMove(C2, E4, WhiteBishop)}},
		"4 bishops capture":                     {m: newCaptureMove(G2, E4, WhiteBishop), uci: "g2e4", san: "Bg2xe4", lan: "Bg2xe4", legalMoves: []Move{newCaptureMove(G6, E4, WhiteBishop), newCaptureMove(G2, E4, WhiteBishop), newCaptureMove(C6, E4, WhiteBishop), newCaptureMove(C2, E4, WhiteBishop)}},
		"2 queens ambiguity":                    {m: newQuietMove(H1, H5, WhiteQueen), uci: "h1h5", san: "Qhh5", lan: "Qh1-h5", legalMoves: []Move{newQuietMove(H1, H5, WhiteQueen), newQuietMove(E2, H5, WhiteQueen)}},
		"2 queens ambiguity capture":            {m: newCaptureMove(H1, H5, WhiteQueen), uci: "h1h5", san: "Qhxh5", lan: "Qh1xh5", legalMoves: []Move{newCaptureMove(H1, H5, WhiteQueen), newCaptureMove(E2, H5, WhiteQueen)}},
		"2 queens ambiguity same file":          {m: newQuietMove(F1, G2, WhiteQueen), uci: "f1g2", san: "Qfg2", lan: "Qf1-g2", legalMoves: []Move{newQuietMove(F1, G2, WhiteQueen), newQuietMove(H1, G2, WhiteQueen)}},
		"2 queens ambiguity same file capture":  {m: newCaptureMove(F1, G2, WhiteQueen), uci: "f1g2", san: "Qfxg2", lan: "Qf1xg2", legalMoves: []Move{newCaptureMove(F1, G2, WhiteQueen), newCaptureMove(H1, G2, WhiteQueen)}},
		"2 queens ambiguity same rank":          {m: newQuietMove(H1, G2, WhiteQueen), uci: "h1g2", san: "Q1g2", lan: "Qh1-g2", legalMoves: []Move{newQuietMove(H1, G2, WhiteQueen), newQuietMove(H3, G2, WhiteQueen)}},
		"2 queens ambiguity same rank capture":  {m: newCaptureMove(H1, G2, WhiteQueen), uci: "h1g2", san: "Q1xg2", lan: "Qh1xg2", legalMoves: []Move{newCaptureMove(H1, G2, WhiteQueen), newCaptureMove(H3, G2, WhiteQueen)}},
		"4 queens ambiguity":                    {m: newQuietMove(G2, F3, WhiteQueen), uci: "g2f3", san: "Qg2f3", lan: "Qg2-f3", legalMoves: []Move{newQuietMove(G2, F3, WhiteQueen), newQuietMove(G4, F3, WhiteQueen), newQuietMove(E2, F3, WhiteQueen), newQuietMove(E4, F3, WhiteQueen)}},
		"4 queens ambiguity capture":            {m: newCaptureMove(G2, F3, WhiteQueen), uci: "g2f3", san: "Qg2xf3", lan: "Qg2xf3", legalMoves: []Move{newCaptureMove(G2, F3, WhiteQueen), newCaptureMove(G4, F3, WhiteQueen), newCaptureMove(E2, F3, WhiteQueen), newCaptureMove(E4, F3, WhiteQueen)}},
		"ambiguity pawn capture":                {m: newCaptureMove(E2, F3, WhitePawn), uci: "e2f3", san: "exf3", lan: "e2xf3", legalMoves: []Move{newCaptureMove(E2, F3, WhitePawn), newCaptureMove(G2, F3, WhitePawn)}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.m.ToUCI() != tc.uci {
				t.Fatalf("invalid move uci: want %s, got %s", tc.uci, tc.m.ToUCI())
			}
			lanGot := tc.m.ToLAN(nil, tc.isCheck, tc.isCheckmate)
			if lanGot != tc.lan {
				t.Fatalf("invalid move lan: want %s, got %s", tc.lan, lanGot)
			}
			legals := make([]Move, 0)
			if len(tc.legalMoves) == 0 {
				legals = append(legals, tc.m)
			} else {
				legals = tc.legalMoves
			}
			sanGot := tc.m.ToSAN(nil, tc.isCheck, tc.isCheckmate, legals)
			if sanGot != tc.san {
				t.Fatalf("invalid move san: want %s, got %s", tc.san, sanGot)
			}
		})
	}
}
