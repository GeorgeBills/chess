package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
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

const (
	etgID       = "id"       // sent to identify the engine
	etgIDName   = "name"     // e.g. "id name Shredder X.Y\n"
	etgIDAuthor = "author"   // e.g. "id author Stefan MK\n"
	etgUCIOK    = "uciok"    // the engine has sent all infos and is ready
	etgReadyOK  = "readyok"  // the engine is ready to accept new commands
	etgBestMove = "bestmove" // engine has stopped searching and found the best move
	etgInfo     = "info"     // engine wants to send information to the GUI
)

// NewParser returns a new parser.
func NewParser(r io.Reader, h *handler, logw io.Writer) *parser {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	return &parser{
		handler: h,
		scanner: scanner,
		logger:  log.New(logw, "parser: ", log.LstdFlags),
	}
}

type parser struct {
	logger  *log.Logger
	handler *handler
	scanner *bufio.Scanner
}

func (p *parser) Run() {
	// we parse UCI with a func-to-func state machine as described in the talk
	// "Lexical Scanning in Go" by Rob Pike (https://youtu.be/HxaD_trXwRE). each
	// state func returns the next state func we are transitioning to.
	for state := waitingForUCI(p); state != nil; {
		state = state(p)
	}
}

type statefn func(p *parser) statefn

func waitingForUCI(p *parser) statefn {
	_ = p.scanner.Scan()
	if err := p.scanner.Err(); err != nil {
		return errorScanning(p, err)
	}
	text := p.scanner.Text()
	switch text {
	case gteUCI:
		return uci
	case gteQuit:
		return nil // no further states
	default:
		p.logger.Printf("unrecognized: %s\n", text)
		return waitingForUCI
	}
}

func uci(p *parser) statefn {
	name, author, rest := p.handler.Identify()

	// print required name and author
	fmt.Println(etgID, etgIDName, name)
	fmt.Println(etgID, etgIDAuthor, author)

	// print rest of our id information in sorted order
	keys := make([]string, 0, len(rest))
	for k := range rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(etgID, k, rest[k])
	}

	fmt.Println(etgUCIOK)
	return waitingForCommand
}

func waitingForCommand(p *parser) statefn {
	_ = p.scanner.Scan()
	if err := p.scanner.Err(); err != nil {
		return errorScanning(p, err)
	}
	text := p.scanner.Text()
	switch text {
	case gteIsReady:
		p.handler.IsReady()
		fmt.Println(etgReadyOK)
		return waitingForCommand
	case gteNewGame:
		p.handler.NewGame()
		return waitingForCommand
	case gtePosition:
		return positionCommand
	case gteGo:
		return goCommand
	case gteQuit:
		return nil // no further states
	default:
		return errorUnrecognized(p, text, waitingForCommand)
	}
}

func positionCommand(p *parser) statefn {
	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	fen := p.scanner.Text()
	if fen == "startpos" {
		p.handler.SetStartingPosition()
	}
	// FIXME: fen isn't a single word...
	// b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	return waitingForCommand
}

func goCommand(p *parser) statefn {
	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	text := p.scanner.Text()
	switch text {
	case gteGoDepth:
		return goDepthCommand
	case gteGoInfinite:
		p.handler.GoInfinite()
	case gteGoNodes:
		return goNodesCommand
	}
	return waitingForCommand
}

func goDepthCommand(p *parser) statefn {
	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	plies, err := strconv.ParseUint(p.scanner.Text(), 10, 8)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoDepth(uint8(plies))
	fmt.Println(etgBestMove, movestr)
	return waitingForCommand
}

func goNodesCommand(p *parser) statefn {
	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	nodes, err := strconv.ParseUint(p.scanner.Text(), 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoNodes(nodes)
	fmt.Println(etgBestMove, movestr)
	return waitingForCommand
}

func errorUnrecognized(p *parser, text string, next statefn) statefn {
	p.logger.Printf("unrecognized: %s\n", text)
	return next
}

func errorParsingNumber(p *parser, err error, next statefn) statefn {
	p.logger.Printf("error parsing number: %v", err)
	return next
}

func errorScanning(p *parser, err error) statefn {
	p.logger.Printf("error scanning input: %v", err)
	return nil // no further states
}
