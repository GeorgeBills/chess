package main

import (
	"bufio"
	"fmt"
	"strconv"
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
		return nil // no further states
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
	case gteGoDepth:
		return goDepthCommand(h, scanner)
	case gteGoInfinite:
		h.GoInfinite()
	case gteGoNodes:
		return goNodesCommand(h, scanner)
	}
	return waitingForCommand
}

func goDepthCommand(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	_ = scanner.Err()
	plies, err := strconv.ParseUint(scanner.Text(), 10, 8)
	if err != nil {
		return errorParsingNumber(err, waitingForCommand)
	}
	h.GoDepth(uint8(plies))
	return waitingForCommand
}

func goNodesCommand(h *handler, scanner *bufio.Scanner) statefn {
	_ = scanner.Scan()
	_ = scanner.Err()
	nodes, err := strconv.ParseUint(scanner.Text(), 10, 64)
	if err != nil {
		return errorParsingNumber(err, waitingForCommand)
	}
	h.GoNodes(nodes)
	return waitingForCommand
}

func errorUnrecognized(text string, next statefn) statefn {
	logger.Printf("unrecognized: %s\n", text)
	return next
}

func errorParsingNumber(err error, next statefn) statefn {
	logger.Printf("error parsing number: %v", err)
	return next
}

func errorScanning(err error) statefn {
	logger.Printf("error scanning input: %v", err)
	return nil
}
