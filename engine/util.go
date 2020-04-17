package engine

// ToAlgebraicNotation converts the index i to algebraic notation (e.g. A1, H8).
// Results for indexes outside of the sane 0...63 range are undefined.
func ToAlgebraicNotation(i uint8) string {
	file := File(i)
	rank := Rank(i)
	return string([]byte{'a' + file, '1' + rank})
}

// Rank returns the rank index (0...7) for a given square index.
func Rank(i uint8) uint8 {
	return i / 8
}

// File returns the file index (0 for A, ..., 7 for H) for a given square index.
func File(i uint8) uint8 {
	return i % 8
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
