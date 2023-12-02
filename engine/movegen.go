package juicer

func (p *Position) generateAllLegalMoves(pseudoMoves []Move) []Move {
	legalMoves := make([]Move, 0)

	for _, m := range pseudoMoves {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			legalMoves = append(legalMoves, m)
		}

		unmakeMove()
	}

	return legalMoves
}
