package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMove_SAN(t *testing.T) {
	tests := []struct {
		name     string
		move     engine.Move
		expected string
	}{
		{"quiet move", engine.NewMove(C3, E5), "c3e5"},
		{"capture", engine.NewCapture(A1, A3), "a1xa3"},
		{"en passant", engine.NewEnPassant(D5, C6), "d5xc6e.p."},
		{"kingside castling", engine.NewKingsideCastle(E1, C1), "O-O"},
		{"queenside castling", engine.NewQueensideCastle(E8, G8), "O-O-O"},
		{"promotion to queen", engine.NewQueenPromotion(A7, A8, false), "a7a8=Q"},
		{"promotion to bishop with capture", engine.NewBishopPromotion(B2, B1, true), "b2xb1=B"},
		{"promotion to rook", engine.NewRookPromotion(C2, C1, false), "c2c1=R"},
		{"promotion to knight with capture", engine.NewKnightPromotion(D7, D8, true), "d7xd8=N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.move.SAN())
		})
	}
}
