package juicer

import "testing"

func TestValidateFen(t *testing.T) {
	testCases := map[string]struct {
		fen     string
		wantErr bool
	}{
		// "valid empty fen":    {fen: FENEmptyPosition, wantErr: false},
		// "valid starting fen": {fen: FENStartingPosition, wantErr: false},
		// "valid Ruy Lopez":               {fen: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/R1BQKBNR b KQkq - 0 3", wantErr: false},
		// "valid Sicilian Defense":        {fen: "rnbqkb1r/ppp1pppp/8/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
		// "valid King's Indian Defense":   {fen: "rnbq1rk1/ppp1ppbp/3p1np1/8/3PP3/2N2N2/PPP2PPP/R1BQKB1R w KQ - 0 6", wantErr: false},
		// "valid Queen's Gambit Declined": {fen: "rnbqkbnr/pp2pppp/2p5/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
		// "valid Italian Game":            {fen: "rnbqkbnr/pp1ppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 4", wantErr: false},
		// "valid Caro-Kann Defense":       {fen: "rnbqkbnr/pp1ppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 4", wantErr: false},
		// "valid English Opening":         {fen: "rnbqkbnr/pppppppp/8/3P4/3P4/8/PPP2PPP/R1BQKBNR b KQkq - 0 3", wantErr: false},
		// "valid Scandinavian Defense":    {fen: "rnbqkbnr/ppp1pppp/8/3p4/3P4/8/PPP2PPP/R1BQKBNR w KQkq - 0 4", wantErr: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := validateFEN(tc.fen, validateFenOps{})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error mismatch, wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
		})
	}
}
