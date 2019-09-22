namespace Chess.Engine
{
    public static class Extensions
    {
        public static char ToChar(this Piece piece)
        {
            if (piece == (Piece.White & Piece.King))
                return '♔';
            if (piece == (Piece.White & Piece.Queen))
                return '♕';
            if (piece == (Piece.White & Piece.Rook))
                return '♖';
            if (piece == (Piece.White & Piece.Bishop))
                return '♗';
            if (piece == (Piece.White & Piece.Knight))
                return '♘';
            if (piece == (Piece.White & Piece.Pawn))
                return '♙';
            if (piece == (Piece.Black & Piece.King))
                return '♚';
            if (piece == (Piece.Black & Piece.Queen))
                return '♛';
            if (piece == (Piece.Black & Piece.Rook))
                return '♜';
            if (piece == (Piece.Black & Piece.Bishop))
                return '♝';
            if (piece == (Piece.Black & Piece.Knight))
                return '♞';
            if (piece == (Piece.Black & Piece.Pawn))
                return '♟';
            return ' ';
        }
    }
}
