package main

import (
	"fmt"
	"time"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

const maxDepth = 5

func main() {
	b := engine.NewBoard()
	g := engine.NewGame(&b)
	start := time.Now()
	n := perft(g, maxDepth)
	elapsed := time.Since(start)
	fmt.Printf("%d nodes, %dms\n", n, elapsed.Milliseconds())
}

func perft(g engine.Game, depth uint8) uint64 {
	var n uint64
	moves := make([]engine.Move, 0, 32)
	moves, _ = g.GenerateMoves(moves)
	if depth == 1 {
		return uint64(len(moves))
	}
	for _, move := range moves {
		g.MakeMove(move)
		n += perft(g, depth-1)
		g.UnmakeMove()
	}
	return n
}
