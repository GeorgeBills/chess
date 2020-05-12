package engine

import (
	"fmt"
	"io"
	"strings"
)

// https://www.chessprogramming.org/Encoding_Moves

// Move represents a chess move.
type Move uint16

// NewMove returns a new move which is not a capture, promotion or castling.
func NewMove(from, to uint8) Move {
	return Move(uint16(from)<<6 | uint16(to))
}

// NewPawnDoublePush returns a new move which represents a pawn double push.
func NewPawnDoublePush(from, to uint8) Move {
	return NewMove(from, to) | moveIsPawnDoubleMove
}

// NewCapture returns a new move which represents a capture.
func NewCapture(from, to uint8) Move {
	return NewMove(from, to) | moveIsCapture
}

// NewEnPassant returns a new move which represents a capture en passant.
func NewEnPassant(from, to uint8) Move {
	return NewMove(from, to) | moveIsEnPassant
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
	// TODO: convert to UCIN, we can provide proper SAN later on; this is neither
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

// FIXME: en passant overlaps with promotion | capture, so ordering of evaluation matters...

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

// IsCapture returns true if the move represents a capture.
func (m Move) IsCapture() bool {
	return m&moveIsCapture == moveIsCapture
}

// IsEnPassant returns true if the move represents a capture en passant.
func (m Move) IsEnPassant() bool {
	return m&moveIsEnPassant == moveIsEnPassant
}

// IsPromotion returns true if the move represents a pawn promotion.
func (m Move) IsPromotion() bool {
	return m&moveIsPromotion == moveIsPromotion
}

// IsKingsideCastling returns true if the move represents kingside castling.
func (m Move) IsKingsideCastling() bool {
	return m&moveMetaMask == moveIsKingsideCastle
}

// IsQueensideCastling returns true if the move represents queenside castling.
func (m Move) IsQueensideCastling() bool {
	return m&moveMetaMask == moveIsQueensideCastle
}

// IsPawnDoublePush returns true if the move represents a pawn double push.
func (m Move) IsPawnDoublePush() bool {
	return m&moveMetaMask == moveIsPawnDoubleMove
}

// From returns the from index for the move.
func (m Move) From() uint8 {
	return uint8((m & moveFromMask) >> 6)
}

// To returns the to index for the move.
func (m Move) To() uint8 {
	return uint8((m & moveToMask) >> 0)
}

type Game struct {
	*Board  // TODO: rename to BoardWithHistory or similar
	history []moveCapture
}

type moveCapture struct {
	Move
	capture      Piece
	previousMeta byte
}

func NewGame(b *Board) Game {
	return Game{
		Board:   b,
		history: make([]moveCapture, 0, 128),
	}
}

const (
	uciWhiteKingsideCastle  = "e1g1"
	uciWhiteQueensideCastle = "e1c1"
	uciBlackKingsideCastle  = "e8g8"
	uciBlackQueensideCastle = "e8c8"
)

// ParseNewMoveFromUCIN parses a new move from Universal Chess Interface
// Notation. UCIN is very similar to Long Algebraic Notation, but omits the
// hyphen, the moving piece (can be inferred from the "from" AN) and whether the
// move is a capture (can be inferred from the "to" AN and the current state of
// the board).
func (b *Board) ParseNewMoveFromUCIN(r io.RuneReader) (Move, error) {
	fromRank, fromFile, err := ParseAlgebraicNotation(r)
	if err != nil {
		return 0, err
	}
	fromSq := Square(fromRank, fromFile)

	toRank, toFile, err := ParseAlgebraicNotation(r)
	if err != nil {
		return 0, err
	}
	toSq := Square(toRank, toFile)

	isCapture := !b.isEmptyAt(toSq)

	if b.isPawnAt(fromSq) {
		diff := diff(fromSq, toSq)

		if diff == 16 {
			// pawn moving exactly two ranks: must be a double push
			return NewPawnDoublePush(fromSq, toSq), nil
		}

		if diff != 8 && b.isEmptyAt(toSq) {
			// pawn capturing to an empty square: must be en passant
			return NewEnPassant(fromSq, toSq), nil
		}

		if toRank == rank1 || toRank == rank8 {
			// pawn moving to rank 1 or 8; must be a promotion
			ch, _, err := r.ReadRune()
			if err != nil {
				return 0, err
			}
			switch ch {
			case 'q':
				return NewQueenPromotion(fromSq, toSq, isCapture), nil
			case 'r':
				return NewRookPromotion(fromSq, toSq, isCapture), nil
			case 'n':
				return NewKnightPromotion(fromSq, toSq, isCapture), nil
			case 'b':
				return NewBishopPromotion(fromSq, toSq, isCapture), nil
			default:
				return 0, fmt.Errorf("unexpected promotion '%c', expecting [qrnb]", ch)
			}
		}
	}

	if isCapture {
		return NewCapture(fromSq, toSq), nil
	}

	if b.isKingAt(fromSq) {
		if fromSq == E1 {
			if toSq == G1 { // "e1g1" is white kingside castling
				return WhiteKingsideCastle, nil
			}
			if toSq == C1 { // "e1c1" is white queenside castling
				return WhiteQueensideCastle, nil
			}
		}

		if fromSq == E8 {
			if toSq == G8 { // "e8g8" is black kingside castling
				return BlackKingsideCastle, nil
			}
			if toSq == C8 { // "e8c8" is black queensdie castling
				return BlackQueensideCastle, nil
			}
		}
	}

	return NewMove(fromSq, toSq), nil
}

// MakeMove applies move to the board, updating its state.
func (g *Game) MakeMove(move Move) {
	tomove := g.ToMove()
	from, to := move.From(), move.To()
	var frombit, tobit uint64 = 1 << from, 1 << to

	mc := moveCapture{
		Move:         move,
		capture:      g.PieceAt(to),
		previousMeta: g.meta,
	}
	g.history = append(g.history, mc)

	// remove castling rights if we need to
	switch from {
	case A1: // white queenside rook starting square
		g.meta &^= maskWhiteCastleQueenside
	case E1: // white king starting square
		g.meta &^= maskWhiteCastleKingside | maskWhiteCastleQueenside
	case H1: // white kingside rook starting square
		g.meta &^= maskWhiteCastleKingside
	case A8: // black queenside rook starting square
		g.meta &^= maskBlackCastleQueenside
	case E8: // black king starting square
		g.meta &^= maskBlackCastleKingside | maskBlackCastleQueenside
	case H8: // black kingside rook starting square
		g.meta &^= maskBlackCastleKingside
	}
	switch to {
	case A1: // white queenside rook starting square
		g.meta &^= maskWhiteCastleQueenside
	case H1: // white kingside rook starting square
		g.meta &^= maskWhiteCastleKingside
	case A8: // black queenside rook starting square
		g.meta &^= maskBlackCastleQueenside
	case H8: // black kingside rook starting square
		g.meta &^= maskBlackCastleKingside
	}

	// clear en passant if any was present
	g.meta &^= maskCanEnPassant | maskEnPassantFile

	// TODO: we can short circuit out (via return) in a lot of these special cases
	//       or just move lots of code into the switch default case?
	switch {
	case move.IsPawnDoublePush():
		g.meta |= maskCanEnPassant
		g.meta |= File(from)
	case move.IsKingsideCastling():
		switch tomove {
		case White:
			const togglebits uint64 = 1<<F1 | 1<<H1
			g.white ^= togglebits
			g.rooks ^= togglebits
			g.meta &^= maskWhiteCastleKingside | maskWhiteCastleQueenside
		case Black:
			const togglebits uint64 = 1<<F8 | 1<<H8
			g.black ^= togglebits
			g.rooks ^= togglebits
			g.meta &^= maskBlackCastleKingside | maskBlackCastleQueenside
		}
	case move.IsQueensideCastling():
		switch tomove {
		case White:
			const togglebits uint64 = 1<<A1 | 1<<D1
			g.white ^= togglebits
			g.rooks ^= togglebits
			g.meta &^= maskWhiteCastleKingside | maskWhiteCastleQueenside
		case Black:
			const togglebits uint64 = 1<<A8 | 1<<D8
			g.black ^= togglebits
			g.rooks ^= togglebits
			g.meta &^= maskBlackCastleKingside | maskBlackCastleQueenside
		}
	case move.IsPromotion():
		// swap our pawn out for the piece it's promoting to before it moves
		g.pawns &^= frombit
		switch {
		case move&moveIsQueenPromotion == moveIsQueenPromotion:
			g.queens |= frombit
		case move&moveIsKnightPromotion == moveIsKnightPromotion:
			g.knights |= frombit
		case move&moveIsRookPromotion == moveIsRookPromotion:
			g.rooks |= frombit
		case move&moveIsBishopPromotion == moveIsBishopPromotion:
			g.bishops |= frombit
		default:
			panic(fmt.Errorf("promotion to unknown piece: %b", move))
		}
	case move.IsEnPassant():
		switch tomove {
		case White:
			epCaptureSq := to - 8
			g.black &^= 1 << epCaptureSq
			g.pawns &^= 1 << epCaptureSq
		case Black:
			epCaptureSq := to + 8
			g.white &^= 1 << epCaptureSq
			g.pawns &^= 1 << epCaptureSq
		}
	}

	// remove any opposing piece on our destination square
	g.pawns &^= tobit
	g.knights &^= tobit
	g.bishops &^= tobit
	g.rooks &^= tobit
	g.queens &^= tobit
	g.kings &^= tobit

	// TODO: use (and document) the from|to xor trick throughout

	// update colour masks
	switch tomove {
	case White:
		g.white &^= frombit
		g.white |= tobit
		g.black &^= tobit
	case Black:
		g.black &^= frombit
		g.black |= tobit
		g.white &^= tobit
	}

	// update relevant piece mask
	switch {
	case g.pawns&frombit != 0:
		g.pawns &^= frombit
		g.pawns |= tobit
	case g.bishops&frombit != 0:
		g.bishops &^= frombit
		g.bishops |= tobit
	case g.knights&frombit != 0:
		g.knights &^= frombit
		g.knights |= tobit
	case g.rooks&frombit != 0:
		g.rooks &^= frombit
		g.rooks |= tobit
	case g.queens&frombit != 0:
		g.queens &^= frombit
		g.queens |= tobit
	case g.kings&frombit != 0:
		g.kings &^= frombit
		g.kings |= tobit
	}

	g.total++
}

// UnmakeMove unapplies the most recent move on the board.
func (g *Game) UnmakeMove() {
	tomove := g.ToMove()
	move := g.history[len(g.history)-1]
	g.history = g.history[0 : len(g.history)-1]

	from, to := move.To(), move.From() // flip from and to
	var frombit, tobit uint64 = 1 << from, 1 << to

	// restore previous meta
	g.meta = move.previousMeta

	switch {
	case move.IsKingsideCastling():
		switch tomove {
		case Black:
			const togglebits uint64 = 1<<F1 | 1<<H1
			g.white ^= togglebits
			g.rooks ^= togglebits
		case White:
			const togglebits uint64 = 1<<F8 | 1<<H8
			g.black ^= togglebits
			g.rooks ^= togglebits
		}
	case move.IsQueensideCastling():
		switch tomove {
		case Black:
			const togglebits uint64 = 1<<A1 | 1<<D1
			g.white ^= togglebits
			g.rooks ^= togglebits
		case White:
			const togglebits uint64 = 1<<A8 | 1<<D8
			g.black ^= togglebits
			g.rooks ^= togglebits
		}
	case move.IsPromotion():
		// "unpromote" our promoted piece
		g.pawns |= frombit
		switch {
		case move.Move&moveIsQueenPromotion == moveIsQueenPromotion:
			g.queens &^= frombit
		case move.Move&moveIsKnightPromotion == moveIsKnightPromotion:
			g.knights &^= frombit
		case move.Move&moveIsRookPromotion == moveIsRookPromotion:
			g.rooks &^= frombit
		case move.Move&moveIsBishopPromotion == moveIsBishopPromotion:
			g.bishops &^= frombit
		default:
			panic(fmt.Errorf("promotion to unknown piece: %b", move))
		}
	case move.IsEnPassant():
		var epCaptureSq uint8
		switch tomove {
		case Black:
			epCaptureSq = from - 8
			g.black |= 1 << epCaptureSq
			g.pawns |= 1 << epCaptureSq
		case White:
			epCaptureSq = from + 8
			g.white |= 1 << epCaptureSq
			g.pawns |= 1 << epCaptureSq
		}
	}

	// update colour masks
	switch {
	case g.white&frombit != 0:
		g.white &^= frombit
		g.white |= tobit
		g.black &^= tobit
	case g.black&frombit != 0:
		g.black &^= frombit
		g.black |= tobit
		g.white &^= tobit
	}

	// update relevant piece mask
	switch {
	case g.pawns&frombit != 0:
		g.pawns &^= frombit
		g.pawns |= tobit
	case g.bishops&frombit != 0:
		g.bishops &^= frombit
		g.bishops |= tobit
	case g.knights&frombit != 0:
		g.knights &^= frombit
		g.knights |= tobit
	case g.rooks&frombit != 0:
		g.rooks &^= frombit
		g.rooks |= tobit
	case g.queens&frombit != 0:
		g.queens &^= frombit
		g.queens |= tobit
	case g.kings&frombit != 0:
		g.kings &^= frombit
		g.kings |= tobit
	}

	// resurrect any captured piece
	switch move.capture {
	case PieceWhitePawn:
		g.pawns |= frombit
		g.white |= frombit
	case PieceWhiteKnight:
		g.knights |= frombit
		g.white |= frombit
	case PieceWhiteBishop:
		g.bishops |= frombit
		g.white |= frombit
	case PieceWhiteRook:
		g.rooks |= frombit
		g.white |= frombit
	case PieceWhiteQueen:
		g.queens |= frombit
		g.white |= frombit
	case PieceBlackPawn:
		g.pawns |= frombit
		g.black |= frombit
	case PieceBlackKnight:
		g.knights |= frombit
		g.black |= frombit
	case PieceBlackBishop:
		g.bishops |= frombit
		g.black |= frombit
	case PieceBlackRook:
		g.rooks |= frombit
		g.black |= frombit
	case PieceBlackQueen:
		g.queens |= frombit
		g.black |= frombit
	}

	g.total--
}
