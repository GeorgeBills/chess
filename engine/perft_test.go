package engine_test

import (
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// https://www.chessprogramming.org/Perft

func perft(g *engine.Game, depth uint8) uint64 {
	moves := make([]engine.Move, 0, 96)
	moves, _ = g.GenerateLegalMoves(moves)
	if depth == 1 {
		return uint64(len(moves))
	}
	var n uint64
	for _, move := range moves {
		g.MakeMove(move)
		n += perft(g, depth-1)
		g.UnmakeMove()
	}
	return n
}

func TestPerft(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestPerft() due to -short flag")
	}

	tests := []struct {
		name     string
		fen      string
		depth    uint8
		expected uint64
	}{
		{
			"initial",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			7,
			3_195_901_860,
		},
		{
			"kiwipete", // https://www.chessprogramming.org/Perft_Results
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 123",
			6,
			8_031_647_685,
		},
		{
			"position 3", // https://www.chessprogramming.org/Perft_Results
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 123",
			8,
			3_009_794_393,
		},
		{
			"position 4", // https://www.chessprogramming.org/Perft_Results
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 123",
			6,
			706_045_033,
		},
		{
			"position 5", // https://www.chessprogramming.org/Perft_Results
			"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			5,
			89_941_194,
		},
		{
			"position 6", // https://www.chessprogramming.org/Perft_Results
			"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			6,
			6_923_051_137,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fen, depth, expected := tt.fen, tt.depth, tt.expected

			t.Parallel()

			b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
			require.NoError(t, err)
			require.NotNil(t, b)
			g := engine.NewGame(b)
			n := perft(g, depth)

			assert.Equal(t, expected, n)
		})
	}
}

func BenchmarkPerft(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping BenchmarkPerft() due to -short flag")
	}

	// https://www.chessprogramming.org/Perft_Results#Position_4
	// good mix of material, good mix of move types (including castling, en
	// passant, promotions, checks and mates), plausible looking board state,
	// not too many moves so the perft will be reasonably quick.
	const fen = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	const depth = 6
	const expected = 706_045_033
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	g := engine.NewGame(board)

	var n uint64
	for i := 0; i < b.N; i++ {
		n = perft(g, depth)
	}
	b.StopTimer()

	// sanity check our expected perft
	assert.EqualValues(b, expected, n)
}
