package main

import (
	"bufio"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	h := &handler{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	// we parse UCI with a func-to-func state machine as described in the talk
	// "Lexical Scanning in Go" by Rob Pike (https://youtu.be/HxaD_trXwRE). each
	// state func returns the next state func we are transitioning to.
	for state := waitingForUCI(h, scanner); state != nil; {
		state = state(h, scanner)
	}
}
