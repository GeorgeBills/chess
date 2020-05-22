package engine

import (
	"math/bits"
)

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
