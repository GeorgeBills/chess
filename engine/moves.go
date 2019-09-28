package engine

// n: north; e: east; s: south; w: west
// nn, ss: north north and south south (used by pawns performing double moves)
// nne, een, ssw: north north east, etc (used by knights)
//
// north and south are easy; add or subtract 8 (a full rank). with bitshifting
// you don't even need to check if that's off the board, since bits that leave
// the board will just be shifted out.

// WhitePawnMoves returns the moves a white pawn at index i can make, ignoring
// captures and en passant.
func WhitePawnMoves(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 8) // n
	// if a white pawn is on rank 2 it may move two squares
	if 8 <= i && i <= 15 {
		moves |= 1 << (i + 16) // nn
	}
	return moves
}

// BlackPawnMoves returns the moves a black pawn at index i can make, ignoring
// captures and en passant.
func BlackPawnMoves(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i - 8) // s
	// if a black pawn is on rank 7 it may move two squares
	if 48 <= i && i <= 55 {
		moves |= 1 << (i - 16) // ss
	}
	return moves
}
