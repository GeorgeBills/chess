package main

import (
	"fmt"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

const maxDepth = 5

func main() {
	b := engine.NewBoard()
	g := engine.NewGame(&b)
	n := perft(g, maxDepth)
	fmt.Printf("%d\n", n)
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
