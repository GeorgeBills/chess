using System;

namespace Chess.Engine
{
    [Flags]
    public enum Piece : byte
    {
        None  =  0b00000000,
        White =  0b10000000,
        Black =  0b01000000,
        Pawn =   0b00100000,
        Knight = 0b00010000,
        Bishop = 0b00001000,
        Rook =   0b00000100,
        Queen =  0b00000010,
        King =   0b00000001,
    }
}
