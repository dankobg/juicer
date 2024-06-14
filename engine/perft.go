package juicer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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

func Perft(fen string, depth int) int64 {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	nodes := traverse(p, depth)
	return nodes
}

func Divide(fen string, depth int) {
	p := &Position{}
	if err := p.LoadFromFEN(fen); err != nil {
		panic(err)
	}

	pseudo := p.generateAllPseudoLegalMoves()

	start := time.Now()

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

	fmt.Printf("\nNodes searched: %d\n", nodesSearched)
	fmt.Printf("perft (nps): %v\n\n", (1000000*nodesSearched)/time.Since(start).Microseconds())
}

func CompareWithStockfishPerft(fen string, depth int, sfBinaryPath *string) {
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

	mine := make(map[string]int64)
	sf := make(map[string]int64)
	diff := make(map[string]int64)

	for _, m := range pseudo {
		unmakeMove := p.MakeMove(m)

		if !p.board.IsInCheck(p.turn.Opposite()) {
			nodes := traverse(p, depth-1)
			nodesSearched += nodes
			mine[m.String()] = nodes
		}

		unmakeMove()
	}

	mine["total"] = nodesSearched

	sfPath := "/usr/bin/stockfish"
	if sfBinaryPath != nil {
		sfPath = *sfBinaryPath
	}

	cmd := exec.Command(sfPath)

	in, err := cmd.StdinPipe()
	if err != nil {
		panic("in pipe err: " + err.Error())
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		panic("out pipe err: " + err.Error())
	}
	defer out.Close()

	if err := cmd.Start(); err != nil {
		panic("cmd start err: " + err.Error())
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
		line := scanner.Text()

		if strings.HasPrefix(line, "Nodes searched") {
			pair := strings.Split(line, ": ")
			nodes, _ := strconv.Atoi(pair[1])
			sf["total"] = int64(nodes)
		}

		if !strings.HasPrefix(line, "Stockfish") &&
			!strings.HasPrefix(line, "info") &&
			!strings.HasPrefix(line, "Nodes searched") &&
			line != "" {
			pair := strings.Split(line, ": ")
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

	tw := tabwriter.NewWriter(os.Stdout, 11, 4, 1, ' ', tabwriter.Debug)

	fmt.Fprintf(tw, "+----------------------------------------------+\n")
	fmt.Fprintf(tw, "| Move\t Stockfish\t Juicer\t Diff\t\n")
	fmt.Fprintf(tw, "+----------------------------------------------+\n")

	for k, v := range sf {
		if k != "total" {
			fmt.Fprintf(tw, "| %v\t %v\t %v\t %v\t\n", k, v, mine[k], v-mine[k])
		}
	}
	fmt.Fprintf(tw, "+----------------------------------------------+\n")
	fmt.Fprintf(tw, "| %v\t %v\t %v\t %v\t\n", "Nodes", sf["total"], mine["total"], sf["total"]-mine["total"])
	fmt.Fprintf(tw, "+----------------------------------------------+\n")
	tw.Flush()
}
