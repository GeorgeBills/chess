package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

const maxDepth = 5
const divide = true
const fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
	if err != nil {
		fatal(err)
	}
	g := engine.NewGame(b)
	start := time.Now()
	n := perft(g, maxDepth, divide)
	elapsed := time.Since(start)
	fmt.Printf("%d nodes, %dms\n", n, elapsed.Milliseconds())
}

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}

func perft(g engine.Game, depth uint8, divide bool) uint64 {
	var ret uint64
	moves := make([]engine.Move, 0, 32)
	moves, _ = g.GenerateMoves(moves)
	if depth == 1 {
		return uint64(len(moves))
	}
	for _, move := range moves {
		g.MakeMove(move)
		n := perft(g, depth-1, false)
		if divide {
			fmt.Printf("%s %d\n", move.SAN(), n)
		}
		ret += n
		g.UnmakeMove()
	}
	return ret
}
