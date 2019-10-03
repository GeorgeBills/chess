package engine_test

import (
	"fmt"
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

func TestWhitePawnMoves(t *testing.T) {
	indexes := []struct {
		i        uint8
		expected uint64
	}{
		{A2, 0b00000000_00000000_00000000_00000000_00000001_00000001_00000000_00000000}, // a3, a4
		{B3, 0b00000000_00000000_00000000_00000000_00000010_00000000_00000000_00000000}, // b4
		{H7, 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000}, // h8
	}
	for _, tt := range indexes {
		t.Run(fmt.Sprintf("%d", tt.i), func(t *testing.T) {
			moves := engine.WhitePawnMoves(tt.i)
			assert.Equal(t, tt.expected, moves)
		})
	}
}

func TestBlackPawnMoves(t *testing.T) {
	indexes := []struct {
		i        uint8
		expected uint64
	}{
		{D7, 0b00000000_00000000_00001000_00001000_00000000_00000000_00000000_00000000}, // d6, d5
		{E6, 0b00000000_00000000_00000000_00010000_00000000_00000000_00000000_00000000}, // e5
		{F5, 0b00000000_00000000_00000000_00000000_00100000_00000000_00000000_00000000}, // e4
	}
	for _, tt := range indexes {
		t.Run(fmt.Sprintf("%d", tt.i), func(t *testing.T) {
			moves := engine.BlackPawnMoves(tt.i)
			assert.Equal(t, tt.expected, moves)
		})
	}
}

func BenchmarkPawnMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idx := uint8(i % 64)
		engine.WhitePawnMoves(idx)
		engine.BlackPawnMoves(idx)
	}
}

func TestKingMoves(t *testing.T) {
	indexes := []struct {
		i        uint8
		expected uint64
	}{
		{A1, 0b00000000_00000000_00000000_00000000_00000000_00000000_00000001_00000010}, // a2, b1
		{B2, 0b00000000_00000000_00000000_00000000_00000000_00000010_00000101_00000010}, // b1, b3, a2, c2
		{H8, 0b01000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000}, // g8, h7
	}
	for _, tt := range indexes {
		t.Run(fmt.Sprintf("%d", tt.i), func(t *testing.T) {
			moves := engine.KingMoves(tt.i)
			assert.Equal(t, tt.expected, moves)
		})
	}
}

func BenchmarkKingMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		engine.KingMoves(uint8(i % 64))
	}
}

func TestKnightMoves(t *testing.T) {
	indexes := []struct {
		i        uint8
		expected uint64
	}{
		{D5, 0b00000000_00010100_00100010_00000000_00100010_00010100_00000000_00000000}, // e7, f6, f4, e3, c3, b4, b6, c7
		{A6, 0b00000010_00000100_00000000_00000100_00000010_00000000_00000000_00000000}, // b8, c7, c5, b4
		{H1, 0b00000000_00000000_00000000_00000000_00000000_01000000_00100000_00000000}, // f2, g3
		{G7, 0b00010000_00000000_00010000_10100000_00000000_00000000_00000000_00000000}, // h5, f5, e6, e8
		{C3, 0b00000000_00000000_00000000_00001010_00010001_00000000_00010001_00001010}, // d5, e4, e2, d1, b1, a2, a4, b5
	}
	for _, tt := range indexes {
		t.Run(fmt.Sprintf("%d", tt.i), func(t *testing.T) {
			moves := engine.KnightMoves(tt.i)
			assert.Equal(t, tt.expected, moves)
		})
	}
}

func BenchmarkKnightMoves(b *testing.B) {
	for i := 0; i < b.N; i++ {
		engine.KnightMoves(uint8(i % 64))
	}
}

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
				"rnbqkbnr/pppppppp/8/8/8/N7/PPPPPPPP/R1BQKBNR b KQkq - 1 1",  // 1. Na3
				"rnbqkbnr/pppppppp/8/8/8/2N5/PPPPPPPP/R1BQKBNR b KQkq - 1 1", // 1. Nc3
				"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1", // 1. Nf3
				"rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R b KQkq - 1 1",  // 1. Nh3
				"rnbqkbnr/pppppppp/8/8/8/P7/1PPPPPPP/RNBQKBNR b KQkq - 0 1",  // 1. a3
				"rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq - 0 1",  // 1. a4
				"rnbqkbnr/pppppppp/8/8/8/1P6/P1PPPPPP/RNBQKBNR b KQkq - 0 1", // 1. b3
				"rnbqkbnr/pppppppp/8/8/1P6/8/P1PPPPPP/RNBQKBNR b KQkq - 0 1", // 1. b4
				"rnbqkbnr/pppppppp/8/8/8/2P5/PP1PPPPP/RNBQKBNR b KQkq - 0 1", // 1. c3
				"rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1", // 1. c4
				"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1", // 1. d3
				"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1", // 1. d4
				"rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1", // 1. e3
				"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1", // 1. e4
				"rnbqkbnr/pppppppp/8/8/8/5P2/PPPPP1PP/RNBQKBNR b KQkq - 0 1", // 1. f3
				"rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq - 0 1", // 1. f4
				"rnbqkbnr/pppppppp/8/8/8/6P1/PPPPPP1P/RNBQKBNR b KQkq - 0 1", // 1. g3
				"rnbqkbnr/pppppppp/8/8/6P1/8/PPPPPP1P/RNBQKBNR b KQkq - 0 1", // 1. g4
				"rnbqkbnr/pppppppp/8/8/8/7P/PPPPPPP1/RNBQKBNR b KQkq - 0 1",  // 1 .h3
				"rnbqkbnr/pppppppp/8/8/7P/8/PPPPPPP1/RNBQKBNR b KQkq - 0 1",  // 1. h4
			},
		},
		{
			"1. Nf3 (Réti opening)",
			"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
			[]string{
				"r1bqkbnr/pppppppp/n7/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 2 2",
				"r1bqkbnr/pppppppp/2n5/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 2 2",
				"rnbqkb1r/pppppppp/5n2/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 2 2",
				"rnbqkb1r/pppppppp/7n/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 2 2",
				"rnbqkbnr/1ppppppp/8/p7/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/1ppppppp/p7/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/p1pppppp/8/1p6/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/p1pppppp/1p6/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pp1ppppp/8/2p5/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pp1ppppp/2p5/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppp1pppp/8/3p4/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppp1pppp/3p4/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pppp1ppp/4p3/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppppp1pp/8/5p2/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppppp1pp/5p2/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pppppp1p/8/6p1/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/pppppp1p/6p1/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppppppp1/8/7p/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
				"rnbqkbnr/ppppppp1/7p/8/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 1 2",
			},
		},
		{
			"pawns can't double move if blocked by opposing piece",
			"k7/8/8/8/8/p7/P7/7K w - - 0 123",
			[]string{
				"k7/8/8/8/8/p7/P7/6K1 b - - 1 123",
				"k7/8/8/8/8/p7/P6K/8 b - - 1 123",
			},
		},
		{
			"pawns can't double move if blocked by friendly piece",
			"k7/6p1/6n1/8/8/8/8/7K b - - 0 123",
			[]string{
				"8/k5p1/6n1/8/8/8/8/7K w - - 1 124",
				"1k6/6p1/6n1/8/8/8/8/7K w - - 1 124",
				"k7/6p1/8/8/5n2/8/8/7K w - - 1 124",
				"k7/6p1/8/8/7n/8/8/7K w - - 1 124",
				"k7/6p1/8/4n3/8/8/8/7K w - - 1 124",
				"k7/4n1p1/8/8/8/8/8/7K w - - 1 124",
				"k4n2/6p1/8/8/8/8/8/7K w - - 1 124",
				"k6n/6p1/8/8/8/8/8/7K w - - 1 124",
			},
		},
		{
			"simple rook moves",
			"k7/8/8/8/8/8/1R6/7K w - - 1 123",
			[]string{
				// king
				"k7/8/8/8/8/8/1R6/6K1 b - - 2 123",
				"k7/8/8/8/8/8/1R5K/8 b - - 2 123",
				// rook vertical (along file)
				"k7/8/8/8/8/1R6/8/7K b - - 2 123",
				"k7/8/8/8/1R6/8/8/7K b - - 2 123",
				"k7/8/8/1R6/8/8/8/7K b - - 2 123",
				"k7/8/1R6/8/8/8/8/7K b - - 2 123",
				"k7/1R6/8/8/8/8/8/7K b - - 2 123",
				"kR6/8/8/8/8/8/8/7K b - - 2 123",
				// rook horizontal (along rank)
				"k7/8/8/8/8/8/2R5/7K b - - 2 123",
				"k7/8/8/8/8/8/3R4/7K b - - 2 123",
				"k7/8/8/8/8/8/4R3/7K b - - 2 123",
				"k7/8/8/8/8/8/5R2/7K b - - 2 123",
				"k7/8/8/8/8/8/6R1/7K b - - 2 123",
				"k7/8/8/8/8/8/7R/7K b - - 2 123",
				"k7/8/8/8/8/8/8/1R5K b - - 2 123",
				"k7/8/8/8/8/8/R7/7K b - - 2 123",
			},
		},
		{
			"rook captures",
			"k7/8/8/8/1q6/8/rR2P3/1N5K w - - 1 123",
			[]string{
				// knight
				"k7/8/8/8/1q6/N7/rR2P3/7K b - - 2 123",
				"k7/8/8/8/1q6/2N5/rR2P3/7K b - - 2 123",
				"k7/8/8/8/1q6/8/rR1NP3/7K b - - 2 123",
				// pawn
				"k7/8/8/8/1q6/4P3/rR6/1N5K b - - 1 123",
				"k7/8/8/8/1q2P3/8/rR6/1N5K b - - 1 123",
				// king
				"k7/8/8/8/1q6/8/rR2P2K/1N6 b - - 2 123",
				"k7/8/8/8/1q6/8/rR2P3/1N4K1 b - - 2 123",
				// rook
				"k7/8/8/8/1q6/1R6/r3P3/1N5K b - - 2 123",
				"k7/8/8/8/1R6/8/r3P3/1N5K b - - 1 123",
				"k7/8/8/8/1q6/8/R3P3/1N5K b - - 1 123",
				"k7/8/8/8/1q6/8/r1R1P3/1N5K b - - 2 123",
				"k7/8/8/8/1q6/8/r2RP3/1N5K b - - 2 123",
			},
		},
		{
			"king must not move into check",
			"4k2r/8/8/8/8/8/8/6K1 w - - 1 123",
			[]string{
				"4k2r/8/8/8/8/8/6K1/8 b - - 2 123",
				"4k2r/8/8/8/8/8/8/5K2 b - - 2 123",
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
				"3qk3/8/8/8/3P4/8/8/3K4 b KQkq - 1 123",
				"3qk3/8/8/8/8/3P4/8/3K4 b KQkq - 1 123",
				"3qk3/8/8/8/8/8/3P4/2K5 b KQkq - 2 123",
				"3qk3/8/8/8/8/8/3P4/4K3 b KQkq - 2 123",
			},
		},
		{
			"king free to move: opposing piece blocks check",
			"3qk3/3b4/8/8/8/8/8/3K4 w KQkq - 1 123",
			[]string{
				"3qk3/3b4/8/8/8/8/3K4/8 b KQkq - 2 123",
				"3qk3/3b4/8/8/8/8/8/2K5 b KQkq - 2 123",
				"3qk3/3b4/8/8/8/8/8/4K3 b KQkq - 2 123",
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
				moves = append(moves, move.FEN())
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
