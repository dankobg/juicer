package juicer

import (
	"testing"
)

func TestJuicer(t *testing.T) {
	wk := bitboardEmpty
	wk.setBit(E1)
	wq := bitboardEmpty
	wq.setBit(E2)
	wb := bitboardEmpty
	wb.setBit(C1)
	// wb.setBit(F7)
	wn := bitboardEmpty
	// wn.setBit(B1)
	// wn.setBit(D4)
	wr := bitboardEmpty
	wr.setBit(A1)
	wr.setBit(H1)
	wp := bitboardEmpty
	wp.setBit(A2)
	wp.setBit(A3)
	wp.setBit(C2)
	wp.setBit(D2)
	wp.setBit(E4)
	wp.setBit(F2)
	wp.setBit(G3)
	wp.setBit(H2)

	bk := bitboardEmpty
	bk.setBit(F7)
	bq := bitboardEmpty
	bq.setBit(F6)
	bb := bitboardEmpty
	bb.setBit(C8)
	// bb.setBit(F8)
	bn := bitboardEmpty
	// bn.setBit(D5)
	bn.setBit(G8)
	br := bitboardEmpty
	br.setBit(A8)
	br.setBit(H8)
	bp := bitboardEmpty
	bp.setBit(A6)
	bp.setBit(B5)
	bp.setBit(C7)
	bp.setBit(D7)
	bp.setBit(D4)
	bp.setBit(F7)
	bp.setBit(G7)
	bp.setBit(H7)

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
}
