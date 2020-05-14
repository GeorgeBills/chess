package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	gteUCI        = "uci"        // tell engine to use the universal chess interface
	gteDebug      = "debug"      // switch the debug mode of the engine on and off
	gteIsReady    = "isready"    // used to synchronize the engine with the GUI
	gteSetOption  = "setoption"  // change internal parameters of the engine
	gteNewGame    = "ucinewgame" // the next search will be from a different game
	gtePosition   = "position"   // set up the position described on the internal board
	gteGo         = "go"         // start calculating on the current position
	gteGoDepth    = "depth"      // search x plies only
	gteGoNodes    = "nodes"      // search x nodes only
	gteGoMoveTime = "movetime"   // search exactly x mseconds
	gteGoInfinite = "infinite"   // search until the stop command
	gteStop       = "stop"       // stop calculating as soon as possible
	gteQuit       = "quit"       // quit the program as soon as possible
)

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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
SCAN_INPUT:
	for scanner.Scan() {
		word := scanner.Text()
		switch word {
		case gteUCI:
			id()
			uciok()
		case gteIsReady:
			readyok()
		case gteQuit:
			break SCAN_INPUT
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func id() {
	fmt.Println(etgID, etgIDName, name)
	fmt.Println(etgID, etgIDAuthor, author)
}

func uciok() {
	fmt.Println(etgUCIOK)
}

func readyok() {
	// TODO: block on mutex in engine if we're waiting on anything slow?
	fmt.Println(etgReadyOK)
}
