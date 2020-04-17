package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"strings"
	"testing"
)

func TestMoves(t *testing.T) {
	moves := []struct {
		name     string
		board    string
		expected []string
	}{
		{
			"pawn pushes (white)",
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			[]string{
				// knights
				"b1a3", "b1c3",
				"g1f3", "g1h3",
				// pawns
				"a2a3", "a2a4",
				"b2b3", "b2b4",
				"c2c3", "c2c4",
				"d2d3", "d2d4",
				"e2e3", "e2e4",
				"f2f3", "f2f4",
				"g2g3", "g2g4",
				"h2h3", "h2h4",
			},
		},
		{
			"pawn pushes (black)",
			"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
			[]string{
				// knights
				"b8a6", "b8c6",
				"g8f6", "g8h6",
				// pawns
				"a7a6", "a7a5",
				"b7b6", "b7b5",
				"c7c6", "c7c5",
				"d7d6", "d7d5",
				"e7e6", "e7e5",
				"f7f6", "f7f5",
				"g7g6", "g7g5",
				"h7h6", "h7h5",
			},
		},
		{
			"pawns can't move if blocked by opposing piece",
			"k7/8/8/8/3p4/p5P1/P2Pp3/7K w - - 0 123",
			[]string{
				"h1g1", "h1h2", "h1g2", // king
				"d2d3", "g3g4", // pawns
			},
		},
		{
			"pawns can't move if blocked by friendly piece",
			"k7/5pp1/p5n1/1p3p2/8/2pp4/8/7K b - - 0 123",
			[]string{
				"a8a7", "a8b8", "a8b7", // king
				"g6e5", "g6e7", "g6f4", "g6f8", "g6h4", "g6h8", // knight
				"f7f6", "f5f4", "a6a5", "b5b4", "c3c2", "d3d2", // pawns
			},
		},
		{
			"pawn captures (white)",
			"4k3/8/8/p1bpp2p/P2PPPpP/r5PP/1P6/4K3 w - - 0 123",
			[]string{
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"b2xa3", "b2b3", "b2b4", // b2 pawn
				"d4xc5", "d4xe5", // d4 pawn
				"e4xd5",         // e4 pawn
				"f4f5", "f4xe5", // f4 pawn
				"h3xg4", // h3 pawn
			},
		},
		{
			"pawn captures (black)",
			"4k3/8/8/p2pp2p/P2PPPpP/R5PR/1P6/4K3 b - - 0 123",
			[]string{
				"e8e7", "e8f7", "e8f8", "e8d7", "e8d8", // king
				"d5xe4", "e5xd4", "e5xf4", "g4xh3", // pawns
			},
		},
		// TODO: update pawn promo tests to check for wrapping A => H or H => A file on captures
		{
			"pawn promotions (white)",
			"rn3rk1/P1PP4/4P3/5P2/8/8/8/4K3 w - - 0 123",
			[]string{
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"e6e7", "f5f6", // a6, f5 pawns can't promote
				// a7 pawn can only capture
				"a7xb8=Q", "a7xb8=N", "a7xb8=R", "a7xb8=B",
				// c7 pawn can either capture or advance
				"c7xb8=Q", "c7xb8=N", "c7xb8=R", "c7xb8=B",
				"c7c8=Q", "c7c8=N", "c7c8=R", "c7c8=B",
				// d7 pawn can only advance
				"d7d8=Q", "d7d8=N", "d7d8=R", "d7d8=B",
			},
		},
		{
			"pawn promotions (black)",
			"4k3/8/8/8/8/3P4/1pp4p/bN1QK2R b - - 0 123",
			[]string{
				"e8d7", "e8d8", "e8e7", "e8f7", "e8f8", // king
				// c2 pawn can either capture or advance
				"c2xb1=B", "c2xb1=N", "c2xb1=Q", "c2xb1=R",
				"c2c1=B", "c2c1=N", "c2c1=Q", "c2c1=R",
				"c2xd1=B", "c2xd1=N", "c2xd1=Q", "c2xd1=R"},
		},
		{
			"pawn en passant (black to move)",
			"4k3/8/8/8/3pPp2/3P2P1/8/4K3 b - e3 0 123",
			[]string{
				"e8d7", "e8d8", "e8e7", "e8f7", "e8f8", // king
				"d4xe3e.p.",                  // d4 pawn
				"f4xg3", "f4f3", "f4xe3e.p.", // f4 pawn
			},
		},
		{
			"pawn en passant (white to move)",
			"4k3/8/8/2PpP3/4P3/8/8/K7 w - d6 0 123",
			[]string{
				"a1b1", "a1a2", "a1b2", // king
				"e5xd6e.p.", "e5e6", // e5 pawn
				"c5xd6e.p.", "c5c6", // c5 pawn
				"e4xd5", // e4 pawn
			},
		},
		{
			"knight moves",
			"4k3/8/8/3p4/p7/2N5/P3P3/4K3 w - - 1 123",
			[]string{
				"c3b1", "c3b5", "c3d1", "c3e4", "c3xa4", "c3xd5", // knight
				"e1d1", "e1d2", "e1f1", "e1f2", // king
				"e2e3", "e2e4", "a2a3", // pawns
			},
		},
		{
			"rook moves",
			"k7/8/8/8/8/8/8/1R5K w - - 1 123",
			[]string{
				"h1h2", "h1g1", "h1g2", // king
				"b1b2", "b1b3", "b1b4", "b1b5", "b1b6", "b1b7", "b1b8", // rook vertical
				"b1a1", "b1c1", "b1d1", "b1e1", "b1f1", "b1g1", // rook horizontal
			},
		},
		{
			"rook captures",
			"k7/8/8/8/1q6/8/rR2P3/1N5K w - - 1 123",
			[]string{
				"b1d2", "b1c3", "b1a3", // knight
				"e2e3", "e2e4", // pawn
				"h1g1", "h1h2", "h1g2", // king
				"b2xa2", "b2xb4", "b2b3", "b2c2", "b2d2", // rook
			},
		},
		{
			"rook moves blocked", // two blockers in every direction; this test makes sure we get the block masks correct
			"2bqk3/1pppp3/8/P1pr1P1p/8/3P4/8/3QK3 b - - 0 123",
			[]string{
				"d5xd3", "d5d4", "d5d6", "d5e5", "d5xf5", // rook
				"e8f7", "e8f8", // king
				"c7c6", "d7d6", "e7e6", "e7e5", "b7b6", "b7b5", "c5c4", "h5h4", // pawns
			},
		},
		{
			"rook move edge conditions: rooks h1 and a8",
			"R7/1k6/8/8/8/8/6K1/7R w - - 1 123",
			[]string{
				"a8a1", "a8a2", "a8a3", "a8a4", "a8a5", "a8a6", "a8a7", // a8 rook vertical
				"a8b8", "a8c8", "a8d8", "a8e8", "a8f8", "a8g8", "a8h8", // a8 rook horizontal
				"h1a1", "h1b1", "h1c1", "h1d1", "h1e1", "h1f1", "h1g1", // h1 rook vertical
				"h1h2", "h1h3", "h1h4", "h1h5", "h1h6", "h1h7", "h1h8", // h1 rook horizontal
				"g2f1", "g2f2", "g2f3", "g2g1", "g2g3", "g2h2", "g2h3", // king
			},
		},
		{
			"rook move edge conditions: rooks a1 and h8",
			"7R/1k6/8/8/8/8/6K1/R7 w - - 1 123",
			[]string{
				"a1a2", "a1a3", "a1a4", "a1a5", "a1a6", "a1a7", "a1a8", // a1 rook vertical
				"a1b1", "a1c1", "a1d1", "a1e1", "a1f1", "a1g1", "a1h1", // a1 rook horizontal
				"h8h1", "h8h2", "h8h3", "h8h4", "h8h5", "h8h6", "h8h7", // h8 rook vertical
				"h8a8", "h8b8", "h8c8", "h8d8", "h8e8", "h8f8", "h8g8", // h8 rook horizontal
				"g2f1", "g2f2", "g2f3", "g2g1", "g2g3", "g2h1", "g2h2", "g2h3", // king
			},
		},
		{
			"bishop moves",
			"4k3/3b4/8/8/8/8/8/3K4 b - - 1 123",
			[]string{
				"e8e7", "e8f7", "e8f8", "e8d8", // king
				"d7c8", "d7c6", "d7b5", "d7a4", "d7e6", "d7f5", "d7g4", "d7h3", // bishop
			},
		},
		{
			"bishop captures",
			"4k3/3p4/p7/1B6/8/3K4/8/8 w - - 1 123",
			[]string{
				"d3c2", "d3c3", "d3c4", "d3d2", "d3d4", "d3e2", "d3e3", "d3e4", // king
				"b5c6", "b5c4", "b5a4", "b5xa6", "b5xd7", // bishop
			},
		},
		{
			"bishop moves blocked", // two blockers in every direction; this test makes sure we get the block masks correct
			"rn2k1n1/pP2p3/p3R3/3b4/8/1P3p2/P7/4K2R b - - 0 123",
			[]string{
				"g8f6", "g8h6", // kingside knight
				"e8d7", "e8d8", "e8f7", "e8f8", // king
				"a6a5", "f3f2", // pawns
				"b8c6", "b8d7", // queenside knight
				"d5c4", "d5e4", "d5c6", "d5xb3", "d5xb7", "d5xe6", // bishop
			},
		},
		{
			"bishop move edge conditions: bishops moving to corners of the board",
			"4k3/8/8/3B4/8/2B5/8/6K1 w - - 0 1",
			[]string{
				"c3a1", "c3b2", "c3d4", "c3e5", "c3f6", "c3g7", "c3h8", // c3 bishop rising diagonal
				"c3a5", "c3b4", "c3d2", "c3e1", // c3 bishop falling diagonal
				"d5a2", "d5b3", "d5c4", "d5e6", "d5f7", "d5g8", // d5 bishop rising diagonal
				"d5a8", "d5b7", "d5c6", "d5e4", "d5f3", "d5g2", "d5h1", // d5 bishop falling diagonal
				"g1f1", "g1f2", "g1g2", "g1h1", "g1h2", // king
			},
		},
		{
			"queen moves",
			"3k4/2p3P1/P2q3B/8/8/P7/7P/3BK3 b - - 1 123",
			[]string{
				"c7c5", "c7c6", // pawn
				"d8c8", "d8d7", "d8e7", "d8e8", // king
				"d6d7",         // queen north
				"d6f8", "d6e7", // queen north east
				"d6f6", "d6e6", "d6g6", "d6xh6", // queen east
				"d6f4", "d6e5", "d6g3", "d6xh2", // queen south east
				"d6d2", "d6d3", "d6d4", "d6d5", "d6xd1", // queen south
				"d6b4", "d6xa3", "d6c5", // queen south west
				"d6c6", "d6b6", "d6xa6", // queen west
			},
		},
		{
			"king must not move into check (pawns; black to move)",
			"4k3/6P1/2P1P3/8/8/8/8/4K3 b - - 0 123",
			[]string{
				"e8e7", "e8d8",
			},
		},
		{
			"king must not move into check (pawns; white to move)",
			"4k3/8/8/8/8/2p3p1/8/4K3 w - - 0 123",
			[]string{
				"e1e2", "e1d1", "e1f1",
			},
		},
		{
			"king must not move into check (bishop)",
			"4k3/8/8/3p4/4Kb2/8/8/8 w - - 1 123",
			[]string{
				"e4xd5", "e4f5",
				"e4d4", "e4xf4",
				"e4d3", "e4f3",
			},
		},
		{
			"king must not move into check (rook)",
			"4k2r/8/8/8/8/8/8/6K1 w - - 1 123",
			[]string{"g1g2", "g1f1", "g1f2"},
		},
		{
			"king must not move into check (knights)",
			"4k3/8/6N1/3N4/8/8/8/4K3 b - - 1 123",
			[]string{"e8d7", "e8d8", "e8f7"},
		},
		{
			"king must not move into check (opposing king)",
			"k7/2K5/7R/8/8/8/8/8 b - - 1 123",
			[]string{"a8a7"},
		},
		{
			"king must not move into check (pawns)",
			"4k3/8/8/8/8/3ppp2/8/4K3 w - - 0 123",
			[]string{"e1d1", "e1f1"},
		},
		// TODO: test for queen threat
		{
			"king must not move into check: stalemate (no moves possible)",
			"4k1r1/8/8/8/8/8/r7/7K w - - 1 123",
			nil,
		},
		{
			"king must not move into check: free to move, own piece blocks check",
			"3qk3/8/8/8/8/8/3P4/3K4 w - - 1 123",
			[]string{
				"d2d3", "d2d4", // pawn
				"d1c1", "d1e1", "d1c2", "d1e2", // king
			},
		},
		{
			"king must not move into check: free to move, opposing piece blocks check",
			"3qk3/3b4/8/8/8/8/8/3K4 w - - 1 123",
			[]string{
				"d1c1", "d1e1", "d1d2", "d1c2", "d1e2", // king
			},
		},
		{
			"king must not move into check: king may not capture a covered piece",
			"4k3/5B2/6P1/8/8/8/8/4K2R b K - 0 123",
			[]string{
				"e8f8", "e8e7", "e8d7", "e8d8",
			},
		},
		{
			"clearing check: must capture to clear check",
			"r1b1k3/1P6/8/8/4n3/6P1/2nPP2P/R2QKBN1 w - - 0 123",
			[]string{
				"d1xc2", // queen must capture knight
			},
		},
		{
			"clearing check: pawn must capture ne to clear check",
			"4k3/8/8/8/1b6/P7/4PP2/3BKB2 w - - 0 123",
			[]string{
				"a3xb4", // pawn must capture bishop
			},
		},
		// TODO: pawn must capture nw to clear check
		// TODO: pawn must capture se to clear check
		// TODO: pawn must capture sw to clear check
		{
			"clearing check: king must capture to clear check",
			"3qkb2/3ppP2/2PP4/8/8/8/8/4K2R b K - 0 123",
			[]string{
				"e8xf7", // king must capture pawn to clear check
			},
		},
		{
			"clearing check: piece must block to clear check (bishop)",
			"4k3/8/8/8/1b6/8/4PP2/1N1BKB2 w - - 0 123",
			[]string{
				"b1c3", "b1d2", // knight must sacrifice itself
			},
		},
		{
			"clearing check: piece must block to clear check (rook)",
			"R3kb2/2brpp2/8/8/8/8/8/4K3 b - - 0 123",
			[]string{
				"c7b8", "c7d8", "d7d8", // either bishop or rook must sacrifice itself
			},
		},
		{
			// This test exposes a bug in treating "threat" rays the same as
			// "block" rays. F1 - the square east of the king - is not a valid
			// square for the king to move to, so it must be included in the
			// mask of "threat" rays. Moving the queen to that square does NOT
			// clear the king from check however, and so is not a legal move.
			"clearing check: moving behind king (MSB) doesn't clear check",
			"4k3/8/8/8/8/3Q4/8/q3K3 w - - 0 123",
			[]string{
				"d3b1", "d3d1", // queen must move in between queen and king
				"e1d2", "e1e2", "e1f2", // king may move up a rank
			},
		},
		{
			// As above this test makes sure that we're treating "threat" rays
			// (squares the king may not move to) separately to "block" rays
			// (squares that we may move a piece to in order to clear check).
			// except now the bit for the square we must move to is more
			// significant than the bit the king is occupying. D1 must be marked
			// as a "threatened" (king may not move to square), but ♘D1 is not a
			// valid move.
			"clearing check: moving behind king (LSB) doesn't clear check",
			"4k3/8/8/8/8/4N3/r7/4K2q w - - 0 123",
			[]string{
				"e3f1", // knight must move in between king and queen
			},
		},
		// TODO: pawn must push (single or double) to clear check (black)
		// TODO: pawn must push (single or double) to clear check (white)
		// TODO: pawn must promote to clear check (black)
		// TODO: pawn must promote to clear check (white)
		{
			"clearing check: capturing or blocking piece doesn't work if double check; king must move",
			"4k3/1pp5/2B5/1b3n2/8/3r1p2/4R2r/4K3 b - - 0 123",
			[]string{
				"e8d8", "e8f7", "e8f8", // only valid moves are king moves
			},
		},
		// TODO: en passant to take checking piece https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/#gotcha-en-passant-check-evasions
		// TODO: en passant to block check
		{
			// https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/#gotcha-king-moves-away-from-a-checking-slider
			"clearing check: king may not move away from checking slider while still on ray",
			"4k3/1q6/8/8/2r1K3/8/8/8 w - - 0 1",
			[]string{
				"e4e5", "e4e3", "e4f5", "e4d3",
			},
		},
		{
			"pinning: absolutely pinned piece must stay on ray (bishop SW/NE diagonal)",
			"4k3/8/2b5/8/B7/8/8/4K3 b - - 0 1",
			[]string{
				"c6b5", "c6d7", "c6xa4", // bishop pinned to SW/NE diagonal
				"e8e7", "e8d7", "e8d8", "e8f7", "e8f8", // king
			},
		},
		{
			"pinning: absolutely pinned piece must stay on ray (bishop NW/SE diagonal)",
			"4k3/8/8/8/1b6/2B5/8/4K3 w - - 0 1",
			[]string{
				"c3xb4", "c3d2", // bishop pinned to NW/SE diagonal
				"e1e2", "e1d1", "e1d2", "e1f1", "e1f2", // king
			},
		},
		{
			"pinning: absolutely pinned piece must stay on ray (rook vertical)",
			"4k3/4r3/8/8/8/8/4R3/4K3 w - - 0 1",
			[]string{
				"e2e3", "e2e4", "e2e5", "e2e6", "e2xe7", // rook pinned to vertical
				"e1d1", "e1d2", "e1f1", "e1f2", // king
			},
		},
		{
			"pinning: absolutely pinned piece must stay on ray (rook horizontal)",
			"Rr2k3/8/8/8/8/8/8/4K3 b - - 0 1",
			[]string{
				"b8c8", "b8d8", "b8xa8", // rook pinned to horizontal
				"e8e7", "e8d7", "e8d8", "e8f7", "e8f8", // king
			},
		},
		// TODO: may not en passant if that exposes king https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/#gotcha-en-passant-discovered-check
		// TODO: checkmate
		{
			"castling: white can ks castle only (qs blocked)",
			"4k3/8/8/8/8/8/P6P/R3K1NR w KQ - 1 123",
			[]string{
				"O-O-O",
				"a1b1", "a1c1", "a1d1", // queenside rook
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"g1e2", "g1f3", "g1h3", // knight
				"a2a3", "a2a4", "h2h3", "h2h4", // pawns
			},
		},
		{
			"castling: white can castle both (ks pawn protects from check, qs despite rook being covered)",
			"4kr2/8/8/8/8/r4P2/7P/R3K2R w KQ - 1 123",
			[]string{
				"O-O", "O-O-O",
				"a1a2", "a1xa3", "a1b1", "a1c1", "a1d1", // queenside rook
				"e1d1", "e1d2", "e1e2", "e1f1", "e1f2", // king
				"h1g1", "h1f1", // kingside rook
				"f3f4", "h2h3", "h2h4", // pawns
			},
		},
		{
			"castling: white can castle neither (ks pass through check, qs into check)",
			"2r1kr2/8/8/8/8/8/P6P/R3K2R w KQ - 1 123",
			[]string{
				"a1b1", "a1c1", "a1d1", // queenside rook
				"e1d1", "e1d2", "e1e2", // king
				"h1g1", "h1f1", // kingside rook
				"a2a3", "a2a4", "h2h3", "h2h4", // queenside rook
			},
		},
		{
			"castling: black can qs castle only (board state)",
			"r3k2r/p6p/8/8/8/8/8/4K3 b q - 1 123",
			[]string{
				"O-O-O",
				"a8b8", "a8c8", "a8d8", // queenside rook
				"h8g8", "h8f8", // kingside rook
				"a7a6", "a7a5", "h7h6", "h7h5", // pawns
				"e8d7", "e8d8", "e8e7", "e8f7", "e8f8", // king
			},
		},
		{
			"castling: black can castle neither (king in check)",
			"r3k2r/p6p/3N4/8/8/8/8/4K3 b kq - 0 1",
			[]string{
				// king must move to clear check; no other piece can move
				"e8d7", "e8d8", "e8e7", "e8f8", // king
			},
		},
		{
			"castling: black can ks castle only (qs into check)",
			"r3k2r/p1R4p/8/8/8/8/8/4K3 b kq - 0 123",
			[]string{
				"O-O",
				"a7a5", "a7a6", // queenside pawn
				"h7h5", "h7h6", // kingside pawn
				"a8b8", "a8c8", "a8d8", // queenside rook
				"h8f8", "h8g8", // kingside rook
				"e8d8", "e8f8", // king
			},
		},
	}

	for _, tt := range moves {
		t.Run(tt.name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.board))
			require.NoError(t, err)
			require.NotNil(t, b)
			var moves []string
			for _, move := range b.GenerateMoves(nil) {
				moves = append(moves, move.SAN())
			}
			// sort so we don't need to fiddle with ordering in the test case
			sort.Strings(tt.expected)
			sort.Strings(moves)
			assert.Equal(t, tt.expected, moves)
		})
	}
}

func TestTooManyCheckersPanics(t *testing.T) {
	fen := "4k3/4r3/8/q7/7b/8/8/4K3 w - - 0 123" // 3 checkers
	b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
	require.NoError(t, err)
	require.NotNil(t, b)
	assert.Panics(t, func() { _ = b.GenerateMoves(nil) })
}

func BenchmarkGenerateMoves10(b *testing.B) {
	const fen = "r3k2r/pbqnbppp/1p2pn2/2p1N3/Q1P5/4P3/PB1PBPPP/RN3RK1 w kq - 8 11" // 10 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateMoves(moves)
	}
}

func BenchmarkGenerateMoves20(b *testing.B) {
	const fen = "4rrk1/2qn2pp/pp2pb2/2p2p2/P1P2P2/2NPPR2/1BQ3PP/1R4K1 b - - 0 20" // 20 ply in, black to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateMoves(moves)
	}
}

func BenchmarkGenerateMoves30(b *testing.B) {
	const fen = "3rr1k1/1nq4p/pp4p1/2pP1p2/P4P2/2Q1P3/1R2N1PP/3R2K1 w - - 0 31" // 30 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateMoves(moves)
	}
}

func BenchmarkGenerateMoves40(b *testing.B) {
	const fen = "3r2k1/1n5p/8/pPpq1p1p/5P2/4P3/6PK/1R2QN2 b - - 3 40" // 40 ply in, black to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateMoves(moves)
	}
}

func BenchmarkGenerateMoves50(b *testing.B) {
	const fen = "6k1/1n5p/8/p7/2p2PP1/1r2P1N1/8/R5K1 w - - 3 51" // 50 ply in, white to play
	board, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	moves := make([]engine.Move, 0, 64)
	for i := 0; i < b.N; i++ {
		board.GenerateMoves(moves)
	}
}
