package engine

import (
	"fmt"
)

// Piece represents a chess piece.
type Piece byte

// TODO: seems like it's better to have these as indexes into an array of boards
//       doing that would mean we'd never need to switch on piece type, just index

// Pieces will have a bit set for the colour and a bit set for the type.
const (
	PieceNone   Piece = 0b00000000
	PieceWhite  Piece = 0b10000000
	PieceBlack  Piece = 0b01000000
	PiecePawn   Piece = 0b00100000
	PieceKnight Piece = 0b00010000
	PieceBishop Piece = 0b00001000
	PieceRook   Piece = 0b00000100
	PieceQueen  Piece = 0b00000010
	PieceKing   Piece = 0b00000001

	PieceWhitePawn   Piece = PieceWhite | PiecePawn
	PieceWhiteKnight Piece = PieceWhite | PieceKnight
	PieceWhiteBishop Piece = PieceWhite | PieceBishop
	PieceWhiteRook   Piece = PieceWhite | PieceRook
	PieceWhiteQueen  Piece = PieceWhite | PieceQueen
	PieceWhiteKing   Piece = PieceWhite | PieceKing

	PieceBlackPawn   Piece = PieceBlack | PiecePawn
	PieceBlackKnight Piece = PieceBlack | PieceKnight
	PieceBlackBishop Piece = PieceBlack | PieceBishop
	PieceBlackRook   Piece = PieceBlack | PieceRook
	PieceBlackQueen  Piece = PieceBlack | PieceQueen
	PieceBlackKing   Piece = PieceBlack | PieceKing
)

// Rune returns a rune that uniquely represents the piece colour and type.
func (p Piece) Rune() rune {
	switch p {
	case PieceNone:
		return '□'
	case PieceWhitePawn:
		return '♙'
	case PieceWhiteRook:
		return '♖'
	case PieceWhiteBishop:
		return '♗'
	case PieceWhiteKnight:
		return '♘'
	case PieceWhiteKing:
		return '♔'
	case PieceWhiteQueen:
		return '♕'
	case PieceBlackPawn:
		return '♟'
	case PieceBlackRook:
		return '♜'
	case PieceBlackBishop:
		return '♝'
	case PieceBlackKnight:
		return '♞'
	case PieceBlackKing:
		return '♚'
	case PieceBlackQueen:
		return '♛'
	default:
		panic(fmt.Errorf("invalid piece while generating rune: %b", p))
	}
}
