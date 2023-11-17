package juicer

import (
	"fmt"
	"testing"
)

func TestJuicer(t *testing.T) {
	var wk, wq, wr, wb, wn, wp, bk, bq, br, bb, bn, bp bitboard

	wk.setBit(E1)
	wq.setBit(D1)
	wr.setBit(A1, H1)
	wb.setBit(C1, F1)
	wn.setBit(B1, G1)
	wp.setBit(A2, B2, C2, D2, E2, F2, G2, H2)
	bk.setBit(E8)
	bq.setBit(D8)
	br.setBit(A8, H8)
	bb.setBit(C8, F8)
	bn.setBit(B8, G8)
	bp.setBit(A7, B7, C7, D7, E7, F7, G7, H7)

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

	var bbtest bitboard
	bbtest.setBit(D2, D7, B4, G4)

	fmt.Println(bbtest.draw(nil))
	fmt.Println("LSB: ", bbtest.LS1B())
	fmt.Println("MSB: ", bbtest.MS1B())
}
