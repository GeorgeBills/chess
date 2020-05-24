package uci

import (
	"fmt"
	"io"
	"sync"

	chess "github.com/GeorgeBills/chess/m/v2"
)

type Wrapper struct {
	parser    *Parser
	executor  *Executor
	responder *Responder
}

func New(adapter Adapter, in io.Reader, out, log io.Writer) *Wrapper {
	parser, commandch, stopch := NewParser(in, log)
	executor, responsech := NewExecutor(commandch, stopch, adapter, log)
	responder := NewResponder(responsech, out, log)
	return &Wrapper{
		parser,
		executor,
		responder,
	}
}

func (w *Wrapper) ParseExecuteRespond() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		w.parser.Parse()
		wg.Done()
	}()
	go func() {
		w.executor.ExecuteCommands()
		wg.Done()
	}()
	go func() {
		w.responder.WriteResponses()
		wg.Done()
	}()
	return wg
}

// Move is a 3-tuple representing the "from" and "to" squares of a move, as well
// as which piece the moving piece will promote to (if any). This is the
// absolute minimum information required to unambiguously represent a chess
// move.
type Move struct {
	from, to  rankFile
	promoteTo chess.PromoteTo
}

// rankFile is a tuple representing the rank and file of a square. Both rank and
// file are zero indexed.
type rankFile struct {
	Rank, File uint8
}

// newMove returns a new chess move.
func newMove(from, to rankFile, promoteTo chess.PromoteTo) Move {
	return Move{from, to, promoteTo}
}

// From returns the square index the move is coming from.
func (m Move) From() uint8 { return chess.SquareIndex(m.from.Rank, m.from.File) }

// To returns the square index the move is going to.
func (m Move) To() uint8 { return chess.SquareIndex(m.to.Rank, m.to.File) }

// PromoteTo returns the piece the move will promote to, or PromoteToNone.
func (m Move) PromoteTo() chess.PromoteTo { return m.promoteTo }

// ParseUCIN parses a string in Universal Chess Notation (e.g. "a1h8") as a
// FromToPromote 3-tuple.
func ParseUCIN(ucin string) (Move, error) {
	if len(ucin) != 4 && len(ucin) != 5 {
		return Move{}, fmt.Errorf("invalid length for UCIN: %d", len(ucin))
	}
	fromRank, fromFile, err := chess.ParseAlgebraicNotationString(ucin[0:2])
	if err != nil {
		return Move{}, err
	}
	toRank, toFile, err := chess.ParseAlgebraicNotationString(ucin[2:4])
	if err != nil {
		return Move{}, err
	}

	promoteTo := chess.PromoteToNone
	if len(ucin) == 5 {
		switch ucin[4] {
		case 'q':
			promoteTo = chess.PromoteToQueen
		case 'n':
			promoteTo = chess.PromoteToKnight
		case 'r':
			promoteTo = chess.PromoteToRook
		case 'b':
			promoteTo = chess.PromoteToBishop
		default:
			return Move{}, fmt.Errorf("invalid promote to char: %q", ucin[4])
		}
	}

	from := rankFile{fromRank, fromFile}
	to := rankFile{toRank, toFile}
	return newMove(from, to, promoteTo), nil
}

// ToUCIN returns the move in Universal Chess Interface Notation (e.g. "a7a8q").
// UCIN is very similar to, but not exactly the same as, Long Algebraic
// Notation.
func ToUCIN(move chess.FromToPromoter) string {
	from, to := move.From(), move.To()

	ucin := [5]byte{
		'a' + chess.FileIndex(from), // a...h
		'1' + chess.RankIndex(from), // 1...8
		'a' + chess.FileIndex(to),   // a...h
		'1' + chess.RankIndex(to),   // 1...8
	}

	promoteTo := move.PromoteTo()
	if promoteTo != chess.PromoteToNone {
		switch promoteTo {
		case chess.PromoteToQueen:
			ucin[4] = 'q'
		case chess.PromoteToKnight:
			ucin[4] = 'n'
		case chess.PromoteToRook:
			ucin[4] = 'r'
		case chess.PromoteToBishop:
			ucin[4] = 'b'
		default:
			panic(fmt.Errorf("unrecognized promote to: %b", promoteTo))
		}
		return string(ucin[0:5])
	}

	return string(ucin[0:4])
}
