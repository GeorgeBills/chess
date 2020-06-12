package engine

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"strings"

	chess "github.com/GeorgeBills/chess/m/v2"
)

// https://www.chessprogramming.org/Bitboards

// InitialBoardFEN is the FEN for a board in the initial state.
const InitialBoardFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

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
	maskWhiteCastleKingside  uint8 = 0b10000000
	maskWhiteCastleQueenside uint8 = 0b01000000
	maskBlackCastleKingside  uint8 = 0b00100000
	maskBlackCastleQueenside uint8 = 0b00010000
	maskCastling             uint8 = maskWhiteCastleKingside | maskWhiteCastleQueenside | maskBlackCastleKingside | maskBlackCastleQueenside
	maskCanEnPassant         uint8 = 0b00001000
	maskEnPassantFile        uint8 = 0b00000111 // the last 3 bits of meta indicate the file (zero indexed) for a valid en passant
)

// NewBoard returns a board in the initial state.
func NewBoard() *Board {
	return &Board{
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

// isWhiteAt returns true if there is a piece at index i and it is white.
func (b Board) isWhiteAt(i uint8) bool { return b.white&(1<<i) != 0 }

// isBlackAt returns true if there is a piece at index i and it is black.
func (b Board) isBlackAt(i uint8) bool { return b.black&(1<<i) != 0 }

// IsEmptyAt returns true if the square at index i is empty.
func (b Board) isEmptyAt(i uint8) bool { return !b.isWhiteAt(i) && !b.isBlackAt(i) }

// isPawnAt returns true if there is a piece at index i and it is a pawn.
func (b Board) isPawnAt(i uint8) bool { return b.pawns&(1<<i) != 0 }

// isKnightAt returns true if there is a piece at index i and it is a knight.
func (b Board) isKnightAt(i uint8) bool { return b.knights&(1<<i) != 0 }

// isBishopAt returns true if there is a piece at index i and it is a bishop.
func (b Board) isBishopAt(i uint8) bool { return b.bishops&(1<<i) != 0 }

// isRookAt returns true if there is a piece at index i and it is a rook.
func (b Board) isRookAt(i uint8) bool { return b.rooks&(1<<i) != 0 }

// isQueenAt returns true if there is a piece at index i and it is a queen.
func (b Board) isQueenAt(i uint8) bool { return b.queens&(1<<i) != 0 }

// isKingAt returns true if there is a piece at index i and it is a king.
func (b Board) isKingAt(i uint8) bool { return b.kings&(1<<i) != 0 }

// EnPassant returns the index of the square under threat of en passant, or
// math.MaxUint8 if there is no such square.
// TODO: return bool here, this isn't used anywhere performance critical
func (b Board) EnPassant() uint8 {
	if b.meta&maskCanEnPassant == 0 {
		return math.MaxUint8
	}
	file := uint8(b.meta & maskEnPassantFile)
	switch b.ToMove() {
	case White:
		return chess.SquareIndex(rank6, file)
	case Black:
		return chess.SquareIndex(rank3, file)
	default:
		panic(fmt.Errorf("invalid to move; %#v", b))
	}
}

// PieceAt returns the piece at index i.
func (b Board) PieceAt(i uint8) Piece {
	if b.isWhiteAt(i) {
		switch {
		case b.isPawnAt(i):
			return PieceWhitePawn
		case b.isKnightAt(i):
			return PieceWhiteKnight
		case b.isBishopAt(i):
			return PieceWhiteBishop
		case b.isRookAt(i):
			return PieceWhiteRook
		case b.isQueenAt(i):
			return PieceWhiteQueen
		case b.isKingAt(i):
			return PieceWhiteKing
		default:
			panic(fmt.Errorf("invalid white piece at index %d; %#v", i, b))
		}
	}
	if b.isBlackAt(i) {
		switch {
		case b.isPawnAt(i):
			return PieceBlackPawn
		case b.isKnightAt(i):
			return PieceBlackKnight
		case b.isBishopAt(i):
			return PieceBlackBishop
		case b.isRookAt(i):
			return PieceBlackRook
		case b.isQueenAt(i):
			return PieceBlackQueen
		case b.isKingAt(i):
			return PieceBlackKing
		default:
			panic(fmt.Errorf("invalid black piece at index %d; %#v", i, b))
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

// CanWhiteCastleKingside returns true if white can castle kingside.
func (b Board) CanWhiteCastleKingside() bool { return b.meta&maskWhiteCastleKingside != 0 }

// CanWhiteCastleQueenside returns true if white can castle queenside.
func (b Board) CanWhiteCastleQueenside() bool { return b.meta&maskWhiteCastleQueenside != 0 }

// CanBlackCastleKingside returns true if black can castle kingside.
func (b Board) CanBlackCastleKingside() bool { return b.meta&maskBlackCastleKingside != 0 }

// CanBlackCastleQueenside returns true if black can castle queenside.
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
		poi := chess.PrintOrderedIndex(i)
		r := b.PieceAt(poi).Rune()
		sb.WriteRune(r)
	}
	return sb.String()
}

// ColourFlipped returns a board such that all the black pieces are now white
// pieces and all the white pieces are now black pieces, and the players have
// switched sides.
func (b Board) ColourFlipped() Board {
	// https://www.chessprogramming.org/Color_Flipping

	b.white, b.black = b.black, b.white // swap colours

	b.meta ^= maskCastling // swap castling

	var i uint8
	for i = 0; i < 32; i++ {
		j := chess.PrintOrderedIndex(i)

		// swap position of i and j
		pi, pj := b.PieceAt(i), b.PieceAt(j)
		b.removePieceAt(i)
		b.removePieceAt(j)
		b.setPieceAt(i, pj)
		b.setPieceAt(j, pi)
	}
	return b
}

func (b *Board) removePieceAt(i uint8) {
	var bit uint64 = 1 << i
	b.black &^= bit
	b.white &^= bit
	b.pawns &^= bit
	b.knights &^= bit
	b.bishops &^= bit
	b.rooks &^= bit
	b.queens &^= bit
	b.kings &^= bit
}

// setPieceAt sets piece p at index i. It does not remove any existing piece
// first: calling setPieceAt() with an existing piece on the square will result
// in an invalid board.
func (b *Board) setPieceAt(i uint8, p Piece) {
	var bit uint64 = 1 << i
	switch p {
	case PieceWhitePawn:
		b.white |= bit
		b.pawns |= bit
	case PieceWhiteKnight:
		b.white |= bit
		b.knights |= bit
	case PieceWhiteBishop:
		b.white |= bit
		b.bishops |= bit
	case PieceWhiteRook:
		b.white |= bit
		b.rooks |= bit
	case PieceWhiteQueen:
		b.white |= bit
		b.queens |= bit
	case PieceWhiteKing:
		b.white |= bit
		b.kings |= bit
	case PieceBlackPawn:
		b.black |= bit
		b.pawns |= bit
	case PieceBlackKnight:
		b.black |= bit
		b.knights |= bit
	case PieceBlackBishop:
		b.black |= bit
		b.bishops |= bit
	case PieceBlackRook:
		b.black |= bit
		b.rooks |= bit
	case PieceBlackQueen:
		b.black |= bit
		b.queens |= bit
	case PieceBlackKing:
		b.black |= bit
		b.kings |= bit
	case PieceNone:
		// do nothing
	default:
		panic(fmt.Errorf("invalid piece to set: %b", p))
	}
}

// Validate returns an error on an inconsistent or invalid board.
func (b Board) Validate() error {
	numWhiteKings := bits.OnesCount64(b.kings & b.white)
	if numWhiteKings != 1 {
		return fmt.Errorf("%d white kings", numWhiteKings)
	}
	numBlackKings := bits.OnesCount64(b.kings & b.black)
	if numBlackKings != 1 {
		return fmt.Errorf("%d black kings", numBlackKings)
	}
	if b.pawns&maskRank1 != 0 {
		return errors.New("pawns on rank 1")
	}
	if b.pawns&maskRank8 != 0 {
		return errors.New("pawns on rank 8")
	}
	if b.CanBlackCastleKingside() && (b.PieceAt(E8) != PieceBlackKing || b.PieceAt(H8) != PieceBlackRook) {
		return errors.New("invalid black castling: kingside implies king at E8 and rook at H8")
	}
	if b.CanBlackCastleQueenside() && (b.PieceAt(E8) != PieceBlackKing || b.PieceAt(A8) != PieceBlackRook) {
		return errors.New("invalid black castling: queenside implies king at E8 and rook at A8")
	}
	if b.CanWhiteCastleKingside() && (b.PieceAt(E1) != PieceWhiteKing || b.PieceAt(H1) != PieceWhiteRook) {
		return errors.New("invalid white castling: kingside implies king at E1 and rook at H1")
	}
	if b.CanWhiteCastleQueenside() && (b.PieceAt(E1) != PieceWhiteKing || b.PieceAt(A1) != PieceWhiteRook) {
		return errors.New("invalid white castling: queenside implies king at E1 and rook at A1")
	}
	return nil
}
