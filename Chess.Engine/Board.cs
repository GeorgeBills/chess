using System;
using System.Text;

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

        public bool IsWhiteAt(uint i) => ((long)white & (1 << (int)i)) != 0;
        public bool IsBlackAt(uint i) => ((long)black & (1 << (int)i)) != 0;
        public bool IsEmptyAt(uint i) => !IsWhiteAt(i) && !IsBlackAt(i);
        public bool IsPawnAt(uint i) => ((long)pawns & (1 << (int)i)) != 0;
        public bool IsKnightAt(uint i) => ((long)knights & (1 << (int)i)) != 0;
        public bool IsBishopAt(uint i) => ((long)bishops & (1 << (int)i)) != 0;
        public bool IsRookAt(uint i) => ((long)rooks & (1 << (int)i)) != 0;
        public bool IsQueenAt(uint i) => ((long)queens & (1 << (int)i)) != 0;
        public bool IsKingAt(uint i) => ((long)kings & (1 << (int)i)) != 0;

        public Piece PieceAt(uint i)
        {
            if (IsWhiteAt(i))
            {
                if (IsPawnAt(i))
                    return Piece.White & Piece.Pawn;
                if (IsKnightAt(i))
                    return Piece.White & Piece.Knight;
                if (IsBishopAt(i))
                    return Piece.White & Piece.Bishop;
                if (IsRookAt(i))
                    return Piece.White & Piece.Rook;
                if (IsQueenAt(i))
                    return Piece.White & Piece.Queen;
                if (IsKingAt(i))
                    return Piece.White & Piece.King;
            }
            if (IsBlackAt(i))
            {
                if (IsPawnAt(i))
                    return Piece.Black & Piece.Pawn;
                if (IsKnightAt(i))
                    return Piece.Black & Piece.Knight;
                if (IsBishopAt(i))
                    return Piece.Black & Piece.Bishop;
                if (IsRookAt(i))
                    return Piece.Black & Piece.Rook;
                if (IsQueenAt(i))
                    return Piece.Black & Piece.Queen;
                if (IsKingAt(i))
                    return Piece.Black & Piece.King;
            }
            return Piece.None;
        }

        public override string ToString()
        {
            var str = new StringBuilder();
            for (uint i = 0; i < 64; i++)
            {
                var ch = this.PieceAt(i).ToChar();
                str.Append(ch);
                if (i % 8 == 0) {
                    str.Append(Environment.NewLine);
                }
            }
            return str.ToString();
        }
    }
}
