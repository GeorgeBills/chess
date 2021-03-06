package engine_test

import (
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/engine"
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
			"depth 0",
			engine.InitialBoardFEN,
			0,
			"-",
		},
		{
			"depth 1: capture queen",
			"3q3k/8/8/8/8/8/8/3QK3 w - - 0 1",
			1,
			"d1xd8",
		},
		{
			"depth 2: mate in one",
			"5k2/4ppp1/8/8/8/8/8/R2bK3 w Q - 0 1",
			2,
			"a1a8",
		},
		{
			"depth 2: must underpromote to avoid stalemate",
			"4k3/8/8/8/r7/7K/6p1/8 b - - 0 1",
			2,
			"g2g1=R",
		},
		{
			// based on Evans vs Reshevsky "The Mother of All Swindles"
			"depth 3: white to force stalemate",
			"7k/3Q4/8/1p2p2p/1P2Pn1P/5Pq1/8/7K w - - 0 1",
			3,
			"d7h7",
		},
		{
			"depth 3: must move even if checkmate is guaranteed",
			"4k3/8/8/8/3Pn3/8/5K2/3b3q w - - 1 15",
			3,
			"f2e3",
		},
		{
			"depth 4: mate in two",
			"r3k3/r5Q1/8/8/8/8/5PPR/7K b - - 0 1",
			4,
			"a7a1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.fen))
			require.NoError(t, err)
			require.NotNil(t, b)
			g := engine.NewGame(b)
			move, _ := g.BestMoveToDepth(tt.depth, nil, nil)
			assert.Equal(t, tt.expected, move.SAN())
		})
	}
}

func BenchmarkBestMoveToDepth(b *testing.B) {
	const depth = 6

	g := engine.NewGame(engine.NewBoard())

	stopch := make(chan struct{})
	statusch := make(chan engine.SearchStatus)

	var move engine.Move
	for i := 0; i < b.N; i++ {
		move, _ = g.BestMoveToDepth(depth, stopch, statusch)
	}
	b.StopTimer()

	// sanity check our best move
	assert.Contains(b, []string{"g1f3", "e2e4", "d2d4", "c2c4"}, move.SAN())
}
