package juicer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestPerftNodes(t *testing.T) {
	testCases := map[string]struct {
		fen       string
		depth     int
		wantNodes int64
	}{
		"Pos1":  {fen: "r6r/1b2k1bq/8/8/7B/8/8/R3K2R b KQ - 3 2", depth: 1, wantNodes: 8},
		"Pos2":  {fen: "8/8/8/2k5/2pP4/8/B7/4K3 b - d3 0 3", depth: 1, wantNodes: 8},
		"Pos3":  {fen: "r1bqkbnr/pppppppp/n7/8/8/P7/1PPPPPPP/RNBQKBNR w KQkq - 2 2", depth: 1, wantNodes: 19},
		"Pos4":  {fen: "r3k2r/p1pp1pb1/bn2Qnp1/2qPN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQkq - 3 2", depth: 1, wantNodes: 5},
		"Pos5":  {fen: "2kr3r/p1ppqpb1/bn2Qnp1/3PN3/1p2P3/2N5/PPPBBPPP/R3K2R b KQ - 3 2", depth: 1, wantNodes: 44},
		"Pos6":  {fen: "rnb2k1r/pp1Pbppp/2p5/q7/2B5/8/PPPQNnPP/RNB1K2R w KQ - 3 9", depth: 1, wantNodes: 39},
		"Pos7":  {fen: "2r5/3pk3/8/2P5/8/2K5/8/8 w - - 5 4", depth: 1, wantNodes: 9},
		"Pos8":  {fen: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8", depth: 3, wantNodes: 62379},
		"Pos9":  {fen: "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", depth: 3, wantNodes: 89890},
		"Pos10": {fen: "3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", depth: 6, wantNodes: 1134888},
		"Pos11": {fen: "8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", depth: 6, wantNodes: 1015133},
		"Pos12": {fen: "8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", depth: 6, wantNodes: 1440467},
		"Pos13": {fen: "5k2/8/8/8/8/8/8/4K2R w K - 0 1", depth: 6, wantNodes: 661072},
		"Pos14": {fen: "3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", depth: 6, wantNodes: 803711},
		"Pos15": {fen: "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", depth: 4, wantNodes: 1274206},
		"Pos16": {fen: "r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", depth: 4, wantNodes: 1720476},
		"Pos17": {fen: "2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", depth: 6, wantNodes: 3821001},
		"Pos18": {fen: "8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", depth: 5, wantNodes: 1004658},
		"Pos19": {fen: "4k3/1P6/8/8/8/8/K7/8 w - - 0 1", depth: 6, wantNodes: 217342},
		"Pos20": {fen: "8/P1k5/K7/8/8/8/8/8 w - - 0 1", depth: 6, wantNodes: 92683},
		"Pos21": {fen: "K1k5/8/P7/8/8/8/8/8 w - - 0 1", depth: 6, wantNodes: 2217},
		"Pos22": {fen: "8/k1P5/8/1K6/8/8/8/8 w - - 0 1", depth: 7, wantNodes: 567584},
		"Pos23": {fen: "8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", depth: 4, wantNodes: 23527},
	}

	InitPrecalculatedTables()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			p := &Position{}
			if err := p.LoadFromFEN(tc.fen); err != nil {
				t.Fatal(err)
			}

			pd := perft2(tc.fen, tc.depth)
			if int64(pd.Nodes) != tc.wantNodes {
				t.Fatalf("want %v, got %v", tc.wantNodes, pd.Nodes)
			}
		})
	}
}

var c int

func play(p *Position, m Move) {
	type tmpp struct {
		TURN  Color
		ENP   Square
		CR    CastleRights
		CHECK bool
		INSF  bool
	}

	pre := tmpp{TURN: p.turn, ENP: p.enpSquare, CR: p.castleRights, CHECK: p.check, INSF: p.insufficientMaterial}
	bbpre, _ := json.Marshal(&pre)

	p.MakeMove(m)
	c++

	after := tmpp{TURN: p.turn, ENP: p.enpSquare, CR: p.castleRights, CHECK: p.check, INSF: p.insufficientMaterial}
	bbafter, _ := json.Marshal(&after)

	num := fmt.Sprintf("%v", c)
	if c < 10 {
		num = fmt.Sprintf(" %v", c)
	}

	fmt.Printf("\n%v. %+v  %+v\n          %+v\n%+v\n%+v\n", num, m, string(bbpre), string(bbafter), p.PrintBoard(), p.Fen())
}

func dbg(w io.Writer, fen string, depth int) {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	pseudo := p.generateAllPseudoLegalMoves()
	legal := make([]Move, 0)

	sort.Slice(pseudo, func(i, j int) bool {
		if pseudo[i].String()[0] != pseudo[j].String()[0] {
			return pseudo[i].String()[0] < pseudo[j].String()[0]
		}
		return pseudo[i].String()[1:] < pseudo[j].String()[1:]
	})

	var nodesSearched int64

	mine := make(map[string]int64)
	sf := make(map[string]int64)
	diff := make(map[string]int64)

	for _, m := range pseudo {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			legal = append(legal, m)
			nodes := traverse(p, depth-1)
			nodesSearched += nodes
			mine[m.String()] = nodes
		}

		unmakeMove()
	}

	mine["total"] = nodesSearched

	cmd := exec.Cmd{Path: "/usr/bin/stockfish"}

	in, err := cmd.StdinPipe()
	if err != nil {
		panic("in err: " + err.Error())
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		panic("out err: " + err.Error())
	}
	defer out.Close()

	if err := cmd.Start(); err != nil {
		panic("start err: " + err.Error())
	}

	send := func(command string) {
		if _, err := in.Write([]byte(command + "\n")); err != nil {
			in.Close()
		}
	}

	send("position fen " + fen)
	send("go perft " + strconv.Itoa(depth))

	in.Close()

	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		txt := scanner.Text()

		if strings.HasPrefix(txt, "Nodes searched") {
			pair := strings.Split(txt, ": ")
			nodes, _ := strconv.Atoi(pair[1])
			sf["total"] = int64(nodes)
		}

		if !strings.HasPrefix(txt, "Stockfish") &&
			!strings.HasPrefix(txt, "Fen") &&
			!strings.HasPrefix(txt, "Key") &&
			!strings.HasPrefix(txt, "Checkers") &&
			!strings.HasPrefix(txt, "Nodes searched") &&
			txt != "" {
			pair := strings.Split(txt, ": ")
			move := pair[0]
			nodes, _ := strconv.Atoi(pair[1])
			sf[move] = int64(nodes)
		}
	}

	for m, n := range mine {
		if m == "total" {
			continue
		}
		if sf[m] != n {
			diff[m] = n
		}
	}

	cmd.Wait()

	display := func() {
		fmt.Printf("%6v | %6v | %6v | %3v | %4v", "move", "sf", "own", "dif", "n")
		fmt.Printf("\n  ------------------------------------\n")

		isDiff := func(a, b int64) string {
			if a != b {
				return "*"
			}
			return " "
		}

		for k, v := range sf {
			if k != "total" {
				fmt.Printf("%6v | %6v | %6v | %3v | %4v\n", k, v, mine[k], isDiff(v, mine[k]), v-mine[k])
			}
		}

		fmt.Printf("\n  Nodes searched: %6v | %6v | %6v\n\n", sf["total"], mine["total"], sf["total"]-mine["total"])
		fmt.Printf("\nlegal:%v\n%v\n", len(legal), legal)
	}

	display()
}

func TestJuicer(t *testing.T) {
	InitPrecalculatedTables()

	fen := "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"

	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		t.Fatal(err)
	}

	dbg(os.Stdout, fen, 5)
}
