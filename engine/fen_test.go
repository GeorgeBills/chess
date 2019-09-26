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
		// {"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z 0 1", ""}, // invalid en passant
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
