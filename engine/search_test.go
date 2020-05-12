package engine_test

import (
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBestMoveToDepth(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    uint8
		expected string
	}{
		{
			"depth 1: capture queen",
			"3q3k/8/8/8/8/8/8/3QK3 w - - 0 1",
			1,
			"d1xd8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.fen))
			require.NoError(t, err)
			require.NotNil(t, b)
			g := engine.NewGame(b)
			move, _ := g.BestMoveToDepth(tt.depth, 1)
			assert.Equal(t, tt.expected, move.SAN())
		})
	}
}
