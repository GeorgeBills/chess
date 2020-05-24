package uci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// GUI-to-engine constants are tokens sent from the GUI to the engine.
const (
	gteUCI        = "uci"        // tell engine to use the universal chess interface
	gteDebug      = "debug"      // switch the debug mode of the engine on and off
	gteIsReady    = "isready"    // used to synchronize the engine with the GUI
	gteSetOption  = "setoption"  // change internal parameters of the engine
	gteNewGame    = "ucinewgame" // the next search will be from a different game
	gtePosition   = "position"   // set up the position described on the internal board
	gteStartPos   = "startpos"   // game was played from the start position
	gteFEN        = "fen"        // position described in fenstring
	gteMoves      = "moves"      // play the moves on the internal chess board
	gteGo         = "go"         // start calculating on the current position
	gteGoDepth    = "depth"      // search x plies only
	gteGoNodes    = "nodes"      // search x nodes only
	gteGoMoveTime = "movetime"   // search exactly x mseconds
	gteGoInfinite = "infinite"   // search until the stop command
	gteWhiteTime  = "wtime"      // white has x msec left on the clock
	gteBlackTime  = "btime"      // black has x msec left on the clock
	gteWhiteInc   = "winc"       // white increment per move in mseconds
	gteBlackInc   = "binc"       //	black increment per move in mseconds
	gteStop       = "stop"       // stop calculating as soon as possible
	gteQuit       = "quit"       // quit the program as soon as possible
)

// NewParser returns a new parser.
func NewParser(r io.Reader, logw io.Writer) (*Parser, <-chan Execer, <-chan struct{}) {
	commandch := make(chan Execer, 0) // unbuffered
	stopch := make(chan struct{}, 1)
	parser := &Parser{
		commandch: commandch,
		stopch:    stopch,
		reader:    bufio.NewReader(r),
		logger:    log.New(logw, "parser: ", log.LstdFlags),
	}
	return parser, commandch, stopch
}

// Parser parses and generates events from UCI.
type Parser struct {
	logger    *log.Logger
	commandch chan<- Execer
	stopch    chan<- struct{}
	reader    io.RuneScanner
}

// Parse starts the parser.
func (p *Parser) Parse() {
	defer close(p.stopch)
	defer close(p.commandch)

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

	p.logger.Println("starting")
	for state := waitingForUCI; state != nil; {
		state = state(p)
	}
	// TODO: give the engine a chance to cleanup here
	p.logger.Println("finished")
}

type statefn func(p *Parser) statefn

func waitingForUCI(p *Parser) statefn { // TODO: "init"
	p.logger.Println("waiting for uci")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteUCI:
		return p.emit(CommandUCI{}, waitingForCommand)
	case gteQuit:
		return commandQuit
	case "": // newline
		return eol(p, waitingForUCI)
	default:
		return errorUnrecognized(p, token, waitingForUCI)
	}
}

func waitingForCommand(p *Parser) statefn { // TODO: "running"
	p.logger.Println("waiting for command")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteIsReady:
		return p.emit(CommandIsReady{}, waitingForCommand)
	case gteNewGame:
		return p.emit(CommandNewGame{}, waitingForCommand)
	case gtePosition:
		return commandPosition
	case gteGo:
		return commandGo
	case gteQuit:
		return commandQuit
	case gteStop:
		return p.emitStop(waitingForCommand)
	case "":
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, waitingForCommand)
	}
}

func (p *Parser) emit(cmd Execer, next statefn) statefn {
	select {
	case p.commandch <- cmd:
		p.logger.Printf("emitted %T command", cmd)
		return next
	default:
		p.logger.Printf("cannot emit %T command; command already in progress or executor not running", cmd)
		return next
	}
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
	case gteStartPos:
		return p.emit(CommandSetStartingPosition{}, eol(p, waitingForCommand))
	case gteFEN:
		return commandPositionFEN
	case "": // newline
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandPosition)
	}
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

	return p.emit(CommandSetPositionFEN{FEN: buf.String()}, commandPositionMoves)
}

func commandPositionMoves(p *Parser) statefn {
	p.logger.Println("command: position moves")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	// could be moves, could be EOL
	switch token {
	case gteMoves:
		return commandPositionMovesMove
	case "": // newline
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandPositionMoves)
	}
}

func commandPositionMovesMove(p *Parser) statefn {
	p.logger.Println("command: position moves move")

	// loop over all moves
	for {
		token, err := nextToken(p.reader)
		if err != nil {
			return errorScanning(p, err)
		}

		if token == "" {
			break // newline
		}

		move, err := ParseUCIN(token)
		if err != nil {
			return errorUnrecognized(p, token, commandPositionMoves) // TODO: pass along err so we get decent logs
		}

		// block on apply move, since we'll likely pass them through faster than
		// executor can take them off the channel
		// TODO: or use the accumulator pattern to make it one big command?
		p.commandch <- CommandApplyMove{move}
	}

	return waitingForCommand
}

// TimeControl represents time controls for playing a chess move.
type TimeControl struct {
	WhiteTime, BlackTime           time.Duration
	WhiteIncrement, BlackIncrement time.Duration
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
		return p.emit(CommandGoInfinite{}, eol(p, waitingForCommand))
	case gteGoNodes:
		return commandGoNodes
	case gteBlackTime:
		return commandGoTimeBlackTime(p, TimeControl{})
	case gteWhiteTime:
		return commandGoTimeWhiteTime(p, TimeControl{})
	case gteBlackInc:
		return commandGoTimeBlackIncrement(p, TimeControl{})
	case gteWhiteInc:
		return commandGoTimeWhiteIncrement(p, TimeControl{})
	case "": // newline
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandGo)
	}
}

func commandGoTime(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	switch token {
	case gteBlackTime:
		return commandGoTimeBlackTime(p, accumulator)
	case gteWhiteTime:
		return commandGoTimeWhiteTime(p, accumulator)
	case gteBlackInc:
		return commandGoTimeBlackIncrement(p, accumulator)
	case gteWhiteInc:
		return commandGoTimeWhiteIncrement(p, accumulator)
	case "": // newline
		// command finished, so run it
		// TODO: check we have at least white and black time set
		p.commandch <- CommandGoTime{accumulator}
		return eol(p, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandGoTime(p, accumulator))
	}
}

func commandGoTimeWhiteTime(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time white time")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	t, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}

	accumulator.WhiteTime = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeBlackTime(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time black time")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	t, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}

	accumulator.BlackTime = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeWhiteIncrement(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time white increment")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	t, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}

	accumulator.WhiteIncrement = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeBlackIncrement(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time black increment")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	t, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return errorParsingNumber(p, err, waitingForCommand)
	}

	accumulator.BlackIncrement = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
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

	p.commandch <- CommandGoDepth{uint8(plies)}
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

	p.commandch <- CommandGoNodes{nodes}
	return waitingForCommand
}

func commandQuit(p *Parser) statefn {
	p.logger.Println("quitting")
	return nil
}

func (p *Parser) emitStop(next statefn) statefn {
	p.logger.Println("stopping current command")
	p.stopch <- struct{}{}
	return next
}

// TODO: expectEOL: read spaces to EOL
//                  error if there are any other tokens
//                  transition to the specified next state
//       most states should finish with expectEOL(expectCommand)

func eol(p *Parser, next statefn) statefn {
	// TODO: add "expected bool" param; if unexpected log a warning
	err := consume(p.reader, isEOL) // consume all newline runes
	if err != nil {
		return errorScanning(p, err)
	}
	return next
}

// TODO: define known types for various errors and pass to generic error handler
// e.g. some methods might return diff errors, some fatal, some not
// should use errors.Is() and similar to inspect the error

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

func accept(r io.RuneScanner, buf *bytes.Buffer, valid string) error {
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		switch {
		case strings.IndexRune(valid, c) != -1:
			buf.WriteRune(c)
		case isSpace(c) || isEOL(c):
			r.UnreadRune()
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
