using System;

namespace Chess.Engine
{
    public class Board
    {
        public static readonly Board Initial = new Board
        {
            white =   0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_11111111,
            black =   0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000,
            pawns =   0b00000000_11111111_00000000_00000000_00000000_00000000_11111111_00000000,
            knights = 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_01000010,
            bishops = 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00100100,
            rooks =   0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001,
            queens =  0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00010000,
            kings =   0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00001000,
        };

        private ulong white;
        private ulong black;
        private ulong pawns;
        private ulong knights;
        private ulong bishops;
        private ulong rooks;
        private ulong queens;
        private ulong kings;
    }
}
