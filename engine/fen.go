package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// FEN returns the Forsyth–Edwards Notation for the board as a string.
func (b Board) FEN() string {
	sb := &strings.Builder{}
	// strings.Builder Write() methods always return a nil error, so this can never error
	b.WriteFEN(sb)
	return sb.String()
}

// WriteFEN writes the Forsyth–Edwards Notation for the board to w.
func (b Board) WriteFEN(w io.Writer) error {
	// bufio.Writer helps us defer error handling till the final Flush()
	// https://blog.golang.org/errors-are-values
	sb := bufio.NewWriter(w)

	var empty int
	var i uint8

	for i = 0; i < 64; i++ {
		poi := PrintOrderedIndex(i)
		p := b.PieceAt(poi)

		// sequences of empty squares are indicated with their count
		if (i%8 == 0 || p != PieceNone) && empty > 0 {
			sb.WriteString(strconv.Itoa(empty))
			empty = 0
		}

		// ranks are separated by /
		if i != 0 && i%8 == 0 {
			sb.WriteRune('/')
		}

		// pieces are indicated with a letter
		// uppercase for white, lowercase for black
		switch p {
		case PieceWhitePawn:
			sb.WriteRune('P')
		case PieceWhiteKnight:
			sb.WriteRune('N')
		case PieceWhiteBishop:
			sb.WriteRune('B')
		case PieceWhiteRook:
			sb.WriteRune('R')
		case PieceWhiteQueen:
			sb.WriteRune('Q')
		case PieceWhiteKing:
			sb.WriteRune('K')
		case PieceBlackPawn:
			sb.WriteRune('p')
		case PieceBlackKnight:
			sb.WriteRune('n')
		case PieceBlackBishop:
			sb.WriteRune('b')
		case PieceBlackRook:
			sb.WriteRune('r')
		case PieceBlackQueen:
			sb.WriteRune('q')
		case PieceBlackKing:
			sb.WriteRune('k')
		case PieceNone:
			empty++
		default:
			panic(fmt.Sprintf("invalid piece %b at index %d while generating FEN; %#v", p, i, b))
		}
	}
	// flush any remaining empty squares
	if empty > 0 {
		sb.WriteString(strconv.Itoa(empty))
	}

	sb.WriteRune(' ')

	if b.ToMove() == White {
		sb.WriteRune('w')
	} else {
		sb.WriteRune('b')
	}

	sb.WriteRune(' ')

	if b.meta&(maskWhiteCastleKingside|maskWhiteCastleQueenside|maskBlackCastleKingside|maskBlackCastleQueenside) == 0 {
		sb.WriteRune('-')
	} else {
		if b.CanWhiteCastleKingside() {
			sb.WriteRune('K')
		}
		if b.CanWhiteCastleQueenside() {
			sb.WriteRune('Q')
		}
		if b.CanBlackCastleKingside() {
			sb.WriteRune('k')
		}
		if b.CanBlackCastleQueenside() {
			sb.WriteRune('q')
		}
	}

	sb.WriteRune(' ')

	// TODO: output correct en passant square
	sb.WriteRune('-')

	sb.WriteRune(' ')

	// number of half moves since pawn movement or piece capture
	sb.WriteString(strconv.Itoa(b.HalfMoves()))

	sb.WriteRune(' ')

	// number of full moves
	sb.WriteString(strconv.Itoa(b.FullMoves()))

	return sb.Flush()
}

// NewBoardFromFEN returns a new board initialised as per the provided
// Forsyth–Edwards Notation. Only 8×8 boards are supported. Only basic
// validation of resulting board state is performed.
func NewBoardFromFEN(fen io.Reader) (*Board, error) {
	b := &Board{}
	r := bufio.NewReader(fen)

	unexpectingEOF := func(err error) error {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
	}

	skipspace := func() error {
		seen := false
		for {
			ch, err := r.ReadByte()
			if err != nil {
				return err
			}
			if ch != ' ' {
				if !seen {
					// require at least one space
					return fmt.Errorf("unexpected '%c', expecting ' '", ch)
				}
				return r.UnreadByte()
			}
			seen = true
		}
	}

	readuint8 := func() (uint8, error) {
		var n uint8
		seen := false
		for {
			ch, err := r.ReadByte()
			if err == io.EOF && seen {
				// expect at least one digit (which might be 0) before EOF
				return n, err
			}
			if err != nil {
				return 0, unexpectingEOF(err)
			}
			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				n = 10*n + uint8(ch-'0') // FIXME: check for overflow
				seen = true
			case ' ':
				err = r.UnreadByte()
				return n, err
			default:
				return 0, fmt.Errorf("unexpected '%c', expecting [0-9]", ch)
			}
		}
	}

	// read the 64 squares first
	var i uint8
READ_SQUARES:
	for i = 0; i < 64; {
		ch, err := r.ReadByte()
		if err != nil {
			return nil, unexpectingEOF(err)
		}

		if i != 0 && i%8 == 0 {
			// expect a / to indicate a new rank, which we immediately skip over
			if ch != '/' {
				return nil, fmt.Errorf("unexpected '%c', expecting '/'", ch)
			}
			ch, err = r.ReadByte()
			if err != nil {
				return nil, unexpectingEOF(err)
			}
		}

		var mask uint64 = 1 << PrintOrderedIndex(i)

		switch ch {
		case 'P':
			b.pawns |= mask
			b.white |= mask
		case 'N':
			b.knights |= mask
			b.white |= mask
		case 'B':
			b.bishops |= mask
			b.white |= mask
		case 'R':
			b.rooks |= mask
			b.white |= mask
		case 'Q':
			b.queens |= mask
			b.white |= mask
		case 'K':
			b.kings |= mask
			b.white |= mask
		case 'p':
			b.pawns |= mask
			b.black |= mask
		case 'n':
			b.knights |= mask
			b.black |= mask
		case 'b':
			b.bishops |= mask
			b.black |= mask
		case 'r':
			b.rooks |= mask
			b.black |= mask
		case 'q':
			b.queens |= mask
			b.black |= mask
		case 'k':
			b.kings |= mask
			b.black |= mask
		case '1', '2', '3', '4', '5', '6', '7', '8':
			i += uint8(ch - '0') // skip empty squares
			continue READ_SQUARES
		default:
			return nil, fmt.Errorf("unexpected '%c', expecting [PNBRQKpnbrqk1-8]", ch)
		}

		i++
	}

	var err error

	if err = skipspace(); err != nil {
		return nil, unexpectingEOF(err)
	}

	// read whose move it is (white or black)
	tomove, err := r.ReadByte()
	if err != nil {
		return nil, unexpectingEOF(err)
	}
	switch tomove {
	case 'w', 'b': // valid; unused till later
	default:
		return nil, fmt.Errorf("unexpected '%c', expecting [wb]", tomove)
	}

	if err = skipspace(); err != nil {
		return nil, unexpectingEOF(err)
	}

	// read castling information
READ_CASTLING:
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return nil, unexpectingEOF(err)
		}
		switch ch {
		case 'K':
			b.meta |= maskWhiteCastleKingside
		case 'Q':
			b.meta |= maskWhiteCastleQueenside
		case 'k':
			b.meta |= maskBlackCastleKingside
		case 'q':
			b.meta |= maskBlackCastleQueenside
		case '-':
			// '-' indicates that castling is unavailable
			// if present it must be the one and only byte
			if b.meta&(maskWhiteCastleKingside|maskWhiteCastleQueenside|maskBlackCastleKingside|maskBlackCastleQueenside) != 0 {
				return nil, errors.New("castling '-' must be solitary if present")
			}
			break READ_CASTLING
		case ' ':
			err := r.UnreadByte()
			if err != nil {
				return nil, unexpectingEOF(err)
			}
			break READ_CASTLING
		default:
			return nil, fmt.Errorf("unexpected '%c', expecting [KQkq]", ch)
		}
	}

	if err = skipspace(); err != nil {
		return nil, unexpectingEOF(err)
	}

	// read en passant square
	ch, err := r.ReadByte()
	if err != nil {
		return nil, unexpectingEOF(err)
	}
	if ch != '-' {
		// should be a file; store the file as the last 4 bits in the board meta
		switch ch {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			b.meta |= uint8(ch-'a') + 1
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
			b.meta |= uint8(ch-'A') + 1
		default:
			return nil, fmt.Errorf("unexpected '%c', expecting [a-hA-H]", ch)
		}

		// The next char should indicate either rank 3 or rank 6.
		//
		// We don't encode this rank into the board state, since it can be
		// inferred from which colour is to move next.
		//
		// If black is to move and an en passant is possible, then white must
		// have moved a pawn two spaces forward from its starting rank 2 on the
		// last move, and so the rank of the square under threat of en passant
		// must be 3. Similarly if white is to move then the rank can be
		// inferred to be 7.
		//
		// This means we need to validate the state here, otherwise the board is
		// inconsistent.
		ch, err := r.ReadByte()
		if err != nil {
			return nil, unexpectingEOF(err)
		}

		switch ch {
		case '3':
			if tomove == 'w' {
				return nil, fmt.Errorf("invalid board state: black moved last; en passant on rank 3")
			}
		case '6':
			if tomove == 'b' {
				return nil, fmt.Errorf("invalid board state: white moved last; en passant on rank 6")
			}
		default:
			return nil, fmt.Errorf("unexpected '%c', expecting [36]", ch)
		}
	}

	if err = skipspace(); err != nil {
		return nil, unexpectingEOF(err)
	}

	// read number of half moves
	if b.half, err = readuint8(); err != nil {
		return nil, unexpectingEOF(err)
	}

	if err = skipspace(); err != nil {
		return nil, unexpectingEOF(err)
	}

	// read number of full moves
	full, err := readuint8()
	if err != nil && err != io.EOF {
		// we're expecting an io.EOF err here
		return nil, err
	}
	b.total = uint16(2 * (full - 1))
	if tomove == 'b' {
		b.total++
	}

	err = b.Validate()
	if err != nil {
		return nil, err
	}

	return b, nil
}
