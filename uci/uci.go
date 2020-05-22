package uci

import (
	"fmt"

	chess "github.com/GeorgeBills/chess/m/v2"
)

// ParseUCIN parses a string in Universal Chess Notation (e.g. "a1h8") as a
// FromToPromote 3-tuple.
func ParseUCIN(ucin string) (chess.Move, error) {
	if len(ucin) != 4 && len(ucin) != 5 {
		return chess.Move{}, fmt.Errorf("invalid length for UCIN: %d", len(ucin))
	}
	from, err := chess.ParseAlgebraicNotationString(ucin[0:2])
	if err != nil {
		return chess.Move{}, err
	}
	to, err := chess.ParseAlgebraicNotationString(ucin[2:4])
	if err != nil {
		return chess.Move{}, err
	}

	promoteTo := chess.PromoteToNone
	if len(ucin) == 5 {
		switch ucin[4] {
		case 'q':
			promoteTo = chess.PromoteToQueen
		case 'n':
			promoteTo = chess.PromoteToKnight
		case 'r':
			promoteTo = chess.PromoteToRook
		case 'b':
			promoteTo = chess.PromoteToBishop
		default:
			return chess.Move{}, fmt.Errorf("invalid promote to char: %q", ucin[4])
		}
	}

	return chess.NewMove(from, to, promoteTo), nil
}

// ToUCIN returns the move in Universal Chess Interface Notation (e.g. "a7a8q").
// UCIN is very similar to, but not exactly the same as, Long Algebraic
// Notation.
func ToUCIN(move chess.FromToPromoter) string {
	from, to := move.From(), move.To()

	ucin := [5]byte{
		'a' + chess.FileIndex(from), // a...h
		'1' + chess.RankIndex(from), // 1...8
		'a' + chess.FileIndex(to),   // a...h
		'1' + chess.RankIndex(to),   // 1...8
	}

	promoteTo := move.PromoteTo()
	if promoteTo != chess.PromoteToNone {
		switch promoteTo {
		case chess.PromoteToQueen:
			ucin[4] = 'q'
		case chess.PromoteToKnight:
			ucin[4] = 'n'
		case chess.PromoteToRook:
			ucin[4] = 'r'
		case chess.PromoteToBishop:
			ucin[4] = 'b'
		default:
			panic(fmt.Errorf("unrecognized promote to: %b", promoteTo))
		}
		return string(ucin[0:5])
	}

	return string(ucin[0:4])
}
