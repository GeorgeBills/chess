package engine_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMove_SAN(t *testing.T) {
	tests := []struct {
		name        string
		move        engine.Move
		expected    string
		isEnPassant bool
		isPromotion bool
	}{
		{"quiet move", engine.NewMove(C3, E5), "c3e5", false, false},
		{"capture", engine.NewCapture(A1, A3), "a1xa3", false, false},
		{"en passant", engine.NewEnPassant(D5, C6), "d5xc6e.p.", true, false},
		{"kingside castling", engine.WhiteKingsideCastle, "O-O", false, false},
		{"queenside castling", engine.BlackQueensideCastle, "O-O-O", false, false},
		{"promotion to queen", engine.NewQueenPromotion(A7, A8, false), "a7a8=Q", false, true},
		{"promotion to bishop with capture", engine.NewBishopPromotion(B2, B1, true), "b2xb1=B", false, true},
		{"promotion to rook", engine.NewRookPromotion(C2, C1, false), "c2c1=R", false, true},
		{"promotion to knight with capture", engine.NewKnightPromotion(D7, D8, true), "d7xd8=N", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.move.SAN())
			assert.Equal(t, tt.isEnPassant, tt.move.IsEnPassant())
			assert.Equal(t, tt.isPromotion, tt.move.IsPromotion())
		})
	}
}

// FIXME: better tests for parsing
//        including parsing and attempting to play invalid moves

func TestParseMakeUnmakeMove(t *testing.T) {
	// TODO: test capturing A8, A1, H1 for completeness
	var tests map[string]struct {
		Before string
		Move   string
		After  string
	}

	f, err := os.Open("testdata/moves-before-after.json")
	require.NoError(t, err)
	err = json.NewDecoder(f).Decode(&tests)
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := NewBoardFromFEN(strings.NewReader(tt.Before))
			require.NoError(t, err)
			require.NotNil(t, b)

			g := NewGame(b)

			parsed, err := ParseLongAlgebraicNotationString(tt.Move)
			require.NoError(t, err)

			// parse the move
			move, err := b.HydrateMove(parsed)
			require.NoError(t, err)
			require.NotNil(t, move)

			// make the move and check the FEN is correct
			g.MakeMove(move)
			assert.Equal(t, tt.After, b.FEN(), "FEN after MakeMove() should match 'after' FEN'")

			// reverse the move and check we're back to the original FEN
			g.UnmakeMove()
			assert.Equal(t, tt.Before, b.FEN(), "FEN after UnmakeMove() should match 'before' FEN'")
		})
	}
}

func TestMakeUnmakeMoveHistory(t *testing.T) {
	// this is the Opera Game: https://en.wikipedia.org/wiki/Opera_Game
	opera := []string{
		"e2e4", "e7e5",
		"g1f3", "d7d6",
		"d2d4", "c8g4",
		"d4e5", "g4f3",
		"d1f3", "d6e5",
		"f1c4", "g8f6",
		"f3b3", "d8e7",
		"b1c3", "c7c6",
		"c1g5", "b7b5",
		"c3b5", "c6b5",
		"c4b5", "b8d7",
		"e1c1", "a8d8",
		"d1d7", "d8d7",
		"h1d1", "e7e6",
		"b5d7", "f6d7",
		"b3b8", "d7b8",
		"d1d8",
	}
	b := engine.NewBoard()
	g := engine.NewGame(&b)
	for _, ucin := range opera {
		parsed, err := ParseLongAlgebraicNotationString(ucin)
		require.NoError(t, err)
		move, err := b.HydrateMove(parsed)
		require.NoError(t, err)
		g.MakeMove(move)
	}
	for i := 0; i < len(opera); i++ {
		g.UnmakeMove()
	}
	assert.Equal(t, engine.InitialBoardFEN, b.FEN())
}
