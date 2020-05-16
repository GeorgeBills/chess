package uci

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
)

// GUI-to-engine constants are tokens sent from the GUI to the engine.
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

// Engine-to-GUI constants are tokens sent from the engine to the GUI.
const (
	etgID       = "id"       // sent to identify the engine
	etgIDName   = "name"     // e.g. "id name Shredder X.Y\n"
	etgIDAuthor = "author"   // e.g. "id author Stefan MK\n"
	etgUCIOK    = "uciok"    // the engine has sent all infos and is ready
	etgReadyOK  = "readyok"  // the engine is ready to accept new commands
	etgBestMove = "bestmove" // engine has stopped searching and found the best move
	etgInfo     = "info"     // engine wants to send information to the GUI
)

// Handler handles events generated from parsing UCI.
type Handler interface {
	Identify() (name, author string, other map[string]string)
	IsReady()
	NewGame()
	SetStartingPosition()
	SetPosition(fen string)
	GoDepth(plies uint8) string
	GoNodes(nodes uint64) string
	GoInfinite()
}

// NewParser returns a new parser.
func NewParser(h Handler, r io.Reader, outw, logw io.Writer) *Parser {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	return &Parser{
		handler: h,
		scanner: scanner,
		logger:  log.New(logw, "parser: ", log.LstdFlags),
		out:     bufio.NewWriter(outw),
	}
}

// Parser parses and generates events from UCI.
type Parser struct {
	logger  *log.Logger
	handler Handler
	scanner *bufio.Scanner
	out     *bufio.Writer
}

// Run starts the parser.
func (p *Parser) Run() error {
	// We parse UCI with a func-to-func state machine as described in the talk
	// "Lexical Scanning in Go" by Rob Pike (https://youtu.be/HxaD_trXwRE). Each
	// state func returns the next state func we are transitioning to. We start
	// off in the "waiting for UCI" state. We perpetually evaluate states until
	// we receive a nil (terminal) state.
	for state := waitingForUCI(p); state != nil; {
		state = state(p)
		err := p.out.Flush() // flush output for each state
		if err != nil {
			return err
		}
	}
	// TODO: give the engine a chance to cleanup here
	p.logger.Println("finished")
	return nil
}

type statefn func(p *Parser) statefn

func waitingForUCI(p *Parser) statefn {
	p.logger.Println("waiting for uci")

	_ = p.scanner.Scan()
	if err := p.scanner.Err(); err != nil {
		return errorScanning(p, err)
	}
	text := p.scanner.Text()
	switch text {
	case gteUCI:
		return commandUCI
	case gteQuit:
		return commandQuit
	default:
		return errorUnrecognized(p, text, waitingForUCI)
	}
}

func waitingForCommand(p *Parser) statefn {
	p.logger.Println("waiting for command")

	_ = p.scanner.Scan()
	if err := p.scanner.Err(); err != nil {
		return errorScanning(p, err)
	}
	text := p.scanner.Text()
	switch text {
	case gteIsReady:
		p.handler.IsReady()
		fmt.Fprintln(p.out, etgReadyOK)
		return waitingForCommand
	case gteNewGame:
		p.handler.NewGame()
		return waitingForCommand
	case gtePosition:
		return commandPosition
	case gteGo:
		return commandGo
	case gteQuit:
		return commandQuit
	default:
		return errorUnrecognized(p, text, waitingForCommand)
	}
}

func commandUCI(p *Parser) statefn {
	p.logger.Println("uci")

	name, author, rest := p.handler.Identify()

	// print required name and author
	fmt.Fprintln(p.out, etgID, etgIDName, name)
	fmt.Fprintln(p.out, etgID, etgIDAuthor, author)

	// print rest of our id information in sorted order
	keys := make([]string, 0, len(rest))
	for k := range rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintln(p.out, etgID, k, rest[k])
	}

	fmt.Fprintln(p.out, etgUCIOK)
	return waitingForCommand
}

func commandPosition(p *Parser) statefn {
	p.logger.Println("command: position")

	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	fen := p.scanner.Text()
	if fen == "startpos" {
		p.handler.SetStartingPosition()
	}
	// FIXME: fen isn't a single word...
	// regexp.FindReaderSubmatchIndex?
	// b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	return waitingForCommand
}

func commandGo(p *Parser) statefn {
	p.logger.Println("command: go")

	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	text := p.scanner.Text()
	switch text {
	case gteGoDepth:
		return commandGoDepth
	case gteGoInfinite:
		p.handler.GoInfinite()
	case gteGoNodes:
		return commandGoNodes
	}
	return waitingForCommand
}

func commandGoDepth(p *Parser) statefn {
	p.logger.Println("command: go depth")

	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	plies, err := strconv.ParseUint(p.scanner.Text(), 10, 8)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoDepth(uint8(plies))
	fmt.Fprintln(p.out, etgBestMove, movestr)
	return waitingForCommand
}

func commandGoNodes(p *Parser) statefn {
	p.logger.Println("command: go nodes")

	_ = p.scanner.Scan()
	_ = p.scanner.Err()
	nodes, err := strconv.ParseUint(p.scanner.Text(), 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoNodes(nodes)
	fmt.Fprintln(p.out, etgBestMove, movestr)
	return waitingForCommand
}

func commandQuit(p *Parser) statefn {
	p.logger.Println("quitting")
	return nil
}

func errorUnrecognized(p *Parser, text string, next statefn) statefn {
	p.logger.Printf("unrecognized: %s\n", text)
	return next
}

func errorParsingNumber(p *Parser, err error, next statefn) statefn {
	p.logger.Printf("error parsing number: %v", err)
	return next
}

func errorScanning(p *Parser, err error) statefn {
	p.logger.Printf("error scanning input: %v", err)
	return nil // no further states
}
