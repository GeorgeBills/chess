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

func TestParseMakeUnmakeMove(t *testing.T) {
	tests := []struct {
		name          string
		before, after string
		move          string
	}{
		{
			"pawn single push (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			"d2d3",
		},
		{
			"pawn single push (black)",
			"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			"rnbqkbnr/ppp1pppp/3p4/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			"d7d6",
		},
		{
			"pawn double push (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			"e2e4",
		},
		{
			"pawn double push (black)",
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
			"e7e5",
		},
		{
			"en passant (white)",
			"r1bqkbnr/pppp1ppp/2n5/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"r1bqkbnr/pppp1ppp/2n1P3/8/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 3",
			"d5e6",
		},
		{
			"en passant (black)",
			"rnbqkbnr/pppp1ppp/8/8/3PpP2/2N5/PPP1P1PP/R1BQKBNR b KQkq f3 0 3",
			"rnbqkbnr/pppp1ppp/8/8/3P4/2N2p2/PPP1P1PP/R1BQKBNR w KQkq - 0 4",
			"e4f3",
		},
		{
			"en passant declined (white)",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 0 3",
			"g1f3",
		},
		{
			"white pawn capture black pawn",
			"rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			"rnbqkbnr/pppp1ppp/8/4P3/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2",
			"d4e5",
		},
		{
			"white bishop capture black knight",
			"r1bqkb1r/pppppppp/2n2n2/1B6/8/4P3/PPPP1PPP/RNBQK1NR w KQkq - 0 3",
			"r1bqkb1r/pppppppp/2B2n2/8/8/4P3/PPPP1PPP/RNBQK1NR b KQkq - 0 3",
			"b5c6",
		},
		{
			"black queen capture white rook",
			"4k3/8/5q2/8/8/8/8/R3K3 b - - 0 123",
			"4k3/8/8/8/8/8/8/q3K3 w - - 0 124",
			"f6a1",
		},
		{
			"castle queenside (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/2KR1BNR b kq - 0 5",
			"e1c1",
		},
		{
			"castle queenside (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"2kr1bnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQ - 0 6",
			"e8c8",
		},
		{
			"castle kingside (white)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQ1RK1 b kq - 0 4",
			"e1g1",
		},
		{
			"castle kingside (black)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R b KQkq - 0 4",
			"rnbq1rk1/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R w KQ - 0 5",
			"e8g8",
		},
		{
			"moving king - can no longer castle (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R2K1BNR b kq - 0 5",
			"e1d1",
		},
		{
			"moving queenside rook - can no longer castle queenside (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/1R2KBNR b Kkq - 0 5",
			"a1b1",
		},
		{
			"moving kingside rook - can no longer castle kingside (white)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK1R1 b Qkq - 0 4",
			"h1g1",
		},
		{
			"moving king - can no longer castle (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"r2k1bnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQ - 0 6",
			"e8d8",
		},
		{
			"moving queenside rook - can no longer castle queenside (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"1r2kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQk - 0 6",
			"a8b8",
		},
		{
			"moving kingside rook - can no longer castle kingside (black)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R b KQkq - 0 4",
			"rnbqk1r1/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R w KQq - 0 5",
			"h8g8",
		},
		{
			"capturing kingside rook - can no longer castle kingside (white to move)",
			"4k2r/8/6N1/8/8/8/8/2KR4 w k - 1 125",
			"4k2N/8/8/8/8/8/8/2KR4 b - - 1 125",
			"g6h8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBoardFromFEN(strings.NewReader(tt.before))
			require.NoError(t, err)
			require.NotNil(t, b)

			g := NewGame(b)

			move, err := b.ParseNewMoveFromUCIN(strings.NewReader(tt.move))
			require.NoError(t, err)
			require.NotNil(t, move)

			// make the move and check the FEN is correct
			g.MakeMove(move)
			assert.Equal(t, tt.after, b.FEN(), "FEN after MakeMove() should match 'after' FEN'")

			// reverse the move and check we're back to the original FEN
			g.UnmakeMove()
			assert.Equal(t, tt.before, b.FEN(), "FEN after UnmakeMove() should match 'before' FEN'")
		})
	}
}

func TestMakeUnmakeMove(t *testing.T) {
	tests := []struct {
		name          string
		before, after string
		move          engine.Move
	}{
		// FIXME: capturing rook starting position should unset castling flag
		{
			"pawn single push (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			"rnbqkbnr/pppppppp/8/8/8/3P4/PPP1PPPP/RNBQKBNR b KQkq - 0 1",
			NewMove(D2, D3),
		},
		{
			"pawn single push (black)",
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
			"en passant (white)",
			"r1bqkbnr/pppp1ppp/2n5/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"r1bqkbnr/pppp1ppp/2n1P3/8/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 3",
			NewEnPassant(D5, E6),
		},
		{
			"en passant (black)",
			"rnbqkbnr/pppp1ppp/8/8/3PpP2/2N5/PPP1P1PP/R1BQKBNR b KQkq f3 0 3",
			"rnbqkbnr/pppp1ppp/8/8/3P4/2N2p2/PPP1P1PP/R1BQKBNR w KQkq - 0 4",
			NewEnPassant(E4, F3),
		},
		{
			"en passant declined (white)",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 0 3",
			NewMove(G1, F3),
		},
		{
			"white pawn capture black pawn",
			"rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2",
			"rnbqkbnr/pppp1ppp/8/4P3/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2",
			NewCapture(D4, E5),
		},
		{
			"white bishop capture black knight",
			"r1bqkb1r/pppppppp/2n2n2/1B6/8/4P3/PPPP1PPP/RNBQK1NR w KQkq - 0 3",
			"r1bqkb1r/pppppppp/2B2n2/8/8/4P3/PPPP1PPP/RNBQK1NR b KQkq - 0 3",
			NewCapture(B5, C6),
		},
		{
			"black queen capture white rook",
			"4k3/8/5q2/8/8/8/8/R3K3 b - - 0 123",
			"4k3/8/8/8/8/8/8/q3K3 w - - 0 124",
			NewCapture(F6, A1),
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
		{
			"moving king - can no longer castle (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R2K1BNR b kq - 0 5",
			NewMove(E1, D1),
		},
		{
			"moving queenside rook - can no longer castle queenside (white)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/R3KBNR w KQkq - 0 5",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPB3/PPPQPPPP/1R2KBNR b Kkq - 0 5",
			NewMove(A1, B1),
		},
		{
			"moving kingside rook - can no longer castle kingside (white)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK2R w KQkq - 0 4",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/3BPN2/PPPP1PPP/RNBQK1R1 b Qkq - 0 4",
			NewMove(H1, G1),
		},
		{
			"moving king - can no longer castle (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"r2k1bnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQ - 0 6",
			NewMove(E8, D8),
		},
		{
			"moving queenside rook - can no longer castle queenside (black)",
			"r3kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R b KQkq - 0 5",
			"1r2kbnr/pppqpppp/2npb3/8/8/2NPBN2/PPPQPPPP/R3KB1R w KQk - 0 6",
			NewMove(A8, B8),
		},
		{
			"moving kingside rook - can no longer castle kingside (black)",
			"rnbqk2r/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R b KQkq - 0 4",
			"rnbqk1r1/pppp1ppp/3bpn2/8/8/2NBPN2/PPPP1PPP/R1BQK2R w KQq - 0 5",
			NewMove(H8, G8),
		},
		{
			"capturing kingside rook - can no longer castle kingside (white to move)",
			"4k2r/8/6N1/8/8/8/8/2KR4 w k - 1 125",
			"4k2N/8/8/8/8/8/8/2KR4 b - - 1 125",
			NewMove(G6, H8),
		},
		// TODO: test capturing A8, A1, H1 for completeness
		{
			"promotion to queen (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"rQ2k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewQueenPromotion(B7, B8, false),
		},
		{
			"promotion to bishop (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"rB2k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewBishopPromotion(B7, B8, false),
		},
		{
			"promotion to rook (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"rR2k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewRookPromotion(B7, B8, false),
		},
		{
			"promotion to knight (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"rN2k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewKnightPromotion(B7, B8, false),
		},
		{
			"promotion to queen with capture (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"Q3k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewQueenPromotion(B7, A8, true),
		},
		{
			"promotion to bishop with capture  (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"B3k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewBishopPromotion(B7, A8, true),
		},
		{
			"promotion to rook with capture (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"R3k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewRookPromotion(B7, A8, true),
		},
		{
			"promotion to knight with capture (white)",
			"r3k3/1P6/8/8/8/8/8/4K3 w - - 0 123",
			"N3k3/8/8/8/8/8/8/4K3 b - - 0 123",
			NewKnightPromotion(B7, A8, true),
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
