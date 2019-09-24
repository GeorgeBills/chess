package engine

import "strings"

// Board represents an 8×8 chess board.
//
// The 0th index represents A1 and the 64th index represents H8.
type Board struct {
	black   uint64
	white   uint64
	pawns   uint64
	knights uint64
	bishops uint64
	rooks   uint64
	queens  uint64
	kings   uint64
}

// NewBoard returns a board in the initial state.
func NewBoard() Board {
	return Board{
		white:   0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
		black:   0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_11111111,
		pawns:   0b00000000_11111111_00000000_00000000_00000000_00000000_11111111_00000000,
		knights: 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
		bishops: 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
		rooks:   0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
		queens:  0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
		kings:   0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
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
			panic(b) // invalid board state
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
			panic(b) // invalid board state
		}
	}
	return 0
}

func (b Board) String() string {
	var sb strings.Builder
	var i uint8
	for i = 0; i < 64; i++ {
		if i != 0 && i%8 == 0 {
			sb.WriteRune('\n')
		}
		r := b.PieceAt(i).Rune()
		sb.WriteRune(r)
	}
	return sb.String()
}
