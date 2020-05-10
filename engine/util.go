package engine

import (
	"fmt"
	"io"
	"math/bits"
)

// ToAlgebraicNotation converts the index i to algebraic notation (e.g. A1, H8).
// Results for indexes outside of the sane 0...63 range are undefined.
func ToAlgebraicNotation(i uint8) string {
	file := File(i)
	rank := Rank(i)
	return string([]byte{'a' + file, '1' + rank})
}

// ParseAlgebraicNotation reads two bytes from r and parses them as Algebraic
// Notation, returning the rank and file (both zero indexed).
func ParseAlgebraicNotation(r io.RuneReader) (rank, file uint8, err error) {
	file, err = parseFile(r)
	if err != nil {
		return 0, 0, err
	}
	rank, err = parseRank(r)
	if err != nil {
		return 0, 0, err
	}
	return rank, file, err
}

func parseFile(r io.RuneReader) (uint8, error) {
	ch, _, err := r.ReadRune() // read file
	if err != nil {
		return 0, unexpectingEOF(err)
	}

	switch ch {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
		return uint8(ch - 'a'), nil
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
		return uint8(ch - 'A'), nil
	default:
		return 0, fmt.Errorf("unexpected '%c', expecting [a-hA-H]", ch)
	}
}

func parseRank(r io.RuneReader) (uint8, error) {
	ch, _, err := r.ReadRune() // read rank
	if err != nil {
		return 0, unexpectingEOF(err)
	}

	switch ch {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		return uint8(ch - '1'), nil
	default:
		return 0, fmt.Errorf("unexpected '%c', expecting [1-8]", ch)
	}
}

func unexpectingEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

// Rank returns the rank index (0...7) for a given square index.
func Rank(i uint8) uint8 {
	return i / 8
}

// File returns the file index (0 for A, ..., 7 for H) for a given square index.
func File(i uint8) uint8 {
	return i % 8
}

// Square returns the square index (0 for A1, 63 for H8) for a rank and file.
func Square(rank, file uint8) uint8 {
	return rank*8 + file
}

// PrintOrderedIndex takes an index i (which must be in the range 0...63), and
// returns the index with its rank "reversed" such that PrintOrderedIndex(i)
// loops through rank 8 first, then 7, 6, ..., 1.
//
// A1 is the 0th index in the bitboard, and H8 is the 63rd index. If we were to
// print in 0...63 index order we would print the 1st rank first, then the 2nd
// rank, 3rd, 4th, ..., 8th. However when outputting the board we want to print
// it so that we view it from white's perspective, with the 1st rank at the
// bottom (printed last) and the 8th rank at the top (printed first). This is
// also required when printing and parsing FEN, where the 8th rank is again
// rendered first.
//
// The floor of i divided by 8 (number of files) gets us the current rank, in
// the range 0...7 (i.e. indexed by 0).
//
// 7 (max rank if we're indexing by 0) minus the current rank gets us ranks in
// reverse order (0 ⇨ 7, 1 ⇨ 6, ..., 7 ⇨ 0).
//
// i modulus 8 (number of ranks) gets us the current file, again in the range
// 0...7.
//
// The corresponding "reversed rank" index for an input index is then the
// reverse rank times 8 (ranks) plus the current file.
//
// The full formula is thus:
//
//       8×( 7 -    ⌊i÷8⌋) + (i mod 8)
//     = 8×( 7 -    ⌊i÷8⌋) + i - 8×⌊i÷8⌋
//     =    56 -  8×⌊i÷8⌋  + i - 8×⌊i÷8⌋
//     =    56 - 16×⌊i÷8⌋  + i
//
// This gets us indexes (for input 0...63) in the following order:
//
//     56, 57, 58, 59, 60, 61, 62, 63, (rank 8)
//     48, 49, 50, ...,            55,      (7)
//     40, 41, 42, ...,            47,      (6)
//     32, 33, 34, ...,            39,      (5)
//     24, 25, 26, ...,            31,      (4)
//     16, 17, 18, ...,            23,      (3)
//      8,  9, 10, ...,            15,      (2)
//      0,  1,  2, ...,             6.      (1)
func PrintOrderedIndex(i uint8) uint8 {
	return i + 56 - 16*(i/8)
}

func diff(sq1, sq2 uint8) uint8 {
	diff := int(sq1) - int(sq2)
	if diff < 0 {
		diff *= -1
	}
	return uint8(diff)
}

// popLSB finds the Least Significant Bit in x, returning the index of that bit,
// a bitmask with only that bit set, and mutating x to unset that bit.
func popLSB(x *uint64) (uint8, uint64) {
	idx := uint8(bits.TrailingZeros64(*x))
	var bit uint64 = 1 << idx
	*x &^= bit
	return idx, bit
}

// popMSB finds the Most Significant Bit in x, returning the index of that bit,
// a bitmask with only that bit set, and mutating x to unset that bit.
func popMSB(x *uint64) (uint8, uint64) {
	idx := uint8(63 - bits.LeadingZeros64(*x))
	var bit uint64 = 1 << idx
	*x &^= bit
	return idx, bit
}
