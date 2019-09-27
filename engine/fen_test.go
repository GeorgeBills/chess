package engine_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFEN(t *testing.T) {
	fen := engine.NewBoard().FEN()
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	assert.Equal(t, expected, fen)
}

func TestNewBoardFromFEN(t *testing.T) {
	fen := strings.NewReader("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	expected := engine.NewBoard()
	b, err := engine.NewBoardFromFEN(fen)
	require.NoError(t, err)
	assert.Equal(t, &expected, b)
}

func TestEnPassant(t *testing.T) {
	fen := []struct {
		fen      string
		expected uint8
	}{
		{"rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R b KQkq - 1 1", 0},
		// white pawns
		{"rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq A3 0 1", engine.A3},
		{"rnbqkbnr/pppppppp/8/8/1P6/8/P1PPPPPP/RNBQKBNR b KQkq b3 0 1", engine.B3},
		{"rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR b KQkq C3 0 1", engine.C3},
		{"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1", engine.D3},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq E3 0 1", engine.E3},
		{"rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 1", engine.F3},
		{"rnbqkbnr/pppppppp/8/8/6P1/8/PPPPPP1P/RNBQKBNR b KQkq G3 0 1", engine.G3},
		{"rnbqkbnr/pppppppp/8/8/P7/8/PPPPPPP1/RNBQKBNR b KQkq h3 0 1", engine.H3},
		// black pawns
		{"rnbqkbnr/1ppppppp/8/p7/8/7N/PPPPPPPP/RNBQKB1R w KQkq a6 0 2", engine.A6},
		{"rnbqkbnr/p1pppppp/8/1p6/8/7N/PPPPPPPP/RNBQKB1R w KQkq B6 0 2", engine.B6},
		{"rnbqkbnr/pp1ppppp/8/2p5/8/7N/PPPPPPPP/RNBQKB1R w KQkq c6 0 2", engine.C6},
		{"rnbqkbnr/ppp1pppp/8/3p4/8/7N/PPPPPPPP/RNBQKB1R w KQkq D6 0 2", engine.D6},
		{"rnbqkbnr/pppp1ppp/8/4p3/8/7N/PPPPPPPP/RNBQKB1R w KQkq e6 0 2", engine.E6},
		{"rnbqkbnr/ppppp1pp/8/5p2/8/7N/PPPPPPPP/RNBQKB1R w KQkq F6 0 2", engine.F6},
		{"rnbqkbnr/pppppp1p/8/6p1/8/7N/PPPPPPPP/RNBQKB1R w KQkq g6 0 2", engine.G6},
		{"rnbqkbnr/ppppppp1/8/7p/8/7N/PPPPPPPP/RNBQKB1R w KQkq H6 0 2", engine.H6},
	}
	for i, tt := range fen {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.fen))
			require.NoError(t, err)
			require.NotNil(t, b)
			ep := b.EnPassant()
			assert.Equal(t, tt.expected, ep)
		})
	}
}

func TestNewBoardFromInvalidFEN(t *testing.T) {
	invalid := []struct{ fen, expected string }{
		{"", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 ", "unexpected EOF"},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNRwKQkq-01", "unexpected 'w', expecting ' '"},                    // no whitespace
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1n", "unexpected 'n', expecting [0-9]"},            // invalid full moves char
		{"rnbqkbnr/pppppppp/7/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "unexpected '/', expecting [PNBRQKpnbrqk1-8]"}, // empty n too low
		{"rnbqkbnr/pppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "unexpected '/', expecting [PNBRQKpnbrqk1-8]"},     // rank char too short
		{"pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "unexpected ' ', expecting '/'"},                        // missing rank
		{"rnbqkbnx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "unexpected 'x', expecting [PNBRQKpnbrqk1-8]"}, // invalid piece char
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR ? KQkq - 0 1", "unexpected '?', expecting [wb]"},              // invalid to move char
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w abcd - 0 1", "unexpected 'a', expecting [KQkq]"},            // invalid castling chars
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z1 0 1", "unexpected 'z', expecting [a-hA-H]"},         // invalid en passant
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq a1 0 1", "unexpected '1', expecting [36]"},             // invalid en passant
		{"rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq a6 0 1", "invalid board state: white moved last; en passant on rank 6"},
		{"rnbqkbnr/1ppppppp/8/p7/8/7N/PPPPPPPP/RNBQKB1R w KQkq a3 0 2", "invalid board state: black moved last; en passant on rank 3"},
		// {"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 999999", ""}, // number too large
	}
	for i, tt := range invalid {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.fen))
			assert.EqualError(t, err, tt.expected)
			assert.Nil(t, b)
		})
	}
}
