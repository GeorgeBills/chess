package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

func TestMoves(t *testing.T) {
	moves := []struct {
		name     string
		board    string
		expected []string
	}{
		{
			"initial board",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			[]string{
				// knights
				"b1a3", "b1c3",
				"g1f3", "g1h3",
				// pawns
				"a2a3", "a2a4",
				"b2b3", "b2b4",
				"c2c3", "c2c4",
				"d2d3", "d2d4",
				"e2e3", "e2e4",
				"f2f3", "f2f4",
				"g2g3", "g2g4",
				"h2h3", "h2h4",
			},
		},
		{
			"1. Nf3 (Réti opening)",
			"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
			[]string{
				// knights
				"b8a6", "b8c6",
				"g8f6", "g8h6",
				// pawns
				"a7a6", "a7a5",
				"b7b6", "b7b5",
				"c7c6", "c7c5",
				"d7d6", "d7d5",
				"e7e6", "e7e5",
				"f7f6", "f7f5",
				"g7g6", "g7g5",
				"h7h6", "h7h5",
			},
		},
		{
			"pawns can't move if blocked by opposing piece",
			"k7/8/8/8/8/p7/P7/7K w - - 0 123",
			[]string{
				"h1g1", "h1h2", "h1g2", // king
			},
		},
		{
			"pawns can't move if blocked by friendly piece",
			"k7/6p1/6n1/8/8/8/8/7K b - - 0 123",
			[]string{
				"a8a7", "a8b8", "a8b7", // king
				"g6e5", "g6e7", "g6f4", "g6f8", "g6h4", "g6h8", // knight
			},
		},
		{
			"knight moves",
			"4k3/8/8/3p4/p7/2N5/P3P3/4K3 w - - 1 123",
			[]string{
				"c3b1", "c3b5", "c3d1", "c3e4", "c3xa4", "c3xd5", // knight
				"e1d1", "e1d2", "e1f1", "e1f2", // king
				"e2e3", "e2e4", "a2a3", // pawns
			},
		},
		{
			"rook moves",
			"k7/8/8/8/8/8/1R6/7K w - - 1 123",
			[]string{
				"h1h2", "h1g1", "h1g2", // king
				// rook vertical (along file)
				"b2b1", "b2b3", "b2b4", "b2b5", "b2b6", "b2b7", "b2b8",
				// rook horizontal (along rank)
				"b2a2", "b2c2", "b2d2", "b2e2", "b2f2", "b2g2", "b2h2",
			},
		},
		{
			"rook captures",
			"k7/8/8/8/1q6/8/rR2P3/1N5K w - - 1 123",
			[]string{
				"b1d2", "b1c3", "b1a3", // knight
				"e2e3", "e2e4", // pawn
				"h1g1", "h1h2", "h1g2", // king
				"b2xa2", "b2xb4", "b2b3", "b2c2", "b2d2", // rook
			},
		},
		{
			"bishop moves",
			"4k3/3b4/8/8/8/8/8/3K4 b KQkq - 1 123",
			[]string{
				"e8e7", "e8f7", "e8f8", "e8d8", // king
				"d7c8", "d7c6", "d7b5", "d7a4", "d7e6", "d7f5", "d7g4", "d7h3", // bishop
			},
		},
		{
			"bishop captures",
			"4k3/3p4/p7/1B6/8/3K4/8/8 w - - 1 123",
			[]string{
				"d3c2", "d3c3", "d3c4", "d3d2", "d3d4", "d3e2", "d3e3", "d3e4", // king
				"b5c6", "b5c4", "b5a4", "b5xa6", "b5xd7", // bishop
			},
		},
		{
			"king must not move into check",
			"4k2r/8/8/8/8/8/8/6K1 w - - 1 123",
			[]string{
				"g1g2", "g1f1", "g1f2", // king
			},
		},
		{
			"stalemate (no moves possible)",
			"4k1r1/8/8/8/8/8/r7/7K w KQkq - 1 123",
			nil,
		},
		{
			"king free to move: own piece blocks check",
			"3qk3/8/8/8/8/8/3P4/3K4 w KQkq - 1 123",
			[]string{
				"d2d3", "d2d4", // pawn
				"d1c1", "d1e1", "d1c2", "d1e2", // king
			},
		},
		{
			"king free to move: opposing piece blocks check",
			"3qk3/3b4/8/8/8/8/8/3K4 w KQkq - 1 123",
			[]string{
				"d1c1", "d1e1", "d1d2", "d1c2", "d1e2", // king
			},
		},
	}

	for _, tt := range moves {
		t.Run(tt.name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.board))
			require.NoError(t, err)
			require.NotNil(t, b)
			var moves []string
			for _, move := range b.Moves() {
				moves = append(moves, move.SAN())
			}
			// sort so we don't need to fiddle with ordering in the test case
			sort.Strings(tt.expected)
			sort.Strings(moves)
			assert.Equal(t, tt.expected, moves)
		})
	}

}

func BenchmarkInitialMoves(b *testing.B) {
	board := engine.NewBoard()
	for i := 0; i < b.N; i++ {
		board.Moves()
	}
}
