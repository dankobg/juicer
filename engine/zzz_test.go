package engine

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

	bb := bitboardEmpty

	bb.setBit(A1)
	bb.setBit(A8)
	bb.setBit(H1)
	bb.setBit(H8)
	bb.setBit(B4)
	bb.setBit(E7)
	bb.setBit(F6)

	fmt.Println("count: ", bb.countBits())
	fmt.Println(bb.drawPretty())
}
