package main

import (
	"fmt"
	"os"
)

func main() {
	logf, err := os.Create("uci.log")
	if err != nil {
		fatal(err)
	}

	h := NewHandler(logf)
	parser := NewParser(os.Stdin, h, logf)
	parser.Run()
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
