package engine

import (
	"math/bits"
)

// Pregenerated masks for moves in any of the compass directions from any given
// square, and for kings and knights. Takes up 10 * 64 * 64 = 40kb of memory,
// which should fit in L1 cache on a modern CPU.
var (
	movesNorth     [64]uint64
	movesNorthEast [64]uint64
	movesEast      [64]uint64
	movesSouthEast [64]uint64
	movesSouth     [64]uint64
	movesSouthWest [64]uint64
	movesWest      [64]uint64
	movesNorthWest [64]uint64
	movesKing      [64]uint64
	movesKnights   [64]uint64
)

func init() {
	// n: north; e: east; s: south; w: west nn, ss: north north and south south
	// (used by pawns performing double moves) nne, een, ssw: north north east,
	// etc (used by knights)
	//
	// North and South are easy; add or subtract 8 (a full rank). With bit
	// shifting you don't even need to check if that's off the board: shifting
	// above 64 will shift the bits out, and subtracting past 0 with an unsigned
	// int ends up wrapping into a very large shift up and off the board again.
	//
	// East and West require checking if you'll leave the board or not, since
	// leaving the board in either of those directions will wrap around; e.g. a
	// king shouldn't be able to move one square West from A2 (8) and end up on
	// H1 (7).

	var from uint8

	for from = 0; from < 64; from++ {
		rank := Rank(from)
		file := File(from)

		// horizontal: rooks, queens
		for n := from + 8; n < 64; n += 8 {
			movesNorth[from] |= 1 << n
		}
		for e := from + 1; e < (rank+1)*8; e++ {
			movesEast[from] |= 1 << e
		}
		for s := from - 8; s < 64; s -= 8 {
			movesSouth[from] |= 1 << s
		}
		for w := from - 1; w != (rank*8)-1; w-- {
			movesWest[from] |= 1 << w
		}

		// diagonal: bishops, queens
		for ne := from + 9; ne < 64 && File(ne) != fileA; ne += 9 {
			movesNorthEast[from] |= 1 << ne
		}
		for se := from - 7; se < 64 && File(se) != fileA; se -= 7 {
			movesSouthEast[from] |= 1 << se
		}
		for sw := from - 9; sw < 64 && File(sw) != fileH; sw -= 9 {
			movesSouthWest[from] |= 1 << sw
		}
		for nw := from + 7; nw < 64 && File(nw) != fileH; nw += 7 {
			movesNorthWest[from] |= 1 << nw
		}

		// king
		movesKing[from] |= 1 << (from + 8) // n
		movesKing[from] |= 1 << (from - 8) // s
		// can't move east if we're on file h
		if file != fileH {
			movesKing[from] |= 1 << (from + 1) // e
			movesKing[from] |= 1 << (from + 9) // ne
			movesKing[from] |= 1 << (from - 7) // se
		}
		// can't move west if we're on file a
		if file != fileA {
			movesKing[from] |= 1 << (from - 1) // w
			movesKing[from] |= 1 << (from + 7) // nw
			movesKing[from] |= 1 << (from - 9) // sw
		}

		// knights
		if file > fileA {
			movesKnights[from] |= 1 << (from + 15) // nnw (+2×8, -1)
			movesKnights[from] |= 1 << (from - 17) // ssw (-2×8, -1)
			if file > fileB {
				movesKnights[from] |= 1 << (from + 6)  // wwn (+8, -2×1)
				movesKnights[from] |= 1 << (from - 10) // wws (-8, -2×1)
			}
		}
		if file < fileH {
			movesKnights[from] |= 1 << (from + 17) // nne (+2×8, +1)
			movesKnights[from] |= 1 << (from - 15) // sse (-2×8, +1)
			if file < fileG {
				movesKnights[from] |= 1 << (from + 10) // een (+8, +2×1)
				movesKnights[from] |= 1 << (from - 6)  // ees (-8, +2×1)
			}
		}
	}
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

	maskFileA = 1<<A1 | 1<<A2 | 1<<A3 | 1<<A4 | 1<<A5 | 1<<A6 | 1<<A7 | 1<<A8
	maskFileH = 1<<H1 | 1<<H2 | 1<<H3 | 1<<H4 | 1<<H5 | 1<<H6 | 1<<H7 | 1<<H8
)

// "One may not castle out of, through, or into check".
//
// The castle threat mask constants track the squares for the kings starting
// position, the square the king will move through, and the kings final
// position. If one of these squares is under threat then the relevant castling
// is not legal.
const (
	whiteKingsideCastleThreatMask  uint64 = 1<<E1 | 1<<F1 | 1<<G1
	whiteQueensideCastleThreatMask uint64 = 1<<C1 | 1<<D1 | 1<<E1
	blackKingsideCastleThreatMask  uint64 = 1<<E8 | 1<<F8 | 1<<G8
	blackQueensideCastleThreatMask uint64 = 1<<C8 | 1<<D8 | 1<<E8
)

// In order to castle there must be no spaces in between the king and the rook.
// The castle block mask constants track the squares in between the king and the
// rook. If one of these squares is occupied then the relevant castling is not
// legal.
const (
	whiteKingsideCastleBlockMask  uint64 = 1<<F1 | 1<<G1
	whiteQueensideCastleBlockMask uint64 = 1<<B1 | 1<<C1 | 1<<D1
	blackKingsideCastleBlockMask  uint64 = 1<<F8 | 1<<G8
	blackQueensideCastleBlockMask uint64 = 1<<B8 | 1<<C8 | 1<<D8
)

// Moves returns a slice of possible moves from the current board state.
func (b Board) Moves(moves []Move) []Move {
	moves = moves[:0] // empty passed in slice

	// checkers is a mask for pieces giving check. if there is more than one bit
	// set then we're in double check.
	// var checkers uint64

	var from uint8
	var frombit, tobit uint64
	var colour, opposing uint64
	var pawnssgl, pawnsdbl, pawnspromo, pawnscaptureEast, pawnscaptureWest uint64

	occupied := b.white | b.black

	tomove := b.ToMove()
	var pawns uint64
	// TODO: use pointers so we don't need to check tomove later
	if tomove == White {
		colour = b.white
		opposing = b.black
		pawns = b.pawns & b.white
		pawnssgl = pawns &^ (occupied >> 8)
		pawnsdbl = pawnssgl & rank2mask &^ ((occupied & rank4mask) >> 16)
		pawnspromo = pawns & rank7mask

		pawnscaptureEast = pawns & (opposing >> 9) &^ maskFileH // ne
		pawnscaptureWest = pawns & (opposing >> 7) &^ maskFileA // nw
	} else {
		colour = b.black
		opposing = b.white
		pawns = b.pawns & b.black
		pawnssgl = pawns &^ (occupied << 8)
		pawnsdbl = pawnssgl & rank7mask &^ ((occupied & rank5mask) << 16)
		pawnspromo = pawns & rank2mask

		pawnscaptureEast = pawns & (opposing << 9) &^ maskFileH // se
		pawnscaptureWest = pawns & (opposing << 7) &^ maskFileA // sw
	}

	king := b.kings & colour

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

	// TODO: "covered" seems to be the appropriate chess term
	var threatened uint64 // threatened tracks squares we may not move our king to
	{
		opposingknights := b.knights & opposing
		opposingking := b.kings & opposing
		opposingrooks := (b.rooks | b.queens) & opposing
		opposingbishops := (b.bishops | b.queens) & opposing
	FIND_THREAT:
		for from = 0; from < 64; from++ {
			frombit = 1 << from // TODO: *=2 frombit each round and calc from only when needed?

			if opposing&frombit == 0 {
				continue FIND_THREAT
			}

			if opposingknights&frombit != 0 { // is there a knight on this square?
				threatened |= movesKnights[from]
				// if threatened&king != 0 {
				// 	check++
				// 	// TODO: now need to mark the checking piece; reverse the moves from king
				// }
				continue FIND_THREAT
			}

			if opposingking&frombit != 0 { // is there a king on this square?
				// note that there's a subtle bug here if we start using
				// threatened for more than just a "can our king move to this
				// square?" check - not all of these moves will be legal for the
				// opposing king to make.
				// note that kings can never check another king
				threatened |= movesKing[from]
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
				for s := from - 8; s < 64; s -= 8 { // uint wraps below 0
					tobit = 1 << s
					threatened |= tobit
					if occupied&tobit != 0 {
						break
					}
				}
				for w := from - 1; w != (rank*8)-1; w-- {
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
				for se := from - 7; se < 64 && File(se) != fileA; se -= 7 {
					tobit = 1 << se
					threatened |= tobit
					// if king&tobit != 0 {
					// 	// king in threat
					// 	checkers |= frombit
					// 	break
					// }
					if occupied&tobit != 0 {
						break
					}
				}
				for sw := from - 9; sw < 64 && File(sw) != fileH; sw -= 9 {
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
	}

	// Check for castling.
	if tomove == White {
		if b.CanWhiteCastleKingSide() && threatened&whiteKingsideCastleThreatMask == 0 && occupied&whiteKingsideCastleBlockMask == 0 {
			moves = append(moves, NewWhiteKingsideCastle())
		}
		if b.CanWhiteCastleQueenSide() && threatened&whiteQueensideCastleThreatMask == 0 && occupied&whiteQueensideCastleBlockMask == 0 {
			moves = append(moves, NewWhiteQueensideCastle())
		}
	} else {
		if b.CanBlackCastleKingSide() && threatened&blackKingsideCastleThreatMask == 0 && occupied&blackKingsideCastleBlockMask == 0 {
			moves = append(moves, NewBlackKingsideCastle())
		}
		if b.CanBlackCastleQueenSide() && threatened&blackQueensideCastleThreatMask == 0 && occupied&blackQueensideCastleBlockMask == 0 {
			moves = append(moves, NewBlackQueensideCastle())
		}
	}

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

	maybeCapture := func(from, to uint8, tobit uint64) {
		if opposing&tobit != 0 {
			moves = append(moves, NewCapture(from, to))
		}
	}

	addQuietMoves := func(from uint8, quiet uint64) {
		for quiet != 0 {
			to := uint8(bits.TrailingZeros64(quiet))
			quiet ^= (1 << to) // unset bit
			moves = append(moves, NewMove(from, to))
		}
	}

	addCaptures := func(from uint8, captures uint64) {
		for captures != 0 {
			to := uint8(bits.TrailingZeros64(captures))
			captures ^= (1 << to) // unset bit
			moves = append(moves, NewCapture(from, to))
		}
	}

	// - Find all pieces of the given colour.
	// - For each square, check if we have a piece there.
	// - If there is a piece there, find all the moves that piece can make.
	// - AND NOT the moves with our colour (we can't move to a square we occupy).
	//
	// We evaluate pieces in descending frequency order (pawn, knight, king) to
	// hopefully skip a loop iteration as early as possible.
	{
		// var maymoveto uint64
		// switch bits.OnesCount64(checkers) {
		// case 0:
		// 	// no check:
		// 	// most pieces can move anywhere
		// 	// pinned pieces may only move on their pinned ray
		// 	// king can only move to unthreatened squares
		// 	maymoveto = 0xFFFFFFFF // all squares
		// case 1:
		// 	// single check; we must either:
		// 	// capture the piece giving check
		// 	// move a piece on to the threat ray
		// 	// move our king out of threat
		// 	maymoveto = checkers & threatray
		// case 2:
		// 	// double check: we must move our king
		// 	// TODO: should we just goto the king moves here?
		// 	maymoveto = 0x00000000 // no squares
		// }

		knights := b.knights & colour
		bishops := (b.bishops | b.queens) & colour
		rooks := (b.rooks | b.queens) & colour
	FIND_MOVES:
		for colour != 0 {
			from = uint8(bits.TrailingZeros64(colour))
			frombit = 1 << from
			colour ^= frombit // unset

			if pawnspromo&frombit != 0 { // is there a pawn that can promote on this square?
				// we could calculate masks for the below checks (e.g. pawns
				// that can promote by pushing are just pawns that can push
				// bitwise AND'ed with pawns that can promote), but generating
				// that mask every time we generate moves doesn't pay off when
				// pawns being in position to promote is so rare.
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
					if pawnssgl&frombit != 0 {
						moves = append(moves, NewMove(from, from+8))
					}
					if pawnscaptureEast&frombit != 0 {
						moves = append(moves, NewCapture(from, from+9))
					}
					if pawnscaptureWest&frombit != 0 {
						moves = append(moves, NewCapture(from, from+7))
					}
				} else {
					if pawnsdbl&frombit != 0 {
						moves = append(moves, NewMove(from, from-16))
					}
					if pawnssgl&frombit != 0 {
						moves = append(moves, NewMove(from, from-8))
					}
					if pawnscaptureEast&frombit != 0 {
						moves = append(moves, NewCapture(from, from-9))
					}
					if pawnscaptureWest&frombit != 0 {
						moves = append(moves, NewCapture(from, from-7))
					}
				}
				continue FIND_MOVES
			}

			if knights&frombit != 0 { // is there a knight on this square?
				m := movesKnights[from] &^ colour
				addCaptures(from, m&opposing)
				addQuietMoves(from, m&^occupied)
				continue FIND_MOVES
			}

			if king&frombit != 0 { // is there a king on this square?
				m := movesKing[from] &^ colour &^ threatened
				addCaptures(from, m&opposing)
				addQuietMoves(from, m&^occupied)
				continue FIND_MOVES
			}

			if rooks&frombit != 0 { // is there a rook on this square?
				rank := Rank(from)
				for n := from + 8; n < 64; n += 8 {
					tobit = 1 << n
					if occupied&tobit != 0 {
						maybeCapture(from, n, tobit)
						break
					}
					moves = append(moves, NewMove(from, n))
				}
				for e := from + 1; e < (rank+1)*8; e++ {
					tobit = 1 << e
					if occupied&tobit != 0 {
						maybeCapture(from, e, tobit)
						break
					}
					moves = append(moves, NewMove(from, e))
				}
				for s := from - 8; s < 64; s -= 8 { // uint wraps below 0
					tobit = 1 << s
					if occupied&tobit != 0 {
						maybeCapture(from, s, tobit)
						break
					}
					moves = append(moves, NewMove(from, s))
				}
				for w := from - 1; w != (rank*8)-1; w-- {
					tobit = 1 << w
					if occupied&tobit != 0 {
						maybeCapture(from, w, tobit)
						break
					}
					moves = append(moves, NewMove(from, w))
				}
				// need to fall through here in case we have a queen
			}

			if bishops&frombit != 0 { // is there a bishop on this square?
				var ray, intersect, movesqs uint64

				ray = movesNorthEast[from]
				intersect = ray & occupied
				if intersect != 0 {
					ray ^= movesNorthEast[uint8(bits.TrailingZeros64(intersect))]
				}
				movesqs |= ray

				ray = movesSouthEast[from]
				intersect = ray & occupied
				if intersect != 0 {
					ray ^= movesSouthEast[uint8(63-bits.LeadingZeros64(intersect))]
				}
				movesqs |= ray

				ray = movesSouthWest[from]
				intersect = ray & occupied
				if intersect != 0 {
					ray ^= movesSouthWest[uint8(63-bits.LeadingZeros64(intersect))]
				}
				movesqs |= ray

				ray = movesNorthWest[from]
				intersect = ray & occupied
				if intersect != 0 {
					ray ^= movesNorthWest[uint8(bits.TrailingZeros64(intersect))]
				}
				movesqs |= ray

				addCaptures(from, movesqs&opposing)
				addQuietMoves(from, movesqs&^occupied)
			}
		}
	}

	// TODO: disallow moves placing the king in check

	return moves
}
