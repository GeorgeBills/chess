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
