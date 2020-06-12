package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GeorgeBills/chess/engine"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Errorf("%s <fen>", os.Args[0]))
	}

	board, err := engine.NewBoardFromFEN(strings.NewReader(os.Args[1]))
	if err != nil {
		fatal(err)
	}

	fmt.Print(board.String())
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
