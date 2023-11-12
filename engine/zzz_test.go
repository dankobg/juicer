package juicer

import (
	"fmt"
	"testing"
)

func TestJuicer(t *testing.T) {
	// for k, v := range bitboardEmptyFiles {
	// 	fmt.Printf("\n    EMPTY FILES: %v\n", k)
	// 	fmt.Printf("%v\n", v.drawPretty())
	// }
	// for k, v := range bitboardUniverseFiles {
	// 	fmt.Printf("\n    UNIVERSE FILES: %v\n", k)
	// 	fmt.Printf("%v\n", v.drawPretty())
	// }
	// for k, v := range bitboardEmptyRanks {
	// 	fmt.Printf("\n    EMPTY RANKS: %v\n", k)
	// 	fmt.Printf("%v\n", v.drawPretty())
	// }
	// for k, v := range bitboardUniverseRanks {
	// 	fmt.Printf("\n    UNIVERSE RANKS: %v\n", k)
	// 	fmt.Printf("%v\n", v.drawPretty())
	// }

	wk := bitboardEmpty
	wk.setBit(E1)

	wq := bitboardEmpty
	wq.setBit(D1)

	wr := bitboardEmpty
	wr.setBit(A1, H1)

	wb := bitboardEmpty
	wb.setBit(C1, F1)

	wn := bitboardEmpty
	wn.setBit(B1, G1)

	wp := bitboardEmpty
	wp.setBit(A2, B2, C2, D2, E2, F2, G2, H2)

	bk := bitboardEmpty
	bk.setBit(E8)

	bq := bitboardEmpty
	bq.setBit(D8)

	br := bitboardEmpty
	br.setBit(A8, H8)

	bb := bitboardEmpty
	bb.setBit(C8, F8)

	bn := bitboardEmpty
	bn.setBit(B8, G8)

	bp := bitboardEmpty
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
		whitePiecesOccupancy:  bitboardEmpty,
		blackPiecesOccupancy:  bitboardEmpty,
		allPiecesOccupancy:    bitboardEmpty,
	}

	fmt.Println(wp.drawCompact())
	fmt.Println(wp.drawPretty())

	fmt.Println(b.DrawCompact())
	fmt.Println(b.DrawPretty())
}
