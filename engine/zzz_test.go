package juicer

import (
	"fmt"
	"testing"
)

func TestJuicer(t *testing.T) {
	InitAllAttackMasksTables()

	p := Position{}

	fen := FENStartingPosition

	if err := p.LoadFromFEN(fen); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", p.PrintBoard())

	moves := p.generateAllPseudoLegalMoves()

	fmt.Println("pseudo moves:", len(moves))
	for _, m := range moves {
		fmt.Println(m)
	}
}
