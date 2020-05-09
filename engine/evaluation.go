package engine

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
	pstWhiteKing = [64]int16{
		020, 030, 010, 000, 000, 010, 030, 020,
		020, 020, 000, 000, 000, 000, 020, 020,
		-10, -20, -20, -20, -20, -20, -20, -10,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
	}
	pstBlackPawn   [64]int16
	pstBlackKnight [64]int16
	pstBlackBishop [64]int16
	pstBlackRook   [64]int16
	pstBlackQueen  [64]int16
	pstBlackKing   [64]int16
)

func init() {
	var i uint8 = 0
	var j uint8 = 63
	for i < 64 {
		pstBlackPawn[j] = pstWhitePawn[i]
		pstBlackKnight[j] = pstWhiteKnight[i]
		pstBlackBishop[j] = pstWhiteBishop[i]
		pstBlackRook[j] = pstWhiteRook[i]
		pstBlackQueen[j] = pstWhiteQueen[i]
		pstBlackKing[j] = pstWhiteKing[i]
		i++
		j--
	}
}

// Evaluate returns a score evaluation for the board. A positive number
// indicates that white is winning, a negative number indicates that black is
// winning. Larger numbers indicate that the side is winning by a wider margin
// than lower numbers.
func (b *Board) Evaluate() int16 {
	var score int16

	for whitePawns := b.white & b.pawns; whitePawns != 0; {
		idx, _ := popLSB(&whitePawns)
		score += valPawn
		score += pstWhitePawn[idx]
	}

	for whiteBishops := b.white & b.bishops; whiteBishops != 0; {
		idx, _ := popLSB(&whiteBishops)
		score += valBishop
		score += pstWhiteBishop[idx]
	}

	for whiteKnights := b.white & b.knights; whiteKnights != 0; {
		idx, _ := popLSB(&whiteKnights)
		score += valKnight
		score += pstWhiteKnight[idx]
	}

	for whiteQueens := b.white & b.queens; whiteQueens != 0; {
		idx, _ := popLSB(&whiteQueens)
		score += valQueen
		score += pstWhiteQueen[idx]
	}

	for whiteKings := b.white & b.kings; whiteKings != 0; {
		idx, _ := popLSB(&whiteKings)
		score += valKing
		score += pstWhiteKing[idx]
	}

	for whiteRooks := b.white & b.rooks; whiteRooks != 0; {
		idx, _ := popLSB(&whiteRooks)
		score += valRook
		score += pstWhiteRook[idx]
	}

	for blackPawns := b.black & b.pawns; blackPawns != 0; {
		idx, _ := popLSB(&blackPawns)
		score -= valPawn
		score -= pstBlackPawn[idx]
	}

	for blackBishops := b.black & b.bishops; blackBishops != 0; {
		idx, _ := popLSB(&blackBishops)
		score -= valBishop
		score -= pstBlackBishop[idx]
	}

	for blackKnights := b.black & b.knights; blackKnights != 0; {
		idx, _ := popLSB(&blackKnights)
		score -= valKnight
		score -= pstBlackKnight[idx]
	}

	for blackQueens := b.black & b.queens; blackQueens != 0; {
		idx, _ := popLSB(&blackQueens)
		score -= valQueen
		score -= pstBlackQueen[idx]
	}

	for blackKings := b.black & b.kings; blackKings != 0; {
		idx, _ := popLSB(&blackKings)
		score -= valKing
		score -= pstBlackKing[idx]
	}

	for blackRooks := b.black & b.rooks; blackRooks != 0; {
		idx, _ := popLSB(&blackRooks)
		score -= valRook
		score -= pstBlackRook[idx]
	}

	return score
}
