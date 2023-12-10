package main

import (
	"fmt"
	"os"

	juicer "github.com/dankobg/juicer/engine"
)

func main() {
	juicer.InitPrecalculatedTables()

	p := juicer.Position{}

	fen := juicer.FENStartingPosition

	if err := p.LoadFromFEN(fen); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", p.PrintBoard())
}
