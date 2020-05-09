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

			// parse the move
			move, err := b.ParseNewMoveFromUCIN(strings.NewReader(tt.Move))
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
