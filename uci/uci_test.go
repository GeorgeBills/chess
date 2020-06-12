package uci_test

import (
	"testing"
	"time"

	"github.com/GeorgeBills/chess"
	"github.com/GeorgeBills/chess/uci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	processing = 20 * time.Millisecond
	timeout    = 1 * time.Second
)

func mustParseMove(ucin string) *uci.Move {
	m, err := uci.ParseUCIN(ucin)
	if err != nil {
		panic(err)
	}
	return m
}

func TestUCIN(t *testing.T) {
	tests := []struct {
		ucin      string
		from, to  uint8
		promoteTo chess.PromoteTo
	}{
		{"0000", 0, 0, chess.PromoteToNone},
		{"a1a2", 0, 8, chess.PromoteToNone},
		{"b3b4", 17, 25, chess.PromoteToNone},
		{"c5c6", 34, 42, chess.PromoteToNone},
		{"d7d8", 51, 59, chess.PromoteToNone},
		{"e7e8q", 52, 60, chess.PromoteToQueen},
		{"f7f8r", 53, 61, chess.PromoteToRook},
		{"g7g8n", 54, 62, chess.PromoteToKnight},
		{"h7h8b", 55, 63, chess.PromoteToBishop},
	}
	for _, tt := range tests {
		t.Run(tt.ucin, func(t *testing.T) {
			// parse from ucin
			parsed, err := uci.ParseUCIN(tt.ucin)
			require.NoError(t, err)
			if tt.ucin != "0000" {
				require.NotNil(t, parsed)
				assert.Equal(t, tt.from, parsed.From(), "parsed.From() != %d", tt.from)
				assert.Equal(t, tt.to, parsed.To(), "parsed.To() != %d", tt.to)
				assert.Equal(t, tt.promoteTo, parsed.PromoteTo(), "parsed.PromoteTo() != %s", tt.promoteTo)
			}

			// round trip back to ucin
			ucin := uci.ToUCIN(parsed)
			assert.Equal(t, tt.ucin, ucin)
		})
	}
}
