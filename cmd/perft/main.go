package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GeorgeBills/chess/engine"
)

func main() {
	depth := flag.Uint("depth", 1, "depth to generate moves to")
	fen := flag.String("fen", engine.InitialBoardFEN, "FEN to start with")
	divide := flag.Bool("divide", false, "whether or not to output node count divided by initial moves")
	validate := flag.Bool("validate", false, "whether or not to validate each board state")

	flag.Parse()

	b, err := engine.NewBoardFromFEN(strings.NewReader(*fen))
	if err != nil {
		fatal(fmt.Errorf("error parsing FEN: %w", err))
	}
	g := engine.NewGame(b)
	start := time.Now()
	var n uint64 = 0
	if *depth > 0 {
		n = perft(g, *depth, *validate, *divide)
	}
	elapsed := time.Since(start)
	fmt.Printf("%d nodes, took %s\n", n, elapsed)
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func perft(g *engine.Game, depth uint, validate, divide bool) uint64 {
	var ret uint64
	moves := make([]engine.Move, 0, 32)
	moves, _ = g.GenerateLegalMoves(moves)
	for _, move := range moves {
		fen := ""
		if validate {
			fen = g.FEN()
		}
		g.MakeMove(move)
		if validate {
			err := g.Validate()
			if err != nil {
				fatal(fmt.Errorf("move %s on board '%v' results in an invalid board: %v", move.SAN(), fen, err))
			}
		}
		var n uint64 = 1
		if depth > 1 {
			n = perft(g, depth-1, validate, false)
		}
		if divide {
			fmt.Printf("%s\t%d\t%s\n", move.SAN(), n, g.FEN())
		}
		ret += n
		g.UnmakeMove()
	}
	return ret
}
