package main

import (
	"fmt"

	juicer "github.com/dankobg/juicer/engine"
)

func main() {
	fen := "2kr1r2/ppp2pp1/n4q2/3p3p/8/1PN2Q1P/P1P1PPP1/2KR1B1R w K - 0 1"
	// fen := juicer.FENStartingPosition
	p := juicer.Position{}

	if err := p.LoadFromFEN(fen); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", p.PrintBoard())
	}
}
