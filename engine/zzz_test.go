package juicer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestJuicer(t *testing.T) {
	InitPrecalculatedTables()

	const maxDepthToCheck = 5

	f, err := os.Open("perft_suite/epds/perft.epd")
	if err != nil {
		t.Fatalf("open perft file err: %v", err)
	}

	t.Cleanup(func() {
		f.Close()
	})

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := sc.Text()

		t.Run(line, func(t *testing.T) {
			t.Parallel()

			if line == "" {
				return
			}

			tokens := strings.Split(line, ";")

			if len(tokens) < 2 {
				t.Fatalf("invalid perft epd format: %v", line)
			}

			fen := strings.TrimSpace(tokens[0])

			for i := 1; i < len(tokens); i++ {
				token := strings.TrimSpace(tokens[i])
				pairs := strings.Split(token, " ")

				if len(pairs) != 2 {
					t.Fatalf("invalid depth-nodes pair: %v", token)
				}

				depthToken, nodesToken := strings.TrimSpace(pairs[0]), strings.TrimSpace(pairs[1])

				if len(depthToken) < 2 {
					t.Fatalf("invalid depth token length")
				}
				if !strings.HasPrefix(depthToken, "D") {
					t.Fatalf("invalid depth, does not start with `D`")
				}

				wantDepth, err := strconv.Atoi(depthToken[1:])
				if err != nil {
					t.Fatalf("invalid depth, depth is not a number")
				}
				wantNodes, err := strconv.ParseInt(nodesToken, 10, 64)
				if err != nil {
					t.Fatalf("invalid test nodes, nodes is not a number")
				}

				if wantDepth <= maxDepthToCheck {
					nodes := Perft(fen, wantDepth)
					if wantNodes != nodes {
						t.Fatalf("invalid perft result at depth %v: want %v, got %v", wantDepth, wantNodes, nodes)
					}
				}
			}
		})
	}

	if err := sc.Err(); err != nil {
		t.Fatalf("scan perft file err: %v", err)
	}
}

func TestWtf(t *testing.T) {
	InitPrecalculatedTables()

	c, err := NewChess(FENStartingPosition)
	if err != nil {
		t.Fatal(err)
	}

	c.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	c.MakeMove(newQuietMove(G8, F6, BlackKnight))
	c.MakeMove(newQuietMove(F3, G1, WhiteKnight))
	c.MakeMove(newQuietMove(F6, G8, BlackKnight))

	c.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	c.MakeMove(newQuietMove(G8, F6, BlackKnight))
	c.MakeMove(newQuietMove(F3, G1, WhiteKnight))
	c.MakeMove(newQuietMove(F6, G8, BlackKnight))

	c.MakeMove(newQuietMove(A2, A3, WhitePawn))

	c.MakeMove(newQuietMove(G8, F6, BlackKnight))
	c.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	c.MakeMove(newQuietMove(F6, G8, BlackKnight))
	c.MakeMove(newQuietMove(F3, G1, WhiteKnight))

	c.MakeMove(newQuietMove(G8, F6, BlackKnight))
	c.MakeMove(newQuietMove(G1, F3, WhiteKnight))
	c.MakeMove(newQuietMove(F6, G8, BlackKnight))
	c.MakeMove(newQuietMove(F3, G1, WhiteKnight))
}
