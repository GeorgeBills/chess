package engine

import (
	"fmt"
	"math"
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
	maskAll  = 0xFFFFFFFFFFFFFFFF
	maskNone = 0x0000000000000000

	rank2 = 1
	rank7 = 6

	maskRank8 uint64 = 0xFF00000000000000
	maskRank7 uint64 = 0x00FF000000000000
	maskRank6 uint64 = 0x0000FF0000000000
	maskRank5 uint64 = 0x000000FF00000000
	maskRank4 uint64 = 0x00000000FF000000
	maskRank3 uint64 = 0x0000000000FF0000
	maskRank2 uint64 = 0x000000000000FF00
	maskRank1 uint64 = 0x00000000000000FF

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
//
// This function will panic if run with certain invalid boards, e.g. if there
// are more than two pieces giving check, or if one side doesn't have a king on
// the board. You should wrap it in a recover, or ideally ensure that you're
// only calling it on valid boards.
func (b Board) Moves(moves []Move) []Move {
	moves = moves[:0] // empty passed in slice

	// checkers is a mask for pieces giving check. if there is more than one bit
	// set then we're in double check.
	var checkers uint64

	// pinned variables are for pieces that must stay on the respective ray.
	var pinnedDiagonalSWNE, pinnedDiagonalNWSE, pinnedVertical, pinnedHorizontal uint64

	var threatRay uint64

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
		pawnsdbl = pawnssgl & maskRank2 &^ ((occupied & maskRank4) >> 16)
		pawnspromo = pawns & maskRank7

		pawnscaptureEast = pawns & (opposing >> 9) &^ maskFileH // ne
		pawnscaptureWest = pawns & (opposing >> 7) &^ maskFileA // nw
	} else {
		colour = b.black
		opposing = b.white
		pawns = b.pawns & b.black
		pawnssgl = pawns &^ (occupied << 8)
		pawnsdbl = pawnssgl & maskRank7 &^ ((occupied & maskRank5) << 16)
		pawnspromo = pawns & maskRank2

		pawnscaptureEast = pawns & (opposing << 9) &^ maskFileA // se
		pawnscaptureWest = pawns & (opposing << 7) &^ maskFileH // sw
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

	rayEvaluateCheckPinForward := func(moves *[64]uint64) uint64 {
		ray := moves[from]
		intersection := ray & occupied
		blockFirst := uint8(bits.TrailingZeros64(intersection))
		var blockFirstBit uint64 = 1 << blockFirst
		blockSecond := uint8(bits.TrailingZeros64(intersection &^ blockFirstBit))
		var blockSecondBit uint64 = 1 << blockSecond

		// occlude the ray based on whether it hits the king or not
		switch {
		case blockFirstBit&king == 0 && blockFirst != 64: // NOT the king
			ray ^= moves[blockFirst]
			// case blockFirstBit&king != 0 && blockSecond != 64:
			// 	// pierce "through" the king
			// 	ray ^= moves[blockSecond]
		}

		threatened |= ray

		// set check and pinned appropriately
		switch {
		case blockFirstBit&king != 0: // king is in check
			checkers |= frombit
			threatRay = ray &^ moves[blockFirst] // occlude threat by king
		case blockSecondBit&king != 0: // piece is pinned
			return blockFirstBit
		}

		return 0
	}

	rayEvaluateCheckPinBackward := func(moves *[64]uint64) uint64 {
		ray := moves[from]
		intersection := ray & occupied
		blockFirst := uint8(63 - bits.LeadingZeros64(intersection))
		var blockFirstBit uint64 = 1 << blockFirst
		blockSecond := uint8(63 - bits.LeadingZeros64(intersection&^blockFirstBit))
		var blockSecondBit uint64 = 1 << blockSecond

		// occlude the ray based on whether it hits the king or not
		switch {
		case blockFirstBit&king == 0 && blockFirst != math.MaxUint8: // NOT the king
			ray ^= moves[blockFirst]
			// case blockFirstBit&king != 0 && blockSecond != math.MaxUint8:
			// 	// pierce "through" the king
			// 	ray ^= moves[blockSecond]
		}

		threatened |= ray

		// set check and pinned appropriately
		switch {
		case blockFirstBit&king != 0: // king is in check
			checkers |= frombit
			threatRay = ray &^ moves[blockFirst] // occlude threat by king
		case blockSecondBit&king != 0: // piece is pinned
			return blockFirstBit
		}

		return 0
	}

	{
		opposingpawns := b.pawns & opposing
		// the king can only be under threat from one pawn at a time; it's not
		// possible to simultaneously move two pawns within capture range, nor
		// is it legal for a king to move to where a pawn could capture it.

		if tomove == White {
			m := (opposingpawns&^maskFileA)>>9 | // sw
				(opposingpawns&^maskFileH)>>7 // se
			if m&king != 0 {
				// flip the check to find out which pawn is checking us
				if frombit = king << 9; frombit&opposingpawns&^maskFileH != 0 {
					checkers |= frombit
				} else {
					checkers |= king << 7 // must be nw
				}
			}
			threatened |= m
		} else {
			m := (opposingpawns&^maskFileH)<<9 | // ne
				(opposingpawns&^maskFileA)<<7 // nw
			if m&king != 0 {
				// flip the check to find out which pawn is checking us
				if frombit = king >> 9; frombit&opposingpawns&^maskFileA != 0 {
					checkers |= frombit
				} else {
					checkers |= king >> 7 // must be sw
				}
			}
			threatened |= m
		}
	}

	{
		opposingknights := b.knights & opposing
		for opposingknights != 0 {
			from = uint8(bits.TrailingZeros64(opposingknights))
			frombit = 1 << from
			opposingknights ^= (1 << from) // unset bit
			knightMoves := movesKnights[from]
			if knightMoves&king != 0 {
				checkers |= frombit
			}
			threatened |= knightMoves
		}
	}

	{
		opposingking := b.kings & opposing // always exactly one king
		from = uint8(bits.TrailingZeros64(opposingking))
		// note that there's a subtle bug here if we start using threatened for more
		// than just a "can our king move to this square?" check - not all of these
		// moves will be legal for the opposing king to make.
		// note that kings can never check another king
		threatened |= movesKing[from]
	}

	{
		opposingrooks := (b.rooks | b.queens) & opposing
		for opposingrooks != 0 {
			from = uint8(bits.TrailingZeros64(opposingrooks))
			frombit = 1 << from
			opposingrooks ^= frombit // unset bit
			pinnedVertical |= rayEvaluateCheckPinForward(&movesNorth)
			pinnedHorizontal |= rayEvaluateCheckPinForward(&movesEast)
			pinnedVertical |= rayEvaluateCheckPinBackward(&movesSouth)
			pinnedHorizontal |= rayEvaluateCheckPinBackward(&movesWest)
		}
	}

	{
		opposingbishops := (b.bishops | b.queens) & opposing
		for opposingbishops != 0 {
			from = uint8(bits.TrailingZeros64(opposingbishops))
			frombit = 1 << from
			opposingbishops ^= frombit // unset bit

			pinnedDiagonalSWNE |= rayEvaluateCheckPinForward(&movesNorthEast)
			pinnedDiagonalNWSE |= rayEvaluateCheckPinBackward(&movesSouthEast)
			pinnedDiagonalSWNE |= rayEvaluateCheckPinBackward(&movesSouthWest)
			pinnedDiagonalNWSE |= rayEvaluateCheckPinForward(&movesNorthWest)
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
	if ep != math.MaxUint8 {
		// ep records the square behind, so we check the squares to the ne and
		// nw (for black) or se and sw (for white) to find pawns adjacent.
		// FIXME: check maymoveto here
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

	rayForward := func(moves *[64]uint64, from uint8) uint64 {
		ray := moves[from]
		intersection := ray & occupied
		if intersection != 0 {
			ray ^= moves[uint8(bits.TrailingZeros64(intersection))]
		}
		return ray
	}

	rayBackward := func(moves *[64]uint64, from uint8) uint64 {
		ray := moves[from]
		intersection := ray & occupied
		if intersection != 0 {
			ray ^= moves[uint8(63-bits.LeadingZeros64(intersection))]
		}
		return ray
	}

	var maymoveto uint64
	switch bits.OnesCount64(checkers) {
	case 0:
		// no check:
		// most pieces can move anywhere
		// pinned pieces may only move on their pinned ray
		// king can only move to unthreatened squares
		maymoveto = maskAll // all squares
	case 1:
		// single check; we must either:
		// capture the piece giving check
		// move a piece on to the threat ray
		// move our king out of threat
		maymoveto = checkers | threatRay
	case 2:
		// double check: we must move our king
		goto KING_MOVES
	default:
		panic(fmt.Sprintf("invalid checkers mask: %b; %#v", checkers, b))
	}

	{
		ppcopy := pawnspromo
		for ppcopy != 0 {
			from = uint8(bits.TrailingZeros64(ppcopy))
			ppcopy ^= (1 << from) // unset bit
			// we could calculate masks for the below checks (e.g. pawns that
			// can promote by pushing are just pawns that can push bitwise
			// AND'ed with pawns that can promote), but generating that mask
			// every time we generate moves doesn't pay off when pawns being in
			// position to promote is so rare.
			// TODO: disjoint masks for promo and not promo, combined for captures and pushes
			// TODO: can mask pawn moves (single, double, capture) on maymoveto en masse
			// TODO: need to add tests for pawn promos capturing a checker
			if tomove == White {
				ne := from + 9
				if tobit = 1 << ne; opposing&maymoveto&tobit&^maskFileH != 0 {
					addpromos(from, ne, true)
				}
				nw := from + 7
				if tobit = 1 << nw; opposing&maymoveto&tobit&^maskFileA != 0 {
					addpromos(from, nw, true)
				}
				push := from + 8
				if tobit = 1 << push; maymoveto&tobit&(^occupied) != 0 {
					addpromos(from, push, false)
				}
			} else {
				se := from - 7
				if tobit = 1 << se; opposing&maymoveto&tobit&^maskFileA != 0 {
					addpromos(from, se, true)
				}
				sw := from - 9
				if tobit = 1 << sw; opposing&maymoveto&tobit&^maskFileH != 0 {
					addpromos(from, sw, true)
				}
				push := from - 8
				if tobit = 1 << push; maymoveto&tobit&(^occupied) != 0 {
					addpromos(from, push, false)
				}
			}
		}
	}

	{
		// TODO: limit pawnsnotpromo when generating
		pawnsnotpromo := pawns &^ pawnspromo // need to set this before we start unsetting bits in pawnspromo
		for pawnsnotpromo != 0 {
			from = uint8(bits.TrailingZeros64(pawnsnotpromo))
			frombit = 1 << from
			pawnsnotpromo ^= frombit // unset bit
			// TODO: set en passant target on double pawn moves
			if tomove == White {
				if pawnsdbl&maymoveto&frombit != 0 {
					moves = append(moves, NewMove(from, from+16))
				}
				if pawnssgl&maymoveto&frombit != 0 {
					moves = append(moves, NewMove(from, from+8))
				}
				if pawnscaptureEast&(maymoveto>>9)&frombit != 0 {
					moves = append(moves, NewCapture(from, from+9))
				}
				if pawnscaptureWest&(maymoveto>>7)&frombit != 0 {
					moves = append(moves, NewCapture(from, from+7))
				}
			} else {
				if pawnsdbl&maymoveto&frombit != 0 {
					moves = append(moves, NewMove(from, from-16))
				}
				if pawnssgl&maymoveto&frombit != 0 {
					moves = append(moves, NewMove(from, from-8))
				}
				if pawnscaptureEast&(maymoveto<<9)&frombit != 0 {
					moves = append(moves, NewCapture(from, from-9))
				}
				if pawnscaptureWest&(maymoveto<<7)&frombit != 0 {
					moves = append(moves, NewCapture(from, from-7))
				}
			}
		}
	}

	{
		knights := b.knights & colour
		for knights != 0 {
			from = uint8(bits.TrailingZeros64(knights))
			knights ^= (1 << from) // unset bit
			m := movesKnights[from] &^ colour & maymoveto
			addCaptures(from, m&opposing)
			addQuietMoves(from, m&^occupied)
		}
	}

	{
		rooks := (b.rooks | b.queens) & colour
		for rooks != 0 {
			from = uint8(bits.TrailingZeros64(rooks))
			frombit = 1 << from
			rooks ^= frombit // unset bit
			var movesqs uint64
			if pinnedHorizontal&frombit == 0 { // not pinned horizontally, can move vertically
				movesqs |= rayForward(&movesNorth, from)
				movesqs |= rayBackward(&movesSouth, from)
			}
			if pinnedVertical&frombit == 0 { // not pinned vertically, can move horizontally
				movesqs |= rayForward(&movesEast, from)
				movesqs |= rayBackward(&movesWest, from)
			}
			movesqs &= maymoveto
			addCaptures(from, movesqs&opposing)
			addQuietMoves(from, movesqs&^occupied)
		}
	}

	{
		bishops := (b.bishops | b.queens) & colour
		for bishops != 0 {
			from = uint8(bits.TrailingZeros64(bishops))
			frombit = 1 << from
			bishops ^= frombit // unset bit
			var movesqs uint64
			if pinnedDiagonalSWNE&frombit == 0 { // not pinned to the SW/NE diagonal, can move on the NW/SE diagonal
				movesqs |= rayForward(&movesNorthWest, from)
				movesqs |= rayBackward(&movesSouthEast, from)
			}
			if pinnedDiagonalNWSE&frombit == 0 { // not pinned to the NW/SE diagonal, can move on the SW/NE diagonal
				movesqs |= rayForward(&movesNorthEast, from)
				movesqs |= rayBackward(&movesSouthWest, from)
			}
			movesqs &= maymoveto
			addCaptures(from, movesqs&opposing)
			addQuietMoves(from, movesqs&^occupied)
		}
	}

KING_MOVES:
	{
		from = uint8(bits.TrailingZeros64(king)) // always exactly one king
		m := movesKing[from] &^ colour &^ threatened
		addCaptures(from, m&opposing)
		addQuietMoves(from, m&^occupied)
	}

	// TODO: disallow moves placing the king in check

	return moves
}
