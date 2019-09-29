package engine

// n: north; e: east; s: south; w: west
// nn, ss: north north and south south (used by pawns performing double moves)
// nne, een, ssw: north north east, etc (used by knights)
//
// North and South are easy; add or subtract 8 (a full rank). With bit shifting
// you don't even need to check if that's off the board: shifting above 64 will
// shift the bits out, and subtracting past 0 with an unsigned int ends up
// wrapping into a very large shift up and off the board again.
//
// East and West require checking if you'll leave the board or not, since
// leaving the board in either of those directions will wrap around; e.g. a king
// shouldn't be able to move one square West from A2 (8) and end up on H1 (7).

// WhitePawnMoves returns the moves a white pawn at index i can make, ignoring
// captures and en passant.
func WhitePawnMoves(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 8) // n
	// if a white pawn is on rank 2 it may move two squares
	if i/8 == 1 {
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
	if i/8 == 6 {
		moves |= 1 << (i - 16) // ss
	}
	return moves
}

// KingMoves returns the moves a king at index i can make, ignoring castling.
func KingMoves(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 8) // n
	moves |= 1 << (i - 8) // s
	// can't move east if we're on file h
	if i%8 != 7 {
		moves |= 1 << (i + 1) // e
	}
	// can't move west if we're on file a
	if i%8 != 0 {
		moves |= 1 << (i - 1) // w
	}
	return moves
}

// KnightMoves returns the moves a knight at index i can make.
func KnightMoves(i uint8) uint64 {
	var moves uint64
	file := i % 8
	if file > 0 {
		moves |= 1 << (i + 15) // nnw (+2×8, -1)
		moves |= 1 << (i - 17) // ssw (-2×8, -1)
		if file > 1 {
			moves |= 1 << (i + 6)  // wwn (+8, -2×1)
			moves |= 1 << (i - 10) // wws (-8, -2×1)
		}
	}
	if file < 7 {
		moves |= 1 << (i + 17) // nne (+2×8, +1)
		moves |= 1 << (i - 15) // sse (-2×8, +1)
		if file < 6 {
			moves |= 1 << (i + 10) // een (+8, +2×1)
			moves |= 1 << (i - 6)  // ees (-8, +2×1)
		}
	}
	return moves
}

// Moves returns a slice of possible board states from the current board state.
func (b Board) Moves() []Board {
	var moves []Board
	var from, to uint8
	var frombit, tobit uint64
	var colour uint64

	tomove := b.ToMove()
	if tomove == White {
		colour = b.white
	} else {
		colour = b.black
	}

	// - Find all pieces of the given colour.
	// - For each square, check if we have a piece there.
	// - If there is a piece there, find all the moves that piece can make.
	// - AND NOT the moves with our colour (we can't move to a square we occupy).
	// - For each remaining move, generate the board state for that move:
	//   - Remove the piece from the square it currently occupies.
	//   - Remove any opposing pieces from the target square.
	//   - Place the piece in the target square.

	// We evaluate pieces in queen, king, rook, bishop, knight, pawn order as a
	// heuristic to roughly sort more likely to be impactful moves early, in
	// order to aid alpha-beta pruning.

	// TODO: queen moves

	// TODO: king moves

	// TODO: castling

	// TODO: rook moves

	// TODO: bishop moves

	// TODO: knight moves

	pawns := b.pawns & colour
	// ignore ranks 1 and 8, pawns can't ever occupy them
	for from = 8; from < 56; from++ {
		frombit = 1 << from
		if pawns&frombit != 0 { // is there a pawn on this square?
			var pawnmoves uint64
			if tomove == White {
				pawnmoves = WhitePawnMoves(from) &^ colour
			} else {
				pawnmoves = BlackPawnMoves(from) &^ colour
			}
			for to = 0; to < 64; to++ {
				tobit = 1 << to
				if pawnmoves&tobit != 0 { // is there a move to this square?
					newboard := b

					// remove piece
					newboard.pawns ^= frombit

					// remove target
					newboard.bishops ^= tobit
					newboard.knights ^= tobit
					newboard.queens ^= tobit
					newboard.rooks ^= tobit

					// place piece
					newboard.pawns |= tobit

					if tomove == White {
						newboard.white ^= frombit // remove piece
						newboard.black ^= tobit   // remove target
						newboard.white |= tobit   // place piece
					} else {
						newboard.black ^= frombit // remove piece
						newboard.white ^= tobit   // remove target
						newboard.black |= tobit   // place piece
					}

					// TODO: pawn captures
					// TODO: en passant captures
					// TODO: set en passant target on double pawn moves
					// TODO: pawn promotion

					newboard.total++

					moves = append(moves, newboard)
				}
			}
		}
	}

	// TODO: record half moves

	// TODO: disallow moves placing the king in check

	return moves
}
