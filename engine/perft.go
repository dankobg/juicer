package juicer

import (
	"fmt"
	"sort"
	"time"
)

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
