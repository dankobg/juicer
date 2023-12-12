package juicer

import "testing"

func TestValidateFen(t *testing.T) {
	testCases := map[string]struct {
		fen     string
		wantErr bool
	}{
		"valid starting fen":            {fen: FENStartingPosition, wantErr: false},
		"valid Ruy Lopez":               {fen: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/R1BQKBNR b KQkq - 0 3", wantErr: false},
		"valid Sicilian Defense":        {fen: "rnbqkb1r/ppp1pppp/8/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
		"valid King's Indian Defense":   {fen: "rnbq1rk1/ppp1ppbp/3p1np1/8/3PP3/2N2N2/PPP2PPP/R1BQKB1R w KQ - 0 6", wantErr: false},
		"valid Queen's Gambit Declined": {fen: "rnbqkbnr/pp2pppp/2p5/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
		"valid Italian Game":            {fen: "rnbqkbnr/pp1ppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 4", wantErr: false},
		"valid Caro-Kann Defense":       {fen: "rnbqkbnr/pp1ppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 4", wantErr: false},
		"valid English Opening":         {fen: "rnbqkbnr/pppppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 3", wantErr: false},
		"valid Scandinavian Defense":    {fen: "rnbqkbnr/ppp1pppp/8/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
		"valid random 1":                {fen: "r6r/1b2k1bq/8/8/7B/8/8/R3K2R b KQ - 3 2", wantErr: false},
		"valid random 2":                {fen: "8/8/8/2k5/2pP4/8/B7/4K3 b - d3 0 3", wantErr: false},
		"valid random 3":                {fen: "r1bqkbnr/pppppppp/n7/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 2 2", wantErr: false},
		"valid random 4":                {fen: "r3k2r/p1pp1pb1/bn2Qnp1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQkq - 3 2", wantErr: false},
		"valid random 5":                {fen: "2kr3r/p1ppqpb1/bn2Qnp1/3PN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQ - 3 2", wantErr: false},
		"valid random 6":                {fen: "rnb2k1r/pp1Pbppp/2p5/q7/2B5/8/PPPQNnPP/RNB1K2R w KQ - 3 9", wantErr: false},
		"valid random 7":                {fen: "2r5/3pk3/8/2P5/8/2K5/8/8 w - - 5 4", wantErr: false},
		"valid random 8":                {fen: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", wantErr: false},
		"valid random 9":                {fen: "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", wantErr: false},
		"valid random 10":               {fen: "3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", wantErr: false},
		"valid random 11":               {fen: "8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", wantErr: false},
		"valid random 12":               {fen: "8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", wantErr: false},
		"valid random 13":               {fen: "5k2/8/8/8/8/8/8/4K2R w K - 0 1", wantErr: false},
		"valid random 14":               {fen: "3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", wantErr: false},
		"valid random 15":               {fen: "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", wantErr: false},
		"valid random 16":               {fen: "r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", wantErr: false},
		"valid random 17":               {fen: "2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", wantErr: false},
		"valid random 18":               {fen: "8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", wantErr: false},
		"valid random 19":               {fen: "4k3/1P6/8/8/8/8/K7/8 w - - 0 1", wantErr: false},
		"valid random 20":               {fen: "8/P1k5/K7/8/8/8/8/8 w - - 0 1", wantErr: false},
		"valid random 21":               {fen: "K1k5/8/P7/8/8/8/8/8 w - - 0 1", wantErr: false},
		"valid random 22":               {fen: "8/k1P5/8/1K6/8/8/8/8 w - - 0 1", wantErr: false},
		"valid random 23":               {fen: "8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", wantErr: false},
		"invalid empty fen":             {fen: FENEmptyPosition, wantErr: true},
		"invalid fen delim length":      {fen: " rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w  KQkq - 0 1 ", wantErr: true},
		"invalid fen active turn W":     {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR W KQkq - 0 1", wantErr: true},
		"invalid fen active turn B":     {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR B KQkq - 0 1", wantErr: true},
		"invalid fen active turn -":     {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR - KQkq - 0 1", wantErr: true},
		"invalid fen half move clock":   {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - -1 1", wantErr: true},
		"invalid fen full move clock 0": {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0", wantErr: true},
		"invalid fen enp square":        {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq k3 0 1", wantErr: true},
		"invalid fen castle junk":       {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w zz - 0 1", wantErr: true},
		"invalid fen castle order QKqk": {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w QKqk - 0 1", wantErr: true},
		"invalid fen castle order QKq":  {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w QKq - 0 1", wantErr: true},
		"invalid fen castle order QKk":  {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w QKk - 0 1", wantErr: true},
		"invalid fen castle order Qqk":  {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Qqk - 0 1", wantErr: true},
		"invalid fen castle order Kqk":  {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kqk - 0 1", wantErr: true},
		"invalid fen castle order qk":   {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w qk - 0 1", wantErr: true},
		"invalid fen castle order QK":   {fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w QK - 0 1", wantErr: true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := validateFEN(tc.fen, validateFenOps{})
			if (err != nil) != tc.wantErr {
				t.Fatalf("invalid fen, error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
		})
	}
}
