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

func main() {
	h := &handler{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	// we parse UCI with a function-to-function state machine as described in
	// the talk "Lexical Scanning in Go" by Rob Pike
	// (https://youtu.be/HxaD_trXwRE). each state func returns the next state
	// func we are transitioning to.
	for state := waitingForUCI(h, scanner); state != nil; {
		state = state(h, scanner)
	}
}

type statefn func(h *handler, scanner *bufio.Scanner) statefn

func waitingForUCI(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	if err := scanner.Err(); err != nil {
		return errorScanning(err)
	}
	text := scanner.Text()
	switch text {
	case gteUCI:
		h.Identify()
		fmt.Println(etgUCIOK)
		return waitingForCommand
	case gteQuit:
		return nil // no further states
	default:
		logger.Printf("unrecognized: %s\n", text)
		return waitingForUCI
	}
}

func waitingForCommand(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	if err := scanner.Err(); err != nil {
		return errorScanning(err)
	}
	text := scanner.Text()
	switch text {
	case gteIsReady:
		h.IsReady()
		return waitingForCommand
	case gteNewGame:
		h.NewGame()
		return waitingForCommand
	case gtePosition:
		return positionCommand
	case gteGo:
		return goCommand
	case gteQuit:
		return nil
	default:
		return errorUnrecognized(text, waitingForCommand)
	}
}

func positionCommand(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	_ = scanner.Err()
	fen := scanner.Text()
	if fen == "startpos" {
		h.SetStartingPosition()
	}
	// FIXME: fen isn't a single word...
	// b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	return waitingForCommand
}

func goCommand(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	_ = scanner.Err()
	text := scanner.Text()
	switch text {
	case "depth":
		_ = scanner.Scan()
		_ = scanner.Err()
		plies, _ := strconv.Atoi(scanner.Text())
		if 0 < plies && plies < maxDepthPlies {
			h.GoDepth(uint8(plies))
		}
	}
	return waitingForCommand
}

func errorUnrecognized(text string, next statefn) statefn {
	logger.Printf("unrecognized: %s\n", text)
	return next
}

func errorScanning(err error) statefn {
	logger.Printf("error scanning input: %v", err)
	return nil
}

type handler struct {
	game *engine.Game
}

func (h *handler) Identify() {
	fmt.Println(etgID, etgIDName, name)
	fmt.Println(etgID, etgIDAuthor, author)
}

func (h *handler) IsReady() {
	// TODO: block on mutex in engine if we're waiting on anything slow?
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
