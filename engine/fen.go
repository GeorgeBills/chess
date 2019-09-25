package engine

import (
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
