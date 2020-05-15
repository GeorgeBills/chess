package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

var logger = log.New(os.Stderr, "", 0)

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
	name          = "github.com/GeorgeBills/chess"
	author        = "George Bills"
	maxDepthPlies = 40
)

// TODO: state machine
//       http://denis.papathanasiou.org/archive/2013.02.10.post.pdf
//       https://www.youtube.com/watch?v=HxaD_trXwRE

func main() {
	h := &handler{}

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
		case gteNewGame:
			h.NewGame()
		case gtePosition:
			position(scanner, h)
		case gteGo:
			ucigo(scanner, h)
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

func ucinewgame(h *handler) {
	h.NewGame()
}

func position(scanner *bufio.Scanner, h *handler) {
	_ = scanner.Scan()
	_ = scanner.Err()
	fen := scanner.Text()
	if fen == "startpos" {
		h.SetStartingPosition()
	}
	// FIXME: fen isn't a single word...
	// b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
}

func ucigo(scanner *bufio.Scanner, h *handler) {
	_ = scanner.Scan()
	_ = scanner.Err()
	which := scanner.Text()
	switch which {
	case "depth":
		_ = scanner.Scan()
		_ = scanner.Err()
		plies, _ := strconv.Atoi(scanner.Text())
		if 0 < plies && plies < maxDepthPlies {
			h.GoDepth(uint8(plies))
		}
	}
}

type handler struct {
	game *engine.Game
}

func (h *handler) NewGame() {
	logger.Println("initialised new game")
	g := engine.NewGame(nil) // TODO: return pointer
	h.game = &g
}

func (h *handler) SetStartingPosition() {
	logger.Println("set starting position")
	b := engine.NewBoard()
	h.game.SetBoard(&b)
}

func (h *handler) SetPosition(fen string) {
	panic("SetPosition not implemented")
}

func (h *handler) GoDepth(plies uint8) {
	m, _ := h.game.BestMoveToDepth(plies * 2)
	fmt.Println(etgBestMove, m.SAN())
}
