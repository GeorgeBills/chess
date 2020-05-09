# an-to-bb

`an-to-bb` takes a list of squares in [algebraic
notation](https://en.wikipedia.org/wiki/Algebraic_notation_(chess)) and outputs
a [bitboard](https://en.wikipedia.org/wiki/Bitboard) variable with the
corresponding bits set. It's useful mainly for generating constants to copy
paste into code, and for passing to `bb-to-visual`.

```
$ .\an-to-bb.exe g8 h7 g7
0b01000000_11000000_00000000_00000000_00000000_00000000_00000000_00000000
```
