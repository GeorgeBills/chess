package uci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
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
	SetPositionFEN(fen string)
	GoDepth(plies uint8) string
	GoNodes(nodes uint64) string
	GoInfinite()
	Quit()
	// TODO: most handler methods should return error
}

// NewParser returns a new parser.
func NewParser(h Handler, r io.Reader, outw, logw io.Writer) *Parser {
	return &Parser{
		handler: h,
		reader:  bufio.NewReader(r),
		logger:  log.New(logw, "parser: ", log.LstdFlags),
		out:     bufio.NewWriter(outw),
	}
}

// Parser parses and generates events from UCI.
type Parser struct {
	// TODO: rename to StateMachine or similar, this isn't just a parser.
	logger  *log.Logger
	handler Handler
	reader  io.RuneScanner
	out     *bufio.Writer
}

// Run starts the parser.
func (p *Parser) Run() error {
	// We parse UCI with a func-to-func state machine as described in the talk
	// "Lexical Scanning in Go" by Rob Pike (https://youtu.be/HxaD_trXwRE). Each
	// state func returns the next state func we are transitioning to. We start
	// off in the "waiting for UCI" state. We perpetually evaluate states until
	// we receive a nil (terminal) state.
	//
	// We manually read from the buffer in lieu of bufio.Scanner for a few
	// reasons: one, newlines are meaningful in UCI as unambiguous command
	// terminators, and two because of the one annoying case (FEN) where a token
	// contains whitespace.
	for state := waitingForUCI; state != nil; {
		state = state(p)
		// TODO: flush as needed
		//       or just manually buffer so the default is synchronous?
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

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteUCI:
		return commandUCI
	case gteQuit:
		return commandQuit
	case "": // newline
		return eol(p, waitingForUCI)
	default:
		return errorUnrecognized(p, token, waitingForUCI)
	}
}

func waitingForCommand(p *Parser) statefn {
	p.logger.Println("waiting for command")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteIsReady:
		return commandIsReady
	case gteNewGame:
		return commandNewGame
	case gtePosition:
		return commandPosition
	case gteGo:
		return commandGo
	case gteQuit:
		return commandQuit
	case "":
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, waitingForCommand)
	}
}

func commandUCI(p *Parser) statefn {
	p.logger.Println("command: uci")

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

func commandIsReady(p *Parser) statefn {
	p.logger.Println("command: isready")

	p.handler.IsReady()
	fmt.Fprintln(p.out, etgReadyOK)
	return waitingForCommand
}

func commandNewGame(p *Parser) statefn {
	p.logger.Println("command: new game")

	p.handler.NewGame()
	return waitingForCommand
}

func commandPosition(p *Parser) statefn {
	p.logger.Println("command: position")

	err := consume(p.reader, isSpace)
	if err != nil {
		return errorScanning(p, err)
	}

	token, err := readToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case "startpos":
		return commandPositionStarting
	case "fen":
		return commandPositionFEN
	case "": // newline
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandPosition)
	}
}

func commandPositionStarting(p *Parser) statefn {
	p.logger.Println("command: position startpos")

	p.handler.SetStartingPosition()
	return waitingForCommand
}

func commandPositionFEN(p *Parser) statefn {
	p.logger.Println("command: position fen")

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}

	var buf bytes.Buffer

	// ranks: [/1-8BKNPQRbknpqr]
	if err := accept(p.reader, &buf, "/12345678BKNPQRbknpqr"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN ranks: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	// to play: [wb]
	if err := accept(p.reader, &buf, "wb"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN to play: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	// castling: [KQkq]
	if err := accept(p.reader, &buf, "KQkq"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN castling: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	// en passant: [-A-Ha-h][1-8]
	if err := accept(p.reader, &buf, "-ABCDEFGHabcdefgh12345678"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN en passant: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	// half moves
	if err := accept(p.reader, &buf, "1234567890"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN half moves: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	// full moves
	if err := accept(p.reader, &buf, "1234567890"); err != nil {
		return errorScanning(p, fmt.Errorf("error scanning FEN full moves: %w", err))
	}

	if err := consume(p.reader, isSpace); err != nil {
		return errorScanning(p, err)
	}
	buf.WriteRune(' ')

	p.handler.SetPositionFEN(buf.String())
	return waitingForCommand
}

func commandGo(p *Parser) statefn {
	p.logger.Println("command: go")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteGoDepth:
		return commandGoDepth
	case gteGoInfinite:
		return commandGoInfinite
	case gteGoNodes:
		return commandGoNodes
	}
	return waitingForCommand
}

func commandGoDepth(p *Parser) statefn {
	p.logger.Println("command: go depth")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	plies, err := strconv.ParseUint(token, 10, 8)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoDepth(uint8(plies))
	fmt.Fprintln(p.out, etgBestMove, movestr)
	return waitingForCommand
}

func commandGoInfinite(p *Parser) statefn {
	p.logger.Println("command: go infinite")

	p.handler.GoInfinite()
	return waitingForCommand
}

func commandGoNodes(p *Parser) statefn {
	p.logger.Println("command: go nodes")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	nodes, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}
	movestr := p.handler.GoNodes(nodes)
	fmt.Fprintln(p.out, etgBestMove, movestr)
	return waitingForCommand
}

func commandQuit(p *Parser) statefn {
	p.logger.Println("quitting")
	p.handler.Quit()
	return nil
}

// TODO: expectEOL: read spaces to EOL
//                  error if there are any other tokens
//                  transition to the specified next state

func eol(p *Parser, next statefn) statefn {
	// TODO: add "expected bool" param; if unexpected log a warning
	err := consume(p.reader, isEOL) // consume all newline runes
	if err != nil {
		return errorScanning(p, err)
	}
	return next
}

// errorUnrecognized logs an "unrecognized token" error and then transitions to
// the specified next state. In general, per the UCI specification which states
// that unrecognized tokens should be ignored, the next state should be the
// state we just transitioned out of (i.e. this should loop back to the
// originating state).
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

// isSpace returns whether or not c represents a space. We do not consider end
// of line characters to be spaces, as newlines are meaningful tokens within
// UCI, representing the end of a command.
func isSpace(c rune) bool {
	return c == ' ' || c == '\t'
}

// isEOL returns whether or not c represents part of a newline.
func isEOL(c rune) bool {
	return c == '\r' || c == '\n'
}

// consume reads characters until the predicate returns false.
func consume(r io.RuneScanner, pred func(rune) bool) error {
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !pred(c) {
			r.UnreadRune()
			return nil
		}
	}
}

func accept(r io.RuneReader, buf *bytes.Buffer, valid string) error {
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		switch {
		case strings.IndexRune(valid, c) != -1:
			buf.WriteRune(c)
		case isSpace(c) || isEOL(c):
			return nil
		default:
			return fmt.Errorf("unrecognized rune: %q", c)
		}
	}
}

func nextToken(r io.RuneScanner) (string, error) {
	consume(r, isSpace)
	return readToken(r)
	// TODO: manually inline readToken in here, it's never called otherwise
	// TODO: special case any EOL to consume the whole thing and return "\n"
}

// readToken reads runes from r until it finds a rune that terminates the token,
// returning a buffer of runes that were read.
func readToken(r io.RuneScanner) (string, error) {
	var buf bytes.Buffer
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if isSpace(c) || isEOL(c) {
			r.UnreadRune()
			return buf.String(), nil
		}
		buf.WriteRune(c)
	}
}
