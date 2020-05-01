package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

const defaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	depth := flag.Uint("depth", 1, "depth to generate moves to")
	fen := flag.String("fen", defaultFEN, "FEN to start with")
	divide := flag.Bool("divide", false, "whether or not to output node count divided by initial moves")

	flag.Parse()

	b, err := engine.NewBoardFromFEN(strings.NewReader(*fen))
	if err != nil {
		fatal(err)
	}
	g := engine.NewGame(b)
	start := time.Now()
	n := perft(g, *depth, *divide)
	elapsed := time.Since(start)
	fmt.Printf("%d nodes, %dms\n", n, elapsed.Milliseconds())
}

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}

func perft(g engine.Game, depth uint, divide bool) uint64 {
	var ret uint64
	moves := make([]engine.Move, 0, 32)
	moves, _ = g.GenerateMoves(moves)
	if depth <= 1 {
		return uint64(len(moves))
	}
	for _, move := range moves {
		g.MakeMove(move)
		n := perft(g, depth-1, false)
		if divide {
			fmt.Printf("%s\t%d\t%s\n", move.SAN(), n, g.FEN())
		}
		ret += n
		g.UnmakeMove()
	}
	return ret
}
