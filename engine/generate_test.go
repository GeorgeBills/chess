package engine_test

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateLegalMoves(t *testing.T) {
	var tests map[string]struct {
		FEN   string
		Moves []string
	}

	f, err := os.Open("testdata/legal-moves.json")
	require.NoError(t, err)
	err = json.NewDecoder(f).Decode(&tests)
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.FEN))
			require.NoError(t, err)
			require.NotNil(t, b)

			var san []string
			moves, _ := b.GenerateLegalMoves(nil)
			for _, move := range moves {
				san = append(san, move.SAN())
			}

			if !(len(tt.Moves) == 0 && len(san) == 0) {
				// sort so we don't need to fiddle with ordering
				sort.Strings(tt.Moves)
				sort.Strings(san)
				assert.Equal(t, tt.Moves, san)
			}
		})
	}
}

func TestTooManyCheckersPanics(t *testing.T) {
	fen := "4k3/4r3/8/q7/7b/8/8/4K3 w - - 0 123" // 3 checkers
	b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
	require.NoError(t, err)
	require.NotNil(t, b)
	assert.Panics(t, func() { b.GenerateLegalMoves(nil) })
}

func BenchmarkGenerateLegalMoves10(b *testing.B) {
	const fen = "r3k2r/pbqnbppp/1p2pn2/2p1N3/Q1P5/4P3/PB1PBPPP/RN3RK1 w kq - 8 11" // 10 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateLegalMoves(moves)
	}
}

func BenchmarkGenerateLegalMoves20(b *testing.B) {
	const fen = "4rrk1/2qn2pp/pp2pb2/2p2p2/P1P2P2/2NPPR2/1BQ3PP/1R4K1 b - - 0 20" // 20 ply in, black to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateLegalMoves(moves)
	}
}

func BenchmarkGenerateLegalMoves30(b *testing.B) {
	const fen = "3rr1k1/1nq4p/pp4p1/2pP1p2/P4P2/2Q1P3/1R2N1PP/3R2K1 w - - 0 31" // 30 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateLegalMoves(moves)
	}
}

func BenchmarkGenerateLegalMoves40(b *testing.B) {
	const fen = "3r2k1/1n5p/8/pPpq1p1p/5P2/4P3/6PK/1R2QN2 b - - 3 40" // 40 ply in, black to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateLegalMoves(moves)
	}
}

func BenchmarkGenerateLegalMoves50(b *testing.B) {
	const fen = "6k1/1n5p/8/p7/2p2PP1/1r2P1N1/8/R5K1 w - - 3 51" // 50 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateLegalMoves(moves)
	}
}
