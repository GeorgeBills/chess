package main

import (
	"fmt"
	"os"

	"github.com/GeorgeBills/chess/m/v2/uci"
)

func main() {
	logf, err := os.Create("uci.log")
	if err != nil {
		fatal(err)
	}

	h := newHandler(logf)
	parser := uci.NewParser(h, os.Stdin, os.Stdout, logf)
	parser.Run()
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
