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
	var colour, opposing uint64

	tomove := b.ToMove()
	if tomove == White {
		colour = b.white
		opposing = b.black
	} else {
		colour = b.black
		opposing = b.white
	}

	occupied := b.white | b.black

	// - Find all pieces of the given colour.
	// - For each square, check if we have a piece there.
	// - If there is a piece there, find all the moves that piece can make.
	// - AND NOT the moves with our colour (we can't move to a square we occupy).
	// - For each remaining move, generate the board state for that move:
	//   - Remove the piece from the square it currently occupies.
	//   - Remove any opposing pieces from the target square.
	//   - Place the piece in the target square.

	// Some general optimisations:
	//
	// We don't remove a piece of our type from the target square. e.g. if we're
	// moving a knight to a3 then the bit for a3 on the knights board will end
	// up set no matter what, so don't bother unsetting that bit before the
	// move.
	//
	// We don't check to see if a bit is set or not before clearing or setting
	// it. The branch is more expensive than just setting or clearing
	// regardless.
	//
	// We do all the colour and half move setting in a single if else block,
	// again to minimise branches.
	//
	// More specific optimisations are called out inline.

	// We evaluate pieces in queen, king, rook, bishop, knight, pawn order as a
	// heuristic to roughly sort more likely to be impactful moves early, in
	// order to aid alpha-beta pruning.

	// TODO: queen moves

	king := b.kings & colour
	for from = 0; from < 64; from++ {
		frombit = 1 << from
		if king&frombit != 0 { // is there a king on this square?
			kingmoves := KingMoves(from) &^ colour
			for to = 0; to < 64; to++ {
				tobit = 1 << to
				if kingmoves&tobit != 0 { // is there a move to this square?
					newboard := b

					// remove piece
					newboard.kings &^= frombit

					// remove target
					newboard.bishops &^= tobit
					newboard.knights &^= tobit
					newboard.pawns &^= tobit
					newboard.queens &^= tobit
					newboard.rooks &^= tobit

					// place piece
					newboard.kings |= tobit

					if tomove == White {
						if newboard.black&tobit == 0 { // not a capture: increment halfmoves
							newboard.half++
						}
						newboard.white &^= frombit // remove piece
						newboard.black &^= tobit   // remove target
						newboard.white |= tobit    // place piece
					} else {
						if newboard.white&tobit == 0 { // not a capture: increment halfmoves
							newboard.half++
						}
						newboard.black &^= frombit // remove piece
						newboard.white &^= tobit   // remove target
						newboard.black |= tobit    // place piece
					}

					newboard.total++

					moves = append(moves, newboard)
				}
			}

			// TODO: is it worth breaking out early on the second king?
		}
	}

	// TODO: castling

	rookmove := func() {
		newboard := b
		newboard.rooks &^= frombit // remove piece
		newboard.rooks |= tobit    // place piece
		if tomove == White {
			newboard.half++            // not a capture: increment halfmoves
			newboard.white &^= frombit // remove piece
			newboard.white |= tobit    // place piece
		} else {
			newboard.half++            // not a capture: increment halfmoves
			newboard.black &^= frombit // remove piece
			newboard.black |= tobit    // place piece
		}

		newboard.total++

		moves = append(moves, newboard)
	}
	rookcapture := func() {
		newboard := b

		// remove piece
		newboard.rooks &^= frombit

		// remove target
		newboard.bishops &^= tobit
		newboard.knights &^= tobit
		newboard.pawns &^= tobit
		newboard.queens &^= tobit

		// place piece
		newboard.rooks |= tobit

		if tomove == White {
			newboard.white &^= frombit // remove piece
			newboard.black &^= tobit   // remove target
			newboard.white |= tobit    // place piece
		} else {
			newboard.black &^= frombit // remove piece
			newboard.white &^= tobit   // remove target
			newboard.black |= tobit    // place piece
		}

		newboard.total++

		moves = append(moves, newboard)
	}

	rooks := b.rooks & colour
	for from = 0; from < 64; from++ {
		frombit = 1 << from
		if rooks&frombit != 0 { // is there a rook on this square?
			rank := from / 8
			for n := from + 8; n < 64; n += 8 {
				tobit = 1 << n
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						rookcapture()
					}
					break
				}
				rookmove()
			}
			for e := from + 1; e < (rank+1)*8; e++ {
				tobit = 1 << e
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						rookcapture()
					}
					break
				}
				rookmove()
			}
			for s := from - 8; s > 0 && s < 64; s -= 8 { // uint wraps below 0
				tobit = 1 << s
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						rookcapture()
					}
					break
				}
				rookmove()
			}
			for w := from - 1; w > (rank*8)-1; w-- {
				tobit = 1 << w
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						rookcapture()
					}
					break
				}
				rookmove()
			}
		}
	}

	// TODO: bishop moves

	knights := b.knights & colour
	for from = 0; from < 64; from++ {
		frombit = 1 << from
		if knights&frombit != 0 { // is there a knight on this square?
			knightmoves := KnightMoves(from) &^ colour
			for to = 0; to < 64; to++ {
				tobit = 1 << to
				if knightmoves&tobit != 0 { // is there a move to this square?
					newboard := b

					// remove piece
					newboard.knights &^= frombit

					// remove target
					newboard.bishops &^= tobit
					newboard.pawns &^= tobit
					newboard.queens &^= tobit
					newboard.rooks &^= tobit

					// place piece
					newboard.knights |= tobit

					if tomove == White {
						if newboard.black&tobit == 0 { // not a capture: increment halfmoves
							newboard.half++
						}
						newboard.white &^= frombit // remove piece
						newboard.black &^= tobit   // remove target
						newboard.white |= tobit    // place piece
					} else {
						if newboard.white&tobit == 0 { // not a capture: increment halfmoves
							newboard.half++
						}
						newboard.black &^= frombit // remove piece
						newboard.white &^= tobit   // remove target
						newboard.black |= tobit    // place piece
					}

					newboard.total++

					moves = append(moves, newboard)
				}
			}
		}
	}

	pawns := b.pawns & colour
	// ignore ranks 1 and 8, pawns can't ever occupy them
	for from = 8; from < 56; from++ {
		frombit = 1 << from
		if pawns&frombit != 0 { // is there a pawn on this square?
			var pawnmoves uint64

			// blockdouble: Find all pieces occupying squares in rank 3 or 7;
			// these pieces would block a double move for white or black
			// respectively. Shift this 8 bits to the left (for white) or right
			// (for black) to get the squares we're blocking a double move to.
			// Remove those squares from candidate moves.
			if tomove == White {
				blockdouble := (occupied & 0x0000000000FF0000) << 8
				pawnmoves = WhitePawnMoves(from) &^ occupied &^ blockdouble
			} else {
				blockdouble := (occupied & 0x0000FF0000000000) >> 8
				pawnmoves = BlackPawnMoves(from) &^ occupied &^ blockdouble
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
						newboard.white &^= frombit // remove piece
						newboard.black &^= tobit   // remove target
						newboard.white |= tobit    // place piece
					} else {
						newboard.black &^= frombit // remove piece
						newboard.white &^= tobit   // remove target
						newboard.black |= tobit    // place piece
					}

					// TODO: pawn captures
					// TODO: en passant captures
					// TODO: set en passant target on double pawn moves
					// TODO: pawn promotion

					// pawn moves don't increment the half move clock
					newboard.total++

					moves = append(moves, newboard)
				}
			}
		}
	}

	// TODO: disallow moves placing the king in check

	return moves
}
