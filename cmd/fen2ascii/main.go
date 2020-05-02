package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Errorf("%s <fen>", os.Args[0]))
	}

	board, err := engine.NewBoardFromFEN(strings.NewReader(os.Args[1]))
	if err != nil {
		fatal(err)
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	_, err = f.WriteString(board.String())
	if err != nil {
		fatal(err)
	}
}

func fatal(v error) {
	fmt.Println(v)
	os.Exit(1)
}
