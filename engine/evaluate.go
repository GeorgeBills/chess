package engine

import (
	"math/bits"
)

// https://www.chessprogramming.org/Evaluation

// We use a very simple evaluation function (barely a step up from a simple
// count of material), the logic for which is taken straight from
// https://www.chessprogramming.org/Simplified_Evaluation_Function.

const (
	valPawn   = 100
	valKnight = 320
	valBishop = 330
	valRook   = 500
	valQueen  = 900
	valKing   = 20_000
)

// Piece-Square tables. These give a value bonus or penalty to pieces based on
// the square they are occupying.
// https://www.chessprogramming.org/Piece-Square_Tables.
//
// Note that these are visually from blacks perspective: the 0th index
// corresponds to A1, and the 63rd index corresponds to H8.
var (
	// pawns get bonuses for either advancing or sheltering the king
	pstWhitePawn = [64]int16{
		000, 000, 000, 000, 000, 000, 000, 000,
		005, 010, 010, -20, -20, 010, 010, 005,
		005, -05, -10, 000, 000, -10, -05, 005,
		000, 000, 000, 020, 020, 000, 000, 000,
		005, 005, 010, 025, 025, 010, 005, 005,
		010, 010, 020, 030, 030, 020, 010, 010,
		050, 050, 050, 050, 050, 050, 050, 050,
		000, 000, 000, 000, 000, 000, 000, 000,
	}

	// knights get bonuses for occupying the center, penalties for the edges
	pstWhiteKnight = [64]int16{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 005, 000, 000, 000, 000, 005, -10,
		-10, 010, 010, 010, 010, 010, 010, -10,
		-10, 000, 010, 010, 010, 010, 000, -10,
		-10, 005, 005, 010, 010, 005, 005, -10,
		-10, 000, 005, 010, 010, 005, 000, -10,
		-10, 000, 000, 000, 000, 000, 000, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	}

	// bishops get bonuses for occupying the center, penalties for the edges
	pstWhiteBishop = [64]int16{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 005, 000, 000, 000, 000, 005, -10,
		-10, 010, 010, 010, 010, 010, 010, -10,
		-10, 000, 010, 010, 010, 010, 000, -10,
		-10, 005, 005, 010, 010, 005, 005, -10,
		-10, 000, 005, 010, 010, 005, 000, -10,
		-10, 000, 000, 000, 000, 000, 000, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	}

	// rooks get bonuses for occupying the opponents pawn rank, and for castling
	// rooks get penalties for occupying the A and H files "in order not to defend pawn b3 from a3"
	pstWhiteRook = [64]int16{
		000, 000, 000, 005, 005, 000, 000, 000,
		-05, 000, 000, 000, 000, 000, 000, -05,
		-05, 000, 000, 000, 000, 000, 000, -05,
		-05, 000, 000, 000, 000, 000, 000, -05,
		-05, 000, 000, 000, 000, 000, 000, -05,
		-05, 000, 000, 000, 000, 000, 000, -05,
		005, 010, 010, 010, 010, 010, 010, 005,
		000, 000, 000, 000, 000, 000, 000, 000,
	}

	// queens get bonuses for occupying the center, penalties for the edges
	pstWhiteQueen = [64]int16{
		-20, -10, -10, -05, -05, -10, -10, -20,
		-10, 000, 005, 000, 000, 000, 000, -10,
		-10, 005, 005, 005, 005, 005, 000, -10,
		-05, 000, 005, 005, 005, 005, 000, -05,
		-05, 000, 005, 005, 005, 005, 000, -05,
		-10, 000, 005, 005, 005, 005, 000, -10,
		-10, 000, 000, 000, 000, 000, 000, -10,
		-20, -10, -10, -05, -05, -10, -10, -20,
	}

	// opening-mid game:
	// kings get bonuses for castling to behind the pawn shelter
	// kings get penalties for occupying anywhere unsafe
	pstWhiteKingOpeningMid = [64]int16{
		020, 030, 010, 000, 000, 010, 030, 020,
		020, 020, 000, 000, 000, 000, 020, 020,
		-10, -20, -20, -20, -20, -20, -20, -10,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
	}

	// end game:
	// kings get bonuses for occupying the center
	// kings get penalties for occupying the back ranks (due to back-rank mates)
	pstWhiteKingEnd = [64]int16{
		-50, -30, -30, -30, -30, -30, -30, -50,
		-30, -30, 000, 000, 000, 000, -30, -30,
		-30, -10, 020, 030, 030, 020, -10, -30,
		-30, -10, 030, 040, 040, 030, -10, -30,
		-30, -10, 030, 040, 040, 030, -10, -30,
		-30, -10, 020, 030, 030, 020, -10, -30,
		-30, -20, -10, 000, 000, -10, -20, -30,
		-50, -40, -30, -20, -20, -30, -40, -50,
	}

	pstBlackPawn           [64]int16
	pstBlackKnight         [64]int16
	pstBlackBishop         [64]int16
	pstBlackRook           [64]int16
	pstBlackQueen          [64]int16
	pstBlackKingOpeningMid [64]int16
	pstBlackKingEnd        [64]int16
)

func init() {
	// mirror white piece-square tables for black
	var i uint8 = 0
	var j uint8 = 63
	for i < 64 {
		pstBlackPawn[j] = pstWhitePawn[i]
		pstBlackKnight[j] = pstWhiteKnight[i]
		pstBlackBishop[j] = pstWhiteBishop[i]
		pstBlackRook[j] = pstWhiteRook[i]
		pstBlackQueen[j] = pstWhiteQueen[i]
		pstBlackKingOpeningMid[j] = pstWhiteKingOpeningMid[i]
		pstBlackKingEnd[j] = pstWhiteKingEnd[i]
		i++
		j--
	}
}

// Evaluate returns a score evaluation for the board. A positive number
// indicates that white is winning, a negative number indicates that black is
// winning. Larger numbers indicate that the side is winning by a wider margin
// than lower numbers.
func (b *Board) Evaluate() int16 {
	minorPieces := b.bishops | b.knights

	numWhiteMinor := bits.OnesCount64(b.white & minorPieces)
	numWhiteQueens := bits.OnesCount64(b.white & b.queens)
	numWhiteRooks := bits.OnesCount64(b.white & b.rooks)
	isLateGameWhite := numWhiteQueens == 0 || numWhiteRooks == 0 && numWhiteMinor <= 1

	numBlackMinor := bits.OnesCount64(b.black & minorPieces)
	numBlackQueens := bits.OnesCount64(b.black & b.queens)
	numBlackRooks := bits.OnesCount64(b.black & b.rooks)
	isLateGameBlack := numBlackQueens == 0 || numBlackRooks == 0 && numBlackMinor <= 1

	// late game: both sides have no queen; or queen + knight; or queen + bishop
	isLateGame := isLateGameWhite && isLateGameBlack

	pstWhiteKing := pstWhiteKingOpeningMid
	pstBlackKing := pstBlackKingOpeningMid
	if isLateGame {
		pstWhiteKing = pstWhiteKingEnd
		pstBlackKing = pstBlackKingEnd
	}

	var score int16
	score += evaluateMaterial(b.white&b.pawns, valPawn, &pstWhitePawn)
	score += evaluateMaterial(b.white&b.bishops, valBishop, &pstWhiteBishop)
	score += evaluateMaterial(b.white&b.knights, valKnight, &pstWhiteKnight)
	score += evaluateMaterial(b.white&b.queens, valQueen, &pstWhiteQueen)
	score += evaluateMaterial(b.white&b.kings, valKing, &pstWhiteKing)
	score += evaluateMaterial(b.white&b.rooks, valRook, &pstWhiteRook)
	score -= evaluateMaterial(b.black&b.pawns, valPawn, &pstBlackPawn)
	score -= evaluateMaterial(b.black&b.bishops, valBishop, &pstBlackBishop)
	score -= evaluateMaterial(b.black&b.knights, valKnight, &pstBlackKnight)
	score -= evaluateMaterial(b.black&b.queens, valQueen, &pstBlackQueen)
	score -= evaluateMaterial(b.black&b.kings, valKing, &pstBlackKing)
	score -= evaluateMaterial(b.black&b.rooks, valRook, &pstBlackRook)
	return score
}

func evaluateMaterial(material uint64, val int16, pst *[64]int16) int16 {
	var score int16
	for material != 0 {
		idx, _ := popLSB(&material)
		score += val
		score += pst[idx]
	}
	return score
}
