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
			"k7/8/8/8/3p4/p5P1/P2Pp3/7K w - - 0 123",
			[]string{
				"h1g1", "h1h2", "h1g2", // king
				"d2d3", "g3g4", // pawns
			},
		},
		{
			"pawns can't move if blocked by friendly piece",
			"k7/5pp1/p5n1/1p3p2/8/2pp4/8/7K b - - 0 123",
			[]string{
				"a8a7", "a8b8", "a8b7", // king
				"g6e5", "g6e7", "g6f4", "g6f8", "g6h4", "g6h8", // knight
				"f7f6", "f5f4", "a6a5", "b5b4", "c3c2", "d3d2", // pawns
			},
		},
		{
			"pawn captures",
			"4k3/8/8/2bpp3/3PPPp1/r5PP/1P6/4K3 w - - 0 123",
			[]string{
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"b2xa3", "b2b3", "b2b4", // b2 pawn
				"d4xc5", "d4xe5", // d4 pawn
				"e4xd5",         // e4 pawn
				"f4f5", "f4xe5", // f4 pawn
				"h3xg4", "h3h4", // h3 pawn
			},
		},
		{
			"pawn promotions",
			"rn3rk1/P1PP4/4P3/5P2/8/8/8/4K3 w - - 0 123",
			[]string{
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"e6e7", "f5f6", // a6, f5 pawns can't promote
				// a7 pawn can only capture
				"a7xb8=Q", "a7xb8=N", "a7xb8=R", "a7xb8=B",
				// c7 pawn can either capture or advance
				"c7xb8=Q", "c7xb8=N", "c7xb8=R", "c7xb8=B",
				"c7c8=Q", "c7c8=N", "c7c8=R", "c7c8=B",
				// d7 pawn can only advance
				"d7d8=Q", "d7d8=N", "d7d8=R", "d7d8=B",
			},
		},
		{
			"pawn en passant (black to move)",
			"4k3/8/8/8/3pPp2/3P2P1/8/4K3 b - e3 0 123",
			[]string{
				"e8d7", "e8d8", "e8e7", "e8f7", "e8f8", // king
				"d4xe3e.p.",                  // d4 pawn
				"f4xg3", "f4f3", "f4xe3e.p.", // f4 pawn
			},
		},
		{
			"pawn en passant (white to move)",
			"4k3/8/8/2PpP3/4P3/8/8/K7 w - d6 0 123",
			[]string{
				"a1b1", "a1a2", "a1b2", // king
				"e5xd6e.p.", "e5e6", // e5 pawn
				"c5xd6e.p.", "c5c6", // c5 pawn
				"e4xd5", // e4 pawn
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
			"k7/8/8/8/8/8/8/1R5K w - - 1 123",
			[]string{
				"h1h2", "h1g1", "h1g2", // king
				"b1b2", "b1b3", "b1b4", "b1b5", "b1b6", "b1b7", "b1b8", // rook vertical
				"b1a1", "b1c1", "b1d1", "b1e1", "b1f1", "b1g1", // rook horizontal
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
			"4k3/3b4/8/8/8/8/8/3K4 b - - 1 123",
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
			"queen moves",
			"4k3/2p3P1/B2q3B/8/8/P7/7P/3QK3 b - - 1 123",
			[]string{
				"c7c5", "c7c6", // pawn
				"e8d7", "e8d8", "e8e7", "e8f7", "e8f8", // king
				"d6d7", "d6d8", // queen north
				"d6f8", "d6e7", // queen north east
				"d6f6", "d6e6", "d6g6", "d6xh6", // queen east
				"d6f4", "d6e5", "d6g3", "d6xh2", // queen south east
				"d6d2", "d6d3", "d6d4", "d6d5", "d6xd1", // queen south
				"d6b4", "d6xa3", "d6c5", // queen south west
				"d6c6", "d6b6", "d6xa6", // queen west
			},
		},
		{
			"king must not move into check (bishop)",
			"4k3/8/8/3p4/4Kb2/8/8/8 w - - 1 123",
			[]string{
				"e4xd5", "e4f5",
				"e4d4", "e4xf4",
				"e4d3", "e4f3",
			},
		},
		{
			"king must not move into check (rook)",
			"4k2r/8/8/8/8/8/8/6K1 w - - 1 123",
			[]string{"g1g2", "g1f1", "g1f2"},
		},
		{
			"king must not move into check (knights)",
			"4k3/8/6N1/3N4/8/8/8/4K3 b - - 1 123",
			[]string{"e8d7", "e8d8", "e8f7"},
		},
		{
			"king must not move into check (opposing king)",
			"k7/2K5/7R/8/8/8/8/8 b - - 1 123",
			[]string{"a8a7"},
		},
		// TODO: test for queen threat
		{
			"stalemate (no moves possible)",
			"4k1r1/8/8/8/8/8/r7/7K w - - 1 123",
			nil,
		},
		{
			"king free to move: own piece blocks check",
			"3qk3/8/8/8/8/8/3P4/3K4 w - - 1 123",
			[]string{
				"d2d3", "d2d4", // pawn
				"d1c1", "d1e1", "d1c2", "d1e2", // king
			},
		},
		{
			"king free to move: opposing piece blocks check",
			"3qk3/3b4/8/8/8/8/8/3K4 w - - 1 123",
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
