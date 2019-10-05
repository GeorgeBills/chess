package engine

// ToAlgebraicNotation converts the index i to algebraic notation (e.g. A1, H8).
// Results for indexes outside of the sane 0...63 range are undefined.
func ToAlgebraicNotation(i uint8) string {
	file := i % 8
	rank := i / 8
	return string([]byte{'a' + file, '1' + rank})
}
