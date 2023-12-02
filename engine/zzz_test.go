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

	fmt.Println(p.PrintBoard())

	var c int
	play := func(m Move) {
		p.MakeMove(m)
		c++
		fmt.Printf("%v.\n %+v\n %+v\n %+v\n", c, m, p, p.PrintBoard())
	}
	_ = play

	// play(NewMove(E2, E4, WhitePawn, PromotionNone, false, true, false, false))
	// play(NewMove(E7, E5, BlackPawn, PromotionNone, false, true, false, false))
	// play(NewMove(G1, F3, WhiteKnight, PromotionNone, false, false, false, false))
	// play(NewMove(B8, C6, BlackKnight, PromotionNone, false, false, false, false))
	// play(NewMove(F1, B5, WhiteBishop, PromotionNone, false, false, false, false))

	pseudo := p.generateAllPseudoLegalMoves()
	legal := p.generateAllLegalMoves(pseudo)

	fmt.Printf("pseudo: %v: %+v\n", len(pseudo), pseudo)
	fmt.Printf(" legal: %v: %+v\n", len(legal), legal)
}
