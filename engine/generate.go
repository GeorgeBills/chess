package engine

// TODO: use init block to pregenerate moves

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

// whitePawnPushes returns the moves a white pawn at index i can make, ignoring
// captures and en passant.
func whitePawnPushes(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 8) // n
	// if a white pawn is on rank 2 it may move two squares
	if Rank(i) == rank2 {
		moves |= 1 << (i + 16) // nn
	}
	return moves
}

// blackPawnPushes returns the moves a black pawn at index i can make, ignoring
// captures and en passant.
func blackPawnPushes(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i - 8) // s
	// if a black pawn is on rank 7 it may move two squares
	if Rank(i) == rank7 {
		moves |= 1 << (i - 16) // ss
	}
	return moves
}

func whitePawnCaptures(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 9) // ne
	moves |= 1 << (i + 7) // nw
	return moves
}

func blackPawnCaptures(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i - 7) // se
	moves |= 1 << (i - 9) // sw
	return moves
}

// kingMoves returns the moves a king at index i can make, ignoring castling.
func kingMoves(i uint8) uint64 {
	var moves uint64
	moves |= 1 << (i + 8) // n
	moves |= 1 << (i - 8) // s
	// can't move east if we're on file h
	if File(i) != fileH {
		moves |= 1 << (i + 1) // e
		moves |= 1 << (i + 9) // ne
		moves |= 1 << (i - 7) // se
	}
	// can't move west if we're on file a
	if File(i) != fileA {
		moves |= 1 << (i - 1) // w
		moves |= 1 << (i + 7) // nw
		moves |= 1 << (i - 9) // sw
	}
	return moves
}

// knightMoves returns the moves a knight at index i can make.
func knightMoves(i uint8) uint64 {
	var moves uint64
	file := File(i)
	if file > fileA {
		moves |= 1 << (i + 15) // nnw (+2×8, -1)
		moves |= 1 << (i - 17) // ssw (-2×8, -1)
		if file > fileB {
			moves |= 1 << (i + 6)  // wwn (+8, -2×1)
			moves |= 1 << (i - 10) // wws (-8, -2×1)
		}
	}
	if file < fileH {
		moves |= 1 << (i + 17) // nne (+2×8, +1)
		moves |= 1 << (i - 15) // sse (-2×8, +1)
		if file < fileG {
			moves |= 1 << (i + 10) // een (+8, +2×1)
			moves |= 1 << (i - 6)  // ees (-8, +2×1)
		}
	}
	return moves
}

const (
	checkNone   = 0
	checkDouble = 0xFFFFFFFF

	rank2 = 1
	rank7 = 6

	rank7mask uint64 = 0x00FF000000000000
	rank6mask uint64 = 0x0000FF0000000000
	rank5mask uint64 = 0x000000FF00000000
	rank4mask uint64 = 0x00000000FF000000
	rank3mask uint64 = 0x0000000000FF0000
	rank2mask uint64 = 0x000000000000FF00

	fileA = 0
	fileB = 1
	fileG = 6
	fileH = 7

	// TODO: can we use filemasks to speed up file checks?
)

// Moves returns a slice of possible moves from the current board state.
func (b Board) Moves() []Move {
	moves := make([]Move, 0, 32)

	var from, to uint8
	var frombit, tobit uint64
	var colour, opposing uint64
	var pawnsdbl, pawnspromo uint64

	occupied := b.white | b.black

	tomove := b.ToMove()
	var pawns uint64
	// TODO: use pointers so we don't need to check tomove later
	if tomove == White {
		colour = b.white
		opposing = b.black
		pawns = b.pawns & b.white
		pawnsdbl = pawns & rank2mask &^ ((occupied & rank3mask) >> 8) &^ ((occupied & rank4mask) >> 16)
		pawnspromo = pawns & rank7mask
	} else {
		colour = b.black
		opposing = b.white
		pawns = b.pawns & b.black
		pawnsdbl = pawns & rank7mask &^ ((occupied & rank6mask) << 8) &^ ((occupied & rank5mask) << 16)
		pawnspromo = pawns & rank2mask
	}

	// Loop over opposing pieces to see if we're in check, and to mark both
	// threatened squares and pinned pieces.
	//
	// Here we only care about captures or the threat of capture. This means
	// that we can completely ignore pawn pushes and castling.
	//
	// "Scanners" (rooks, bishops and queens) need to scan "through" other
	// pieces to see if those piece/s are blocking check on the king (in which
	// case the piece is pinned and may not move out of the line of threat).
	//
	// We must never move our king to a threatened square, through a threatened
	// square (in the case of castling) or make a move that exposes a check on
	// our king. This includes capturing opposing pieces that are blocking
	// check.
	//
	// If our king is in check, then we must either capture the checking piece,
	// move our king to an unthreatened square, or (in the case of scanners)
	// block the threat with another piece.
	//
	// If our king is in double check - check from two separate pieces - then we
	// must move our king to a safe square. It's neither possible to capture nor
	// to block two separate threatening pieces in the same turn, so the only
	// remaining option is to move our king.

	var threatened uint64 // threatened tracks squares we may not move our king to
	opposingknights := b.knights & opposing
	opposingking := b.kings & opposing
	opposingrooks := b.rooks & opposing
	opposingbishops := b.bishops & opposing
FIND_THREAT:
	for from = 0; from < 64; from++ {
		frombit = 1 << from // TODO: *=2 frombit each round and calc from only when needed?

		if opposing&frombit == 0 {
			continue FIND_THREAT
		}

		if opposingknights&frombit != 0 { // is there a knight on this square?
			threatened |= knightMoves(from)
			continue FIND_THREAT
		}

		if opposingking&frombit != 0 { // is there a king on this square?
			threatened |= kingMoves(from)
			continue FIND_THREAT
		}

		if opposingrooks&frombit != 0 { // is there a rook on this square?
			rank := Rank(from)
			for n := from + 8; n < 64; n += 8 {
				tobit = 1 << n
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for e := from + 1; e < (rank+1)*8; e++ {
				tobit = 1 << e
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for s := from - 8; s > 0 && s < 64; s -= 8 { // uint wraps below 0
				tobit = 1 << s
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for w := from - 1; w > (rank*8)-1; w-- {
				tobit = 1 << w
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
		}

		if opposingbishops&frombit != 0 { // is there a bishop on this square?
			for ne := from + 9; ne < 64 && File(ne) != fileA; ne += 9 {
				tobit = 1 << ne
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for se := from - 7; 0 < se && se < 64 && File(se) != fileA; se -= 7 {
				tobit = 1 << se
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for sw := from - 9; 0 < sw && sw < 64 && File(sw) != fileH; sw -= 9 {
				tobit = 1 << sw
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
			for nw := from + 7; nw < 64 && File(nw) != fileH; nw += 7 {
				tobit = 1 << nw
				threatened |= tobit
				if occupied&tobit != 0 {
					break
				}
			}
		}
	}

	// TODO: never loop over squares for destinations

	// TODO: castling

	ep := b.EnPassant()
	if ep != 0 {
		// ep records the square behind, so we check the squares to the ne and
		// nw (for black) or se and sw (for white) to find pawns adjacent.
		if tomove == White {
			if from = ep - 7; pawns&(1<<from) != 0 { // sw
				moves = append(moves, NewEnPassant(from, from+7)) // ne
			}
			if from = ep - 9; pawns&(1<<from) != 0 { // se
				moves = append(moves, NewEnPassant(from, from+9)) // nw
			}
		} else {
			if from = ep + 7; pawns&(1<<from) != 0 {
				moves = append(moves, NewEnPassant(from, from-7)) // se
			}
			if from = ep + 9; pawns&(1<<from) != 0 {
				moves = append(moves, NewEnPassant(from, from-9)) // sw
			}
		}
	}

	addpromos := func(from, to uint8, capture bool) {
		moves = append(moves, NewQueenPromotion(from, to, capture))
		moves = append(moves, NewKnightPromotion(from, to, capture))
		moves = append(moves, NewRookPromotion(from, to, capture))
		moves = append(moves, NewBishopPromotion(from, to, capture))
	}

	// - Find all pieces of the given colour.
	// - For each square, check if we have a piece there.
	// - If there is a piece there, find all the moves that piece can make.
	// - AND NOT the moves with our colour (we can't move to a square we occupy).
	//
	// We evaluate pieces in descending frequency order (pawn, knight, king) to
	// hopefully skip a loop iteration as early as possible.
	knights := b.knights & colour
	bishops := (b.bishops | b.queens) & colour
	rooks := (b.rooks | b.queens) & colour
	king := b.kings & colour
FIND_MOVES:
	for from = 0; from < 64; from++ {
		frombit = 1 << from

		if colour&frombit == 0 {
			continue FIND_MOVES
		}

		if pawnspromo&frombit != 0 { // is there a pawn that can promote on this square?
			if tomove == White {
				if ne := from + 9; opposing&(1<<ne) != 0 {
					addpromos(from, ne, true)
				}
				if nw := from + 7; opposing&(1<<nw) != 0 {
					addpromos(from, nw, true)
				}
				if push := from + 8; occupied&(1<<push) == 0 {
					addpromos(from, push, false)
				}
			} else {
				if se := from - 7; opposing&(1<<se) != 0 {
					addpromos(from, se, true)
				}
				if sw := from - 9; opposing&(1<<sw) != 0 {
					addpromos(from, sw, true)
				}
				if push := from - 8; occupied&(1<<push) == 0 {
					addpromos(from, push, false)
				}
			}
			continue FIND_MOVES
		}

		if pawns&frombit != 0 { // is there a pawn on this square?
			// TODO: set en passant target on double pawn moves
			if tomove == White {
				if pawnsdbl&frombit != 0 {
					moves = append(moves, NewMove(from, from+16))
				}
				if push := from + 8; occupied&(1<<push) == 0 {
					moves = append(moves, NewMove(from, push))
				}
				if ne := from + 9; opposing&(1<<ne) != 0 {
					moves = append(moves, NewCapture(from, ne))
				}
				if nw := from + 7; opposing&(1<<nw) != 0 {
					moves = append(moves, NewCapture(from, nw))
				}
			} else {
				if pawnsdbl&frombit != 0 {
					moves = append(moves, NewMove(from, from-16))
				}
				if push := from - 8; occupied&(1<<push) == 0 {
					moves = append(moves, NewMove(from, push))
				}
				if se := from - 7; opposing&(1<<se) != 0 {
					moves = append(moves, NewCapture(from, se))
				}
				if sw := from - 9; opposing&(1<<sw) != 0 {
					moves = append(moves, NewCapture(from, sw))
				}
			}
			continue FIND_MOVES
		}

		if knights&frombit != 0 { // is there a knight on this square?
			knightMoves := knightMoves(from) &^ colour
			for to = 0; to < 64; to++ {
				tobit = 1 << to
				if knightMoves&tobit != 0 { // is there a move to this square?
					capture := b.black&tobit != 0 || b.white&tobit != 0
					if capture {
						moves = append(moves, NewCapture(from, to))
					} else {
						moves = append(moves, NewMove(from, to))
					}
				}
			}
			continue FIND_MOVES
		}

		if king&frombit != 0 { // is there a king on this square?
			kingMoves := kingMoves(from) &^ colour &^ threatened // king cannot move into check
			for to = 0; to < 64; to++ {
				tobit = 1 << to
				if kingMoves&tobit != 0 { // is there a move to this square?
					capture := b.black&tobit != 0 || b.white&tobit != 0
					if capture {
						moves = append(moves, NewCapture(from, to))
					} else {
						moves = append(moves, NewMove(from, to))
					}
				}
			}
			continue FIND_MOVES
		}

		if rooks&frombit != 0 { // is there a rook on this square?
			rank := Rank(from)
			for n := from + 8; n < 64; n += 8 {
				tobit = 1 << n
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, n))
					}
					break
				}
				moves = append(moves, NewMove(from, n))
			}
			for e := from + 1; e < (rank+1)*8; e++ {
				tobit = 1 << e
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, e))
					}
					break
				}
				moves = append(moves, NewMove(from, e))
			}
			for s := from - 8; s > 0 && s < 64; s -= 8 { // uint wraps below 0
				tobit = 1 << s
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, s))
					}
					break
				}
				moves = append(moves, NewMove(from, s))
			}
			for w := from - 1; w > (rank*8)-1; w-- {
				tobit = 1 << w
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, w))
					}
					break
				}
				moves = append(moves, NewMove(from, w))
			}
		}

		if bishops&frombit != 0 { // is there a bishop on this square?
			for ne := from + 9; ne < 64 && File(ne) != fileA; ne += 9 {
				tobit = 1 << ne
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, ne))
					}
					break
				}
				moves = append(moves, NewMove(from, ne))
			}
			for se := from - 7; 0 < se && se < 64 && File(se) != fileA; se -= 7 {
				tobit = 1 << se
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, se))
					}
					break
				}
				moves = append(moves, NewMove(from, se))
			}
			for sw := from - 9; 0 < sw && sw < 64 && File(sw) != fileH; sw -= 9 {
				tobit = 1 << sw
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, sw))
					}
					break
				}
				moves = append(moves, NewMove(from, sw))
			}
			for nw := from + 7; nw < 64 && File(nw) != fileH; nw += 7 {
				tobit = 1 << nw
				if occupied&tobit != 0 {
					if opposing&tobit != 0 {
						moves = append(moves, NewCapture(from, nw))
					}
					break
				}
				moves = append(moves, NewMove(from, nw))
			}
		}

	}

	// TODO: disallow moves placing the king in check

	return moves
}
