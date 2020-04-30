package engine_test

import (
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

func TestMakeUnmakeMove(t *testing.T) {
	tests := []struct {
		name          string
		before, after string
		move          engine.Move
	}{
		{
			"pawn single push (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			NewMove(D2, D3),
		},
		{
			"pawn single push (black)", // bumps total moves
			"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			"rnbqkbnr/ppp1pppp/3p4/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			NewMove(D7, D6),
		},
		{
			"pawn double push (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			NewPawnDoublePush(E2, E4),
		},
		{
			"pawn capture (white)",
			"rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			"rnbqkbnr/pppp1ppp/8/4P3/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2",
			NewCapture(D4, E5),
		},
		{
			"en passant (white)",
			"r1bqkbnr/pppp1ppp/2n5/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"r1bqkbnr/pppp1ppp/2n1P3/8/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 3",
			NewEnPassant(D5, E6),
		},
		{
			"en passant declined (white)",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 0 3",
			NewMove(G1, F3),
		},
		{
			"castle queenside (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/2KR1BNR b kq - 0 5",
			WhiteQueensideCastle,
		},
		{
			"castle queenside (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"2kr1bnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQ - 0 6",
			BlackQueensideCastle,
		},
		{
			"castle kingside (white)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQ1RK1 b kq - 0 4",
			WhiteKingsideCastle,
		},
		{
			"castle kingside (black)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R b KQkq - 0 4",
			"rnbq1rk1/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R w KQ - 0 5",
			BlackKingsideCastle,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBoardFromFEN(strings.NewReader(tt.before))
			require.NoError(t, err)
			require.NotNil(t, b)

			g := NewGame(b)

			// make the move and check the FEN is correct
			g.MakeMove(tt.move)
			assert.Equal(t, tt.after, b.FEN(), "FEN after MakeMove() should match 'after' FEN'")

			// reverse the move and check we're back to the original FEN
			g.UnmakeMove()
			assert.Equal(t, tt.before, b.FEN(), "FEN after UnmakeMove() should match 'before' FEN'")
		})
	}
}
