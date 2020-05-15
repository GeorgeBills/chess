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

	// we parse UCI with a func-to-func state machine as described in the talk
	// "Lexical Scanning in Go" by Rob Pike (https://youtu.be/HxaD_trXwRE). each
	// state func returns the next state func we are transitioning to.
	for state := waitingForUCI(parser); state != nil; {
		state = state(parser)
	}
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
