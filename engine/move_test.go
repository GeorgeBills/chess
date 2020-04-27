package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"testing"
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
