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
