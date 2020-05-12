package engine_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	board := engine.NewBoard()

	assert.EqualValues(t, White, board.ToMove())

	assert.True(t, board.CanWhiteCastleKingside())
	assert.True(t, board.CanWhiteCastleQueenside())
	assert.True(t, board.CanBlackCastleKingside())
	assert.True(t, board.CanBlackCastleQueenside())

	assert.Equal(t, 0, board.HalfMoves())
	assert.Equal(t, 1, board.FullMoves())

	assert.EqualValues(t, math.MaxUint8, board.EnPassant())

	// expected pieces (other than pawns)
	pieces := []struct {
		idx   uint8
		piece engine.Piece
	}{
		{A1, PieceWhiteRook},
		{B1, PieceWhiteKnight},
		{C1, PieceWhiteBishop},
		{D1, PieceWhiteQueen},
		{E1, PieceWhiteKing},
		{F1, PieceWhiteBishop},
		{G1, PieceWhiteKnight},
		{H1, PieceWhiteRook},
		{A8, PieceBlackRook},
		{B8, PieceBlackKnight},
		{C8, PieceBlackBishop},
		{D8, PieceBlackQueen},
		{E8, PieceBlackKing},
		{F8, PieceBlackBishop},
		{G8, PieceBlackKnight},
		{H8, PieceBlackRook},
	}
	for _, tt := range pieces {
		t.Run(fmt.Sprintf("%c @ %d", tt.piece.Rune(), tt.idx), func(t *testing.T) {
			assert.Equal(t, tt.piece, board.PieceAt(tt.idx))
		})
	}

	// white pawns
	for _, idx := range []uint8{A2, B2, C2, D2, E2, F2, G2, H2} {
		t.Run(fmt.Sprintf("♙ @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceWhitePawn, board.PieceAt(idx))
		})
	}

	// black pawns
	for _, idx := range []uint8{A7, B7, C7, D7, E7, F7, G7, H7} {
		t.Run(fmt.Sprintf("♟ @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceBlackPawn, board.PieceAt(idx))
		})
	}

	// empty
	empty := []uint8{
		A3, B3, C3, D3, E3, F3, G3, H3,
		A4, B4, C4, D4, E4, F4, G4, H4,
		A5, B5, C5, D5, E5, F5, G5, H5,
		A6, B6, C6, D6, E6, F6, G6, H6,
	}
	for _, idx := range empty {
		t.Run(fmt.Sprintf("Empty @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceNone, board.PieceAt(idx))
		})
	}

	// to string
	str := board.String()
	const expected = `♜♞♝♛♚♝♞♜
♟♟♟♟♟♟♟♟
□□□□□□□□
□□□□□□□□
□□□□□□□□
□□□□□□□□
♙♙♙♙♙♙♙♙
♖♘♗♕♔♗♘♖`
	assert.Equal(t, expected, str)
}

func TestBoardColourFlipped(t *testing.T) {
	const fen = "4k3/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQ - 0 1"
	const expected = "rnbqkbnr/pppppppp/8/8/8/8/8/4K3 w kq - 0 1"
	b, err := NewBoardFromFEN(strings.NewReader(fen))
	require.NoError(t, err)
	require.NotNil(t, b)
	flipped := b.ColourFlipped()
	assert.Equal(t, expected, flipped.FEN())
}
