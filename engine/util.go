package engine

// ToAlgebraicNotation converts the index i to algebraic notation (e.g. A1, H8).
// Results for indexes outside of the sane 0...63 range are undefined.
func ToAlgebraicNotation(i uint8) string {
	file := i % 8
	rank := i / 8
	return string([]byte{'a' + file, '1' + rank})
}

// Rank returns the rank number (0...7) for a given index.
func Rank(i uint8) uint8 {
	return i / 8
}

// File returns the file number (0 for A....7 for H) for a given index.
func File(i uint8) uint8 {
	return i % 8
}
