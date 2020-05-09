package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEvaluation(t *testing.T) {
	tests := map[string]string{
		"down a pawn":                  "rnbqkbnr/ppp1pppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"down a knight":                "rnbqkb1r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"down a rook":                  "rnbqkbn1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQq - 0 1",
		"down a queen":                 "rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"down a bishop":                "rn1qkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"queen > bishop + knight":      "rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/R1BQK1NR w KQkq - 0 1",
		"rook > bishop":                "rnbqkbn1/pppppppp/8/8/8/8/PPPPPPPP/RN1QKBNR w KQq - 0 1",
		"pawn about to promote":        "3qkbnr/P2ppppp/8/8/8/8/1PPPP3/RNBQK3 w Qk - 0 1",
		"pawn more central":            "rnbqkbnr/ppppppp1/8/7p/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
		"pawn double pushed vs single": "rnbqkbnr/pppp1ppp/4p3/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
		"knight more central":          "r1bqkbnr/pppppppp/8/n7/3N4/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
		"bishop more central":          "rn1qkbnr/ppp1pppp/8/8/5B2/7b/PPP1PPPP/RN1QKBNR w KQ - 0 1",
	}
	scores := make(map[string]int16, len(tests))
	for name, fen := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
			require.NoError(t, err)
			score := b.Evaluate()

			// white should be winning
			assert.Greater(t, score, int16(0))

			// TODO: use "golden" pattern for scores
			scores[fen] = score

			// TODO: flip the board and check black gets the same score
		})
	}
}

// TODO: benchmark evaluation
