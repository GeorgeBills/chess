package engine_test

import (
	"fmt"
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
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
