package chess_test

import (
	"testing"

	chess "github.com/GeorgeBills/chess/m/v2"
	"github.com/stretchr/testify/assert"
)

func TestToAlgebraicNotation(t *testing.T) {
	tests := []struct {
		i        uint8
		expected string
	}{
		{0, "a1"},
		{1, "b1"},
		{8, "a2"},
		{63, "h8"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			an := chess.SquareIndexToAlgebraicNotation(tt.i)
			assert.Equal(t, tt.expected, an)
		})
	}
}
