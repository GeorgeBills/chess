package engine_test

import (
	"fmt"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardPieceAt(t *testing.T) {
	expected := []struct {
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
	board := engine.NewBoard()
	for _, tt := range expected {
		t.Run(fmt.Sprintf("%c @ %d", tt.piece.Rune(), tt.idx), func(t *testing.T) {
			assert.Equal(t, tt.piece, board.PieceAt(tt.idx))
		})
	}
}

func TestNewBoardPawnAt(t *testing.T) {
	board := engine.NewBoard()
	for _, idx := range []uint8{A2, B2, C2, D2, E2, F2, G2, H2} {
		t.Run(fmt.Sprintf("♙ @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceWhitePawn, board.PieceAt(idx))
		})
	}
	for _, idx := range []uint8{A7, B7, C7, D7, E7, F7, G7, H7} {
		t.Run(fmt.Sprintf("♟ @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceBlackPawn, board.PieceAt(idx))
		})
	}
}

func TestNewBoardEmptyAt(t *testing.T) {
	empty := []uint8{
		A3, B3, C3, D3, E3, F3, G3, H3,
		A4, B4, C4, D4, E4, F4, G4, H4,
		A5, B5, C5, D5, E5, F5, G5, H5,
		A6, B6, C6, D6, E6, F6, G6, H6,
	}
	board := engine.NewBoard()
	for _, idx := range empty {
		t.Run(fmt.Sprintf("Empty @ %d", idx), func(t *testing.T) {
			assert.Equal(t, PieceNone, board.PieceAt(idx))
		})
	}
}

func TestNewBoardToString(t *testing.T) {
	str := engine.NewBoard().String()
	expected := `♜♞♝♛♚♝♞♜
♟♟♟♟♟♟♟♟
□□□□□□□□
□□□□□□□□
□□□□□□□□
□□□□□□□□
♙♙♙♙♙♙♙♙
♖♘♗♕♔♗♘♖`
	assert.Equal(t, expected, str)
}
