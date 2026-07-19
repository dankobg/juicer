package engine

// perft tests can be very slow even some at depth 5 can take ~1m
// they are engine and cpu bound so i would need to super optimize engine first

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"
// )

// type Result struct {
// 	FEN     string
// 	Depth   int
// 	Nodes   int
// 	Calc    int64
// 	Elapsed time.Duration
// }

// var maxDepthCheck = 4

// func main() {
// 	InitPrecalculatedTables()

// 	lines := make(chan string)
// 	results := make(chan Result)
// 	wg := new(sync.WaitGroup)

// 	go func() {
// 		processEpdFile("engine/perft_suite/epds/single_check_1.epd", lines)
// 		close(lines)
// 	}()

// 	numWorkers := 128
// 	for range numWorkers {
// 		wg.Add(1)
// 		go worker(lines, results, wg)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for result := range results {
// 		fmt.Printf("%s, D%d [%d, %d] took %v\n", result.FEN, result.Depth, result.Nodes, result.Calc, result.Elapsed)
// 	}
// }

// func processEpdFile(path string, lines chan<- string) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		log.Fatalf("open epd file err: %v", err)
// 	}
// 	defer f.Close()
// 	sc := bufio.NewScanner(f)
// 	for sc.Scan() {
// 		line := strings.TrimSpace(sc.Text())
// 		if line != "" {
// 			lines <- line
// 		}
// 	}
// 	if err := sc.Err(); err != nil {
// 		log.Fatalf("scan epd file err: %v", err)
// 	}
// }

// func worker(lines <-chan string, results chan<- Result, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for line := range lines {
// 		processLine(line, results)
// 	}
// }

// func processLine(line string, results chan<- Result) {
// 	tokens := strings.SplitN(line, ";", 2)
// 	if len(tokens) < 2 {
// 		log.Fatalf("invalid tokens len")
// 	}
// 	fen := strings.TrimSpace(tokens[0])
// 	pairs := strings.Split(tokens[1:][0], ";")
// 	pairWg := new(sync.WaitGroup)
// 	for _, pair := range pairs {
// 		pairWg.Add(1)
// 		go processPair(pair, fen, results, pairWg)
// 	}
// 	pairWg.Wait()
// }

// func processPair(pair, fen string, results chan<- Result, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	vals := strings.Split(strings.TrimSpace(pair), " ")
// 	if len(vals) != 2 {
// 		log.Fatalf("invalid depth-node values")
// 	}
// 	depthToken, nodesToken := vals[0], vals[1]
// 	if !strings.HasPrefix(depthToken, "D") {
// 		log.Fatalf("no prefix D for depth")
// 	}
// 	depth, err := strconv.Atoi(depthToken[1:])
// 	if err != nil {
// 		log.Fatalf("failed to parse depth")
// 	}
// 	nodes, err := strconv.Atoi(nodesToken)
// 	if err != nil {
// 		log.Fatalf("failed to parse nodes")
// 	}
// 	if depth > maxDepthCheck {
// 		return
// 	}
// 	start := time.Now()
// 	calculated := Perft(fen, depth)
// 	elapsed := time.Since(start)
// 	results <- Result{FEN: fen, Depth: depth, Nodes: nodes, Calc: calculated, Elapsed: elapsed}
// }
