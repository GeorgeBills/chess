package engine

// Piece represents a chess piece.
type Piece byte

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
)

// Rune returns a rune that uniquely represents the piece colour and type.
func (p Piece) Rune() rune {
	switch {
	case p == PieceNone:
		return '□'
	case p&PieceWhite != 0 && p&PiecePawn != 0:
		return '♙'
	case p&PieceWhite != 0 && p&PieceRook != 0:
		return '♖'
	case p&PieceWhite != 0 && p&PieceBishop != 0:
		return '♗'
	case p&PieceWhite != 0 && p&PieceKnight != 0:
		return '♘'
	case p&PieceWhite != 0 && p&PieceKing != 0:
		return '♔'
	case p&PieceWhite != 0 && p&PieceQueen != 0:
		return '♕'
	case p&PieceBlack != 0 && p&PiecePawn != 0:
		return '♟'
	case p&PieceBlack != 0 && p&PieceRook != 0:
		return '♜'
	case p&PieceBlack != 0 && p&PieceBishop != 0:
		return '♝'
	case p&PieceBlack != 0 && p&PieceKnight != 0:
		return '♞'
	case p&PieceBlack != 0 && p&PieceKing != 0:
		return '♚'
	case p&PieceBlack != 0 && p&PieceQueen != 0:
		return '♛'
	default:
		panic(p) // invalid piece
	}
}
