# Perft

`perft` runs a performance test as described at
https://www.chessprogramming.org/Perft. It supports specifying FEN and depth;
toggling "divide" (useful for recursively debugging faulty positions by
comparing results with another engine); and toggling "validate" (validates every
position per the engine `board.Validate()` method, which is slow but can help
debug a faulty position by failing out early).

```
PS C:\Users\georg.GEORGEB-XPS15\Documents\repositories\chess> go run .\cmd\perft\ -depth 7
3195901860 nodes, 153467ms
PS C:\Users\georg.GEORGEB-XPS15\Documents\repositories\chess> go run .\cmd\perft\ -depth 6 -fen 'r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 1 123'
8031647685 nodes, 418629ms
```
