package main

import (
	"fmt"

	juicer "github.com/dankobg/juicer/engine"
)

func main() {
	juicer.InitPrecalculatedTables()

	p := juicer.Position{}

	fen := juicer.FENStartingPosition

	if err := p.LoadFromFEN(fen); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", p.PrintBoard())
	}
}
