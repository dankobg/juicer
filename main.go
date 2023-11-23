package main

import (
	"fmt"

	juicer "github.com/dankobg/juicer/engine"
)

func main() {
	fen := "r1b1kb1r/2ppqppp/2n2n2/p3p3/B3P3/5N2/1PPPQPPP/RNB1K2R w KQkq - 0 1"
	// fen := juicer.FENStartingPosition
	p := juicer.Position{}

	if err := p.LoadFromFEN(fen); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", p.PrintBoard())
	}
}
