package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	logf, err := os.Create("uci.log")
	if err != nil {
		fatal(err)
	}

	logger := log.New(logf, "", 0)

	h := &handler{
		logger: logger,
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	parser := &parser{
		logger:  logger,
		handler: h,
		scanner: scanner,
	}

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
