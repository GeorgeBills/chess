package engine

import (
	"strings"
)

// Move represents a chess move.
type Move uint16

// NewMove returns a new move which is not a capture, promotion or castling.
func NewMove(from, to uint8) Move {
	return Move(uint16(from)<<6 | uint16(to))
}

// NewCapture returns a new move which represents a capture.
func NewCapture(from, to uint8) Move {
	return NewMove(from, to) | moveIsCapture
}

// NewEnPassant returns a new move which represents a capture en passant.
func NewEnPassant(from, to uint8) Move {
	return NewMove(from, to) | moveIsEnPassant | moveIsCapture
}

// Castling moves are represented with constants.
const (
	BlackKingsideCastle  = Move(uint16(E8)<<6|uint16(G8)) | moveIsKingsideCastle
	BlackQueensideCastle = Move(uint16(E8)<<6|uint16(C8)) | moveIsQueensideCastle
	WhiteKingsideCastle  = Move(uint16(E1)<<6|uint16(G1)) | moveIsKingsideCastle
	WhiteQueensideCastle = Move(uint16(E1)<<6|uint16(C1)) | moveIsQueensideCastle
)

func newPromotion(from, to uint8, capture bool) Move {
	if capture {
		return NewCapture(from, to)
	}
	return NewMove(from, to)
}

// NewQueenPromotion returns a new move where the move represents a pawn
// promoting to a queen.
func NewQueenPromotion(from, to uint8, capture bool) Move {
	return newPromotion(from, to, capture) | moveIsQueenPromotion
}

// NewKnightPromotion returns a new move which represents a pawn promoting to a
// knight.
func NewKnightPromotion(from, to uint8, capture bool) Move {
	return newPromotion(from, to, capture) | moveIsKnightPromotion
}

// NewRookPromotion returns a new move which represents a pawn promoting to a
// rook.
func NewRookPromotion(from, to uint8, capture bool) Move {
	return newPromotion(from, to, capture) | moveIsRookPromotion
}

// NewBishopPromotion returns a new move which represents a pawn promoting to a
// bishop.
func NewBishopPromotion(from, to uint8, capture bool) Move {
	return newPromotion(from, to, capture) | moveIsBishopPromotion
}

// SAN returns the move in Somewhat Algebraic Notation, which is very similar to
// (but not quite the same as) Standard Algebraic Notation.
//
// The main difference is that Standard Algebraic Notation indicates a piece
// moving to a square (e.g. "Ka1" to represent the king moving to the A1
// square), omitting the source square when it is unambiguous to do so (i.e.
// when there's only one piece of the type indicated that could have made the
// move).
//
// Somewhat Algebraic Notation always includes both source square and target
// square (similarly to Pure Coordinate Notation), and never includes the piece
// type. Piece type can be unambiguously determined from the source square and
// the current state of the board.
func (m Move) SAN() string {
	if m.IsKingsideCastling() {
		return "O-O"
	}
	if m.IsQueensideCastling() {
		return "O-O-O"
	}
	var san strings.Builder
	san.WriteString(ToAlgebraicNotation(m.From()))
	if m.IsCapture() {
		san.WriteByte('x')
	}
	san.WriteString(ToAlgebraicNotation(m.To()))
	switch {
	case m&moveMetaMask == moveIsEnPassant:
		san.WriteString("e.p.")
	case m&moveIsQueenPromotion == moveIsQueenPromotion:
		san.WriteString("=Q")
	case m&moveIsKnightPromotion == moveIsKnightPromotion:
		san.WriteString("=N")
	case m&moveIsRookPromotion == moveIsRookPromotion:
		san.WriteString("=R")
	case m&moveIsBishopPromotion == moveIsBishopPromotion:
		san.WriteString("=B")
	}
	return san.String()
}

const (
	moveMetaMask          = 0b11110000_00000000 // 4 bits: meta (promotion, capture, castle, ...?)
	moveFromMask          = 0b00001111_11000000 // 6 bits: from square
	moveToMask            = 0b00000000_00111111 // 6 bits: to square
	moveIsPromotion       = 0b1000 << 12
	moveIsQueenPromotion  = (moveIsPromotion | 0b0011<<12)
	moveIsKnightPromotion = (moveIsPromotion | 0b0010<<12)
	moveIsRookPromotion   = (moveIsPromotion | 0b0001<<12)
	moveIsBishopPromotion = (moveIsPromotion | 0b0000<<12)
	moveIsCapture         = 0b0100 << 12
	moveIsEnPassant       = (moveIsCapture | 0b0001<<12)
	moveIsKingsideCastle  = 0b0011 << 12
	moveIsQueensideCastle = 0b0010 << 12
	moveIsPawnDoubleMove  = 0b0001 << 12
)

// IsCapture returns true iff the move represents a capture.
func (m Move) IsCapture() bool {
	return m&moveIsCapture == moveIsCapture
}

// IsEnPassant returns true iff the move represents a capture en passant.
func (m Move) IsEnPassant() bool {
	return m&moveIsEnPassant == moveIsEnPassant
}

// IsPromotion returns true iff the move represents a pawn promotion.
func (m Move) IsPromotion() bool {
	return m&moveIsPromotion == moveIsPromotion
}

// IsKingsideCastling returns true iff the move represents kingside castling.
func (m Move) IsKingsideCastling() bool {
	return m&moveMetaMask == moveIsKingsideCastle
}

// IsQueensideCastling returns true iff the move represents queenside castling.
func (m Move) IsQueensideCastling() bool {
	return m&moveMetaMask == moveIsQueensideCastle
}

// From returns the from index for the move.
func (m Move) From() uint8 {
	return uint8((m & moveFromMask) >> 6)
}

// To returns the to index for the move.
func (m Move) To() uint8 {
	return uint8((m & moveToMask) >> 0)
}

// MakeMove applies move to the board, updating its state.
func (b *Board) MakeMove(move Move) {
	// switch {
	// case move.IsEnPassant():
	// case move.IsKingsideCastling():
	// case move.IsQueensideCastling():
	// case move.IsPromotion():
	// default:
	// }

	var frombit uint64 = 1 << move.From()
	var tobit uint64 = 1 << move.To()

	// remove any opposing piece on our destination square
	b.pawns &^= tobit
	b.knights &^= tobit
	b.bishops &^= tobit
	b.rooks &^= tobit
	b.queens &^= tobit
	b.kings &^= tobit

	// update colour masks
	switch {
	case b.white&frombit != 0:
		b.white &^= frombit
		b.white |= tobit
		b.black &^= tobit
	case b.black&frombit != 0:
		b.black &^= frombit
		b.black |= tobit
		b.white &^= tobit
	}

	// update relevant piece mask
	switch {
	case b.pawns&frombit != 0:
		b.pawns &^= frombit
		b.pawns |= tobit
	case b.bishops&frombit != 0:
		b.bishops &^= frombit
		b.bishops |= tobit
	case b.knights&frombit != 0:
		b.knights &^= frombit
		b.knights |= tobit
	case b.rooks&frombit != 0:
		b.rooks &^= frombit
		b.rooks |= tobit
	case b.queens&frombit != 0:
		b.queens &^= frombit
		b.queens |= tobit
	case b.kings&frombit != 0:
		b.kings &^= frombit
		b.kings |= tobit
	}

	b.total++
}

// // UnmakeMove unapplies move to the board, undoing the state change.
// func (b Board) UnmakeMove(move Move) {
// }
