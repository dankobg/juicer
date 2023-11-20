package juicer

import (
	"testing"
)

func TestJuicer(t *testing.T) {
	var wk, wq, wr, wb, wn, wp, bk, bq, br, bb, bn, bp bitboard

	wk.setBit(E1)
	wq.setBit(D1)
	wr.setBit(A1)
	wr.setBit(H1)
	wb.setBit(C1)
	wb.setBit(F1)
	wn.setBit(B1)
	wn.setBit(G1)
	for _, x := range []Square{A2, B2, C2, D2, E2, F2, G2, H2} {
		wp.setBit(x)
	}
	bk.setBit(E8)
	bq.setBit(D8)
	br.setBit(A8)
	br.setBit(H8)
	bb.setBit(C8)
	bb.setBit(F8)
	bn.setBit(B8)
	bn.setBit(G8)
	for _, x := range []Square{A7, B7, C7, D7, E7, F7, G7, H7} {
		bp.setBit(x)
	}

	b := Board{
		whiteKingOccupancy:    wk,
		whiteQueensOccupancy:  wq,
		whiteRooksOccupancy:   wr,
		whiteBishopsOccupancy: wb,
		whiteKnightsOccupancy: wn,
		whitePawnsOccupancy:   wp,
		blackKingOccupancy:    bk,
		blackQueensOccupancy:  bq,
		blackRooksOccupancy:   br,
		blackBishopsOccupancy: bb,
		blackKnightsOccupancy: bn,
		blackPawnsOccupancy:   bp,
	}

	_ = b
	// fmt.Println(b.Draw(nil))
}
