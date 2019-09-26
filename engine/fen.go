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
	var sb strings.Builder
	var empty int
	var i uint8

	for i = 0; i < 64; i++ {
		idx := i + 56 - 16*(i/8)
		p := b.PieceAt(idx)

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
			panic(fmt.Sprintf("invalid piece: %b", p))
		}
	}

	sb.WriteRune(' ')

	if b.ToMove() == White {
		sb.WriteRune('w')
	} else {
		sb.WriteRune('b')
	}

	sb.WriteRune(' ')

	if b.CanWhiteCastleKingSide() {
		sb.WriteRune('K')
	}

	if b.CanWhiteCastleQueenSide() {
		sb.WriteRune('Q')
	}

	if b.CanBlackCastleKingSide() {
		sb.WriteRune('k')
	}

	if b.CanBlackCastleQueenSide() {
		sb.WriteRune('q')
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

	return sb.String()
}

// NewBoardFromFEN returns a new board initialised as per the provided
// Forsyth–Edwards Notation. Only 8×8 boards are supported. No validation of
// resulting board state is performed; e.g. if the FEN specifies multiple kings
// per side then that's what the board will contain.
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
				n = 10*n + uint8(ch-'0')
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
READ_SQUARES:
	for i := 0; i < 64; {
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

		idx := i + 56 - 16*(i/8)

		switch ch {
		case 'P':
			b.pawns |= 1 << idx
			b.white |= 1 << idx
		case 'N':
			b.knights |= 1 << idx
			b.white |= 1 << idx
		case 'B':
			b.bishops |= 1 << idx
			b.white |= 1 << idx
		case 'R':
			b.rooks |= 1 << idx
			b.white |= 1 << idx
		case 'Q':
			b.queens |= 1 << idx
			b.white |= 1 << idx
		case 'K':
			b.kings |= 1 << idx
			b.white |= 1 << idx
		case 'p':
			b.pawns |= 1 << idx
			b.black |= 1 << idx
		case 'n':
			b.knights |= 1 << idx
			b.black |= 1 << idx
		case 'b':
			b.bishops |= 1 << idx
			b.black |= 1 << idx
		case 'r':
			b.rooks |= 1 << idx
			b.black |= 1 << idx
		case 'q':
			b.queens |= 1 << idx
			b.black |= 1 << idx
		case 'k':
			b.kings |= 1 << idx
			b.black |= 1 << idx
		case '1', '2', '3', '4', '5', '6', '7', '8':
			i += int(ch - '0') // skip empty squares
			continue READ_SQUARES
		default:
			return nil, fmt.Errorf("unexpected '%c', expecting [PNBRQKpnbrqk1-8]", ch)
		}

		i++
	}

	var err error

	err = skipspace()
	if err != nil {
		return nil, unexpectingEOF(err)
	}

	// read whose move it is (white or black)
	tomove, err := r.ReadByte()
	if err != nil {
		return nil, unexpectingEOF(err)
	}
	switch tomove {
	case 'w':
		b.meta |= White
	case 'b':
		b.meta |= Black
	default:
		return nil, fmt.Errorf("unexpected '%c', expecting [wb]", tomove)
	}

	err = skipspace()
	if err != nil {
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
			b.meta |= wcks
		case 'Q':
			b.meta |= wcqs
		case 'k':
			b.meta |= bcks
		case 'q':
			b.meta |= bcqs
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

	err = skipspace()
	if err != nil {
		return nil, unexpectingEOF(err)
	}

	// read en passant square
	ch, err := r.ReadByte()
	if err != nil {
		return nil, unexpectingEOF(err)
	}
	if ch != '-' {
		// TODO: support reading en passant square
		return nil, errors.New("en passant unsupported")
	}

	err = skipspace()
	if err != nil {
		return nil, unexpectingEOF(err)
	}

	// read number of half moves
	b.half, err = readuint8()
	if err != nil {
		return nil, unexpectingEOF(err)
	}

	err = skipspace()
	if err != nil {
		return nil, unexpectingEOF(err)
	}

	// read number of full moves
	b.full, err = readuint8()
	if err != nil && err != io.EOF {
		// we're expecting an io.EOF err here
		return nil, err
	}

	return b, nil
}
