package chess_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess"
	"github.com/stretchr/testify/assert"
)

func TestConversions(t *testing.T) {
	tests := []struct {
		square, rank, file uint8
		an                 string
	}{
		{0, 0, 0, "a1"},
		{33, 4, 1, "B5"},
		{42, 5, 2, "c6"},
		{51, 6, 3, "D7"},
		{12, 1, 4, "e2"},
		{21, 2, 5, "F3"},
		{30, 3, 6, "g4"},
		{63, 7, 7, "H8"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("ParseAlgebraicNotation(%s)", tt.an), func(t *testing.T) {
			rank, file, err := chess.ParseAlgebraicNotation(strings.NewReader(tt.an))
			if assert.NoError(t, err) {
				assert.Equal(t, tt.rank, rank)
				assert.Equal(t, tt.file, file)
			}
		})

		t.Run(fmt.Sprintf("ParseAlgebraicNotationString(%s)", tt.an), func(t *testing.T) {
			rank, file, err := chess.ParseAlgebraicNotationString(tt.an)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.rank, rank)
				assert.Equal(t, tt.file, file)
			}
		})

		t.Run(fmt.Sprintf("SquareIndexToAlgebraicNotation(%d)", tt.square), func(t *testing.T) {
			an := chess.SquareIndexToAlgebraicNotation(tt.square)
			expected := strings.ToLower(tt.an)
			assert.Equal(t, expected, an)
		})

		t.Run(fmt.Sprintf("RankIndex(%d)", tt.square), func(t *testing.T) {
			rank := chess.RankIndex(tt.square)
			assert.Equal(t, tt.rank, rank)
		})

		t.Run(fmt.Sprintf("FileIndex(%d)", tt.square), func(t *testing.T) {
			file := chess.FileIndex(tt.square)
			assert.Equal(t, tt.file, file)
		})

		t.Run(fmt.Sprintf("SquareIndex(%d,%d)", tt.rank, tt.file), func(t *testing.T) {
			square := chess.SquareIndex(tt.rank, tt.file)
			assert.Equal(t, tt.square, square, "rank %d, file %d != sq %d", tt.rank, tt.file, tt.square)
		})
	}
}
