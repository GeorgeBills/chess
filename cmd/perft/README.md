# Perft

`perft` runs a performance test as described at
https://www.chessprogramming.org/Perft.

It supports:
 
 * specifying [FEN](https://en.wikipedia.org/wiki/Forsyth-Edwards_Notation) and
   depth
 * toggling "divide" (useful for recursively debugging faulty positions by
   comparing results with another engine)
 * toggling "validate" (validates every position per the engine
   `board.Validate()` method, which is slow but can help debug an illegal move
   by failing out early)

```
$ .\perft.exe -depth 7 
3195901860 nodes, took 2m37.630801s
$ .\perft.exe -depth 6 -fen 'r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 1 123'
8031647685 nodes, took 6m47.4479226s
```
