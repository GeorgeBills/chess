package engine

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strings"
)

// Colour is used to represent each colour.
type Colour byte

// White and Black are constants defined for the colours.
const (
	White Colour = 'w'
	Black Colour = 'b'
)

// Board represents an 8×8 chess board.
//
// The 0th index represents A1 and the 63rd index represents H8.
//
//          A  B  C  D  E  F  G  H
//     8 | 56 57 58 59 60 61 62 63
//     7 | 48 49 50 51 52 53 54 55
//     6 | 40 41 42 43 44 45 46 47
//     5 | 32 33 34 35 36 37 38 39
//     4 | 24 25 26 27 28 29 30 31
//     3 | 16 17 18 19 20 21 22 23
//     2 |  8  9 10 11 12 13 14 15
//     1 |  0  1  2  3  4  5  6  7
//
// For an index i into the board, i/8 is the rank and i%8 is the file.
type Board struct {
	black   uint64
	white   uint64
	pawns   uint64
	knights uint64
	bishops uint64
	rooks   uint64
	queens  uint64
	kings   uint64

	// half and total both record "half" moves. A "half" move is a move where
	// one colour has moved; a "full" move is a move where both colours have
	// moved. We only record half moves, and calculate full moves if needed.
	//
	// half records the number of half moves since a pawn was moved or a piece
	// was captured, and is used for determing if a draw can be claimed under
	// the fifty-move rule.
	//
	// total records the total number of half moves since the start of the game.
	// It starts at 0 and is incremented to 1 post white's first move, 2 post
	// black's first move, 3, 4... etc. It will always be even if it is white's
	// turn to move and odd if it is black's turn to move.
	half  uint8
	total uint16

	// meta records meta information about the board, specifically castling
	// rights and whether any pawn is vulnerable to en passant.
	meta byte
}

const (
	maskWhiteCastleKingside  = 0b10000000
	maskWhiteCastleQueenside = 0b01000000
	maskBlackCastleKingside  = 0b00100000
	maskBlackCastleQueenside = 0b00010000
	maskEnPassant            = 0b00001111 // the last 4 bits of meta indicate the file for a valid en passant
)

// NewBoard returns a board in the initial state.
func NewBoard() Board {
	return Board{
		white:   maskRank1 | maskRank2,
		black:   maskRank7 | maskRank8,
		pawns:   maskRank2 | maskRank7,
		knights: 1<<B1 | 1<<G1 | 1<<B8 | 1<<G8,
		bishops: 1<<C1 | 1<<F1 | 1<<C8 | 1<<F8,
		rooks:   1<<A1 | 1<<H1 | 1<<A8 | 1<<H8,
		queens:  1<<D1 | 1<<D8,
		kings:   1<<E1 | 1<<E8,
		half:    0,
		total:   0,
		meta:    maskWhiteCastleKingside | maskWhiteCastleQueenside | maskBlackCastleKingside | maskBlackCastleQueenside,
	}
}

// IsWhiteAt returns true iff there is a piece at index i and it is white.
func (b Board) IsWhiteAt(i uint8) bool { return b.white&(1<<i) != 0 }

// IsBlackAt returns true iff there is a piece at index i and it is black.
func (b Board) IsBlackAt(i uint8) bool { return b.black&(1<<i) != 0 }

// IsEmptyAt returns true iff the square at index i is empty.
func (b Board) IsEmptyAt(i uint8) bool { return !b.IsWhiteAt(i) && !b.IsBlackAt(i) }

// IsPawnAt returns true iff there is a piece at index i and it is a pawn.
func (b Board) IsPawnAt(i uint8) bool { return b.pawns&(1<<i) != 0 }

// IsKnightAt returns true iff there is a piece at index i and it is a knight.
func (b Board) IsKnightAt(i uint8) bool { return b.knights&(1<<i) != 0 }

// IsBishopAt returns true iff there is a piece at index i and it is a bishop.
func (b Board) IsBishopAt(i uint8) bool { return b.bishops&(1<<i) != 0 }

// IsRookAt returns true iff there is a piece at index i and it is a rook.
func (b Board) IsRookAt(i uint8) bool { return b.rooks&(1<<i) != 0 }

// IsQueenAt returns true iff there is a piece at index i and it is a queen.
func (b Board) IsQueenAt(i uint8) bool { return b.queens&(1<<i) != 0 }

// IsKingAt returns true iff there is a piece at index i and it is a king.
func (b Board) IsKingAt(i uint8) bool { return b.kings&(1<<i) != 0 }

// EnPassant returns the index of the square under threat of en passant, or
// math.MaxUint8 if there is no such square.
func (b Board) EnPassant() uint8 {
	file := uint8(b.meta & maskEnPassant)
	if file == 0 {
		return math.MaxUint8
	}
	tomove := b.ToMove()
	// index = 8×(rank - 1) + file - 1
	switch tomove {
	case White:
		return 39 + file // rank 6; 8×(6 - 1) + file - 1
	case Black:
		return 15 + file // rank 3; 8×(3 - 1) + file - 1
	default:
		panic(fmt.Sprintf("invalid to move: %b", tomove))
	}
}

// PieceAt returns the piece at index i.
func (b Board) PieceAt(i uint8) Piece {
	if b.IsWhiteAt(i) {
		switch {
		case b.IsPawnAt(i):
			return PieceWhitePawn
		case b.IsKnightAt(i):
			return PieceWhiteKnight
		case b.IsBishopAt(i):
			return PieceWhiteBishop
		case b.IsRookAt(i):
			return PieceWhiteRook
		case b.IsQueenAt(i):
			return PieceWhiteQueen
		case b.IsKingAt(i):
			return PieceWhiteKing
		default:
			panic(fmt.Sprintf("invalid white piece at index %d; %#v", i, b))
		}
	}
	if b.IsBlackAt(i) {
		switch {
		case b.IsPawnAt(i):
			return PieceBlackPawn
		case b.IsKnightAt(i):
			return PieceBlackKnight
		case b.IsBishopAt(i):
			return PieceBlackBishop
		case b.IsRookAt(i):
			return PieceBlackRook
		case b.IsQueenAt(i):
			return PieceBlackQueen
		case b.IsKingAt(i):
			return PieceBlackKing
		default:
			panic(fmt.Sprintf("invalid black piece at index %d; %#v", i, b))
		}
	}
	return PieceNone
}

// ToMove returns the colour whose move it is.
func (b Board) ToMove() Colour {
	if b.total%2 == 0 {
		return White
	}
	return Black
}

// CanWhiteCastleKingside returns true iff white can castle king side.
func (b Board) CanWhiteCastleKingside() bool { return b.meta&maskWhiteCastleKingside != 0 }

// CanWhiteCastleQueenside returns true iff white can castle queen side.
func (b Board) CanWhiteCastleQueenside() bool { return b.meta&maskWhiteCastleQueenside != 0 }

// CanBlackCastleKingside returns true iff black can castle king side.
func (b Board) CanBlackCastleKingside() bool { return b.meta&maskBlackCastleKingside != 0 }

// CanBlackCastleQueenside returns true iff black can castle queen side.
func (b Board) CanBlackCastleQueenside() bool { return b.meta&maskBlackCastleQueenside != 0 }

// HalfMoves returns the number of half moves (moves by one player) since the
// last pawn moved or piece was captured. This is used for determining if a draw
// can be claimed by the fifty move rule.
func (b Board) HalfMoves() int { return int(b.half) }

// FullMoves returns the number of full moves (moves by both players).
func (b Board) FullMoves() int { return int(b.total/2) + 1 }

// String renders the board from whites perspective.
func (b Board) String() string {
	var sb strings.Builder
	var i uint8
	for i = 0; i < 64; i++ {
		if i != 0 && i%8 == 0 {
			sb.WriteRune('\n')
		}
		poi := PrintOrderedIndex(i)
		r := b.PieceAt(poi).Rune()
		sb.WriteRune(r)
	}
	return sb.String()
}

// Validate returns an error on an inconsistent or invalid board.
func (b Board) Validate() error {
	numWhiteKings := bits.OnesCount64(b.kings & b.white)
	if numWhiteKings != 1 {
		return fmt.Errorf("invalid board: %d white kings", numWhiteKings)
	}
	numBlackKings := bits.OnesCount64(b.kings & b.black)
	if numBlackKings != 1 {
		return fmt.Errorf("invalid board: %d black kings", numBlackKings)
	}
	if b.pawns&maskRank1 != 0 {
		return errors.New("invalid board: pawns on rank 1")
	}
	if b.pawns&maskRank8 != 0 {
		return errors.New("invalid board: pawns on rank 8")
	}
	return nil
}
