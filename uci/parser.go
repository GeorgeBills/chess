package uci

import (
	"bufio"
	"bytes"
	"errors"
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
func NewParser(r io.Reader, logw io.Writer) (*Parser, <-chan Command, <-chan struct{}) {
	commandch := make(chan Command, 0) // unbuffered
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
	commandch chan<- Command
	stopch    chan<- struct{}
	reader    io.RuneScanner
}

// ParseInput starts the parser.
func (p *Parser) ParseInput() {
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

func (p *Parser) emit(cmd Command, next statefn) statefn {
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
		partial := &CommandSetStartingPosition{}
		return commandPositionMoves(p, partial)
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

	// loop through our six character classes, accepting a token made up of
	// those characters followed by one or more spaces for each; insert a single
	// space between each token.
	var buf bytes.Buffer
	valids := []string{
		"/12345678BKNPQRbknpqr",     // ranks
		"wb",                        // to play
		"KQkq",                      // castling
		"-ABCDEFGHabcdefgh12345678", // en passant
		"1234567890",                // half moves
		"1234567890",                // full moves
	}
	for i, valid := range valids {
		if err := accept(p.reader, &buf, valid); err != nil {
			return p.handleError(fmt.Errorf("error scanning FEN: %w", err), false, commandPositionFEN)
		}
		if err := consume(p.reader, isSpace); err != nil {
			return errorScanning(p, err)
		}
		if i != len(valids)-1 {
			buf.WriteRune(' ')
		}
	}

	partial := &CommandSetPositionFEN{FEN: buf.String()}
	return commandPositionMoves(p, partial)
}

func commandPositionMoves(p *Parser, partial AppendMoveCommand) statefn {
	p.logger.Println("command: position moves")

	token, err := nextToken(p.reader)
	if err != nil {
		return errorScanning(p, err)
	}

	// could be moves, could be EOL
	switch token {
	case gteMoves:
		return commandPositionMovesMove(p, partial)
	case "":
		// newline: fire off the command
		return p.emit(partial, waitingForCommand)
	default:
		return errorUnrecognized(p, token, commandPositionMoves(p, partial))
	}
}

func commandPositionMovesMove(p *Parser, partial AppendMoveCommand) statefn {
	p.logger.Println("command: position moves move")

	// loop over moves until we get an error or a newline
	for {
		token, err := nextToken(p.reader)
		if err != nil {
			return errorScanning(p, err)
		}

		if token == "" {
			// newline: fire off the command
			return p.emit(partial, waitingForCommand)
		}

		move, err := ParseUCIN(token)
		if err != nil {
			next := commandPositionMovesMove(p, partial)
			return errorUnrecognized(p, token, next) // TODO: pass along err so we get decent logs
		}

		partial.AppendMove(move)
	}
}

// TimeControl represents time controls for playing a chess move.
type TimeControl struct {
	WhiteTime, BlackTime           time.Duration
	WhiteIncrement, BlackIncrement time.Duration
}

func (tc TimeControl) String() string {
	return fmt.Sprintf(
		"wtime %d btime %d winc %d binc %d",
		tc.WhiteTime.Milliseconds(),
		tc.BlackTime.Milliseconds(),
		tc.WhiteIncrement.Milliseconds(),
		tc.BlackIncrement.Milliseconds(),
	)
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

	t, err := nextTokenUint(p.reader, 64)
	if err != nil {
		return p.handleError(err, false, commandGoTimeWhiteTime(p, accumulator))
	}

	accumulator.WhiteTime = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeBlackTime(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time black time")

	t, err := nextTokenUint(p.reader, 64)
	if err != nil {
		return p.handleError(err, false, commandGoTimeBlackTime(p, accumulator))
	}

	accumulator.BlackTime = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeWhiteIncrement(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time white increment")

	t, err := nextTokenUint(p.reader, 64)
	if err != nil {
		return p.handleError(err, false, commandGoTimeWhiteIncrement(p, accumulator))
	}

	accumulator.WhiteIncrement = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoTimeBlackIncrement(p *Parser, accumulator TimeControl) statefn {
	p.logger.Println("command: go time black increment")

	t, err := nextTokenUint(p.reader, 64)
	if err != nil {
		return p.handleError(err, false, commandGoTimeBlackIncrement(p, accumulator))
	}

	accumulator.BlackIncrement = time.Duration(t) * time.Millisecond
	return commandGoTime(p, accumulator)
}

func commandGoDepth(p *Parser) statefn {
	p.logger.Println("command: go depth")

	plies, err := nextTokenUint(p.reader, 8)
	if err != nil {
		return p.handleError(err, false, commandGoDepth)
	}

	p.commandch <- CommandGoDepth{uint8(plies)}
	return waitingForCommand
}

func commandGoNodes(p *Parser) statefn {
	p.logger.Println("command: go nodes")

	nodes, err := nextTokenUint(p.reader, 64)
	if err != nil {
		return p.handleError(err, false, commandGoNodes)
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

// handleError takes an error that we don't know the type of and returns the
// next state based on that error.
func (p *Parser) handleError(err error, eofOK bool, next statefn) statefn {
	var (
		ire *invalidRuneError
		ne  *strconv.NumError
	)
	switch {
	case errors.As(err, &ne): // e.g. strconv.ErrRange, strconv.ErrSyntax
		return errorParsingNumber(p, err, next)
	case errors.As(err, &ire):
		return errorInvalidRune(p, err, next)
	case errors.Is(err, io.EOF):
		if !eofOK {
			return errorScanning(p, io.ErrUnexpectedEOF)
		}
		return nil
	default:
		// assume it's fatal
		return errorScanning(p, err)
	}
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

func errorInvalidRune(p *Parser, err error, next statefn) statefn {
	p.logger.Println(err)
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

type invalidRuneError struct{ c rune }

func (err *invalidRuneError) Error() string {
	return fmt.Sprintf("invalid rune: %q", err.c)
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
			// TODO: advance to whitespace, this token is bad
			return &invalidRuneError{c}
		}
	}
}

func nextToken(r io.RuneScanner) (string, error) {
	consume(r, isSpace)
	return readToken(r)
	// TODO: manually inline readToken in here, it's never called otherwise
	// TODO: special case any EOL to consume the whole thing and return "\n"
}

func nextTokenUint(r io.RuneScanner, bits int) (uint64, error) {
	token, err := nextToken(r)
	if err != nil {
		return 0, err
	}

	n, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
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
