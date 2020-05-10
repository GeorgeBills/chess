package engine_test

import (
	"encoding/json"
	"flag"
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestEvaluate(t *testing.T) {
	tests := map[string]string{
		"material: down a bishop":                        "rn1qkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"material: down a knight":                        "rnbqkb1r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"material: down a pawn":                          "rnbqkbnr/ppp1pppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"material: down a queen":                         "rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"material: down a rook":                          "rnbqkbn1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQq - 0 1",
		"material inequalities: bishop > 3 pawns":        "4k3/3ppp2/8/8/8/8/8/2B1K3 w - - 0 1",
		"material inequalities: knight > 3 pawns":        "4k3/3ppp2/8/8/8/8/8/1N2K3 w - - 0 1",
		"material inequalities: queen > 8 pawns":         "4k3/pppppppp/8/8/8/8/8/3QK3 w - - 0 1",
		"material inequalities: queen > bishop + knight": "1nb1k3/8/8/8/8/8/8/3QK3 w - - 0 1",
		"material inequalities: queen > rook + bishop":   "r1b1k3/8/8/8/8/8/8/3QK3 w q - 0 1",
		"material inequalities: queen > rook + knight":   "4k1nr/8/8/8/8/8/8/3QK3 w k - 0 1",
		"material inequalities: rook > 4 pawns":          "4k3/2pppp2/8/8/8/8/8/4K2R w K - 0 1",
		"material inequalities: rook > bishop":           "2b1k3/8/8/8/8/8/8/4K2R w K - 0 1",
		"material inequalities: rook > knight":           "4k1n1/8/8/8/8/8/8/R3K3 w Q - 0 1",
		"positioning: bishop more central":               "rn1qkbnr/ppp1pppp/8/8/5B2/7b/PPP1PPPP/RN1QKBNR w KQ - 0 1",
		"positioning: king behind pawn shelter":          "5rk1/2ppp3/8/8/8/8/5PPP/5RK1 w - - 0 1",
		"positioning: king back (opening-middle game)":   "rn1q4/3ppp2/8/8/4k3/8/3PPP2/3QK1NR w - - 0 1",
		"positioning: king forward (end game)":           "3qk1n1/3ppp2/8/4K3/8/8/3PPP2/3Q2N1 w - - 0 1",
		"positioning: knight more central":               "r1bqkbnr/pppppppp/8/n7/3N4/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1",
		"positioning: pawn about to promote":             "3qkbnr/P2ppppp/8/8/8/8/1PPPP3/RNBQK3 w Qk - 0 1",
		"positioning: pawn double pushed vs single":      "rnbqkbnr/pppp1ppp/4p3/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
		"positioning: pawn more central":                 "rnbqkbnr/ppppppp1/8/7p/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1",
	}

	scores := make(map[string]int16, len(tests))
	var golden map[string]int16
	f, err := os.Open("testdata/evaluate.golden.json")
	require.NoError(t, err)
	err = json.NewDecoder(f).Decode(&golden)
	require.NoError(t, err)

	for name, fen := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
			require.NoError(t, err)
			score := b.Evaluate()

			// white should be winning
			assert.Greater(t, score, int16(0))

			// swap the board and check black gets the same score, but negative
			swapped := b.Swapped()
			assert.Equal(t, score, -1*swapped.Evaluate())

			scores[fen] = score
		})
	}

	if *update {
		t.Log("updating golden file")
		fw, err := os.Create("testdata/evaluate.golden.json")
		require.NoError(t, err)
		enc := json.NewEncoder(fw)
		enc.SetIndent("", "    ")
		err = enc.Encode(scores)
		require.NoError(t, err)
	} else {
		assert.Equal(t, golden, scores, "evaluation scores didn't match golden file")
	}
}

func BenchmarkEvaluate(b *testing.B) {
	board := engine.NewBoard()
	for i := 0; i < b.N; i++ {
		board.Evaluate()
	}
}
