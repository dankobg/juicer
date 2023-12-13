package juicer

import (
	"fmt"
	"sort"
	"time"
)

func traverse(p *Position, depth int) int64 {
	if depth == 0 {
		return 1
	}

	var num int64
	pseudo := p.generateAllPseudoLegalMoves()

	for i := 0; i < len(pseudo); i++ {
		unmakeMove := p.MakeMove(pseudo[i])

		if !p.board.IsInCheck(p.turn.Opposite()) {
			num += traverse(p, depth-1)
		}

		unmakeMove()
	}

	return num
}

func perft(fen string, depth int) (int64, time.Duration) {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	start := time.Now()
	nodes := traverse(p, depth)

	return nodes, time.Since(start)
}

func perftDivide(fen string, depth int) {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	pseudo := p.generateAllPseudoLegalMoves()

	sort.Slice(pseudo, func(i, j int) bool {
		if pseudo[i].String()[0] != pseudo[j].String()[0] {
			return pseudo[i].String()[0] < pseudo[j].String()[0]
		}
		return pseudo[i].String()[1:] < pseudo[j].String()[1:]
	})

	var nodesSearched int64

	for _, m := range pseudo {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			nodes := traverse(p, depth-1)
			nodesSearched += nodes
			fmt.Printf("%v: %v\n", m, nodes)
		}

		unmakeMove()
	}

	fmt.Printf("\nNodes searched: %d\n\n", nodesSearched)
}

type perftData struct {
	Nodes      int
	Captures   int
	Enpassants int
	Castles    int
	Promotions int
	Checks     int
}

func traverse2(p *Position, depth int, pd *perftData) {
	if depth == 0 {
		pd.Nodes++
		return
	}

	pseudo := p.generateAllPseudoLegalMoves()

	for i := 0; i < len(pseudo); i++ {
		m := pseudo[i]
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			traverse2(p, depth-1, pd)

			if m.IsCapture() {
				pd.Captures++
			}
			if m.IsEnPassant() {
				pd.Enpassants++
			}
			if m.IsCastle() && m.Piece().IsKing() {
				pd.Castles++
			}
			if m.Promotion().IsPromotion() {
				pd.Promotions++
			}
			if p.board.IsInCheck(p.turn) {
				pd.Checks++
			}
		}

		unmakeMove()
	}
}

func perft2(fen string, depth int) perftData {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	var pd perftData
	traverse2(p, depth, &pd)
	return pd
}

func perftDivide2(fen string, depth int) {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	pseudo := p.generateAllPseudoLegalMoves()

	sort.Slice(pseudo, func(i, j int) bool {
		if pseudo[i].String()[0] != pseudo[j].String()[0] {
			return pseudo[i].String()[0] < pseudo[j].String()[0]
		}
		return pseudo[i].String()[1:] < pseudo[j].String()[1:]
	})

	var nodesSearched int64
	var pd perftData

	for _, m := range pseudo {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			traverse2(p, depth-1, &pd)
			fmt.Printf("%v: %+v\n", m, pd)
		}

		unmakeMove()
	}

	fmt.Printf("\nNodes searched: %d\n\n", nodesSearched)
}
