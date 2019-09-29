package engine_test

import (
	"fmt"
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
