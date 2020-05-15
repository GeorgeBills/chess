package main

import (
	"bufio"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

const (
	etgID       = "id"       // sent to identify the engine
	etgIDName   = "name"     // e.g. "id name Shredder X.Y\n"
	etgIDAuthor = "author"   // e.g. "id author Stefan MK\n"
	etgUCIOK    = "uciok"    // the engine has sent all infos and is ready
	etgReadyOK  = "readyok"  // the engine is ready to accept new commands
	etgBestMove = "bestmove" // engine has stopped searching and found the best move
	etgInfo     = "info"     // engine wants to send information to the GUI
)

const (
	name   = "github.com/GeorgeBills/chess"
	author = "George Bills"
)

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
