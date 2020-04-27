package main

import (
	"bufio"
	"fmt"
	"github.com/GeorgeBills/chess/m/v2/engine"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Sprintf("%s <fen>", os.Args[0]))
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

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}
