package engine

// Colour is used to represent each colour.
type Colour byte

// White and Black are constants defined for the colours.
const (
	White Colour = 'w'
	Black Colour = 'b'
)

// A1...H8 are constants defined for hardcoding an index by its rank and file.
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

const (
	rank1 = iota
	rank2
	rank3
	rank4
	rank5
	rank6
	rank7
	rank8
)

const (
	fileA = iota
	fileB
	fileC
	fileD
	fileE
	fileF
	fileG
	fileH
)

const (
	maskAll  = 0xFFFFFFFFFFFFFFFF
	maskNone = 0x0000000000000000
)

const (
	maskRank1 uint64 = 0x00000000000000FF
	maskRank2 uint64 = 0x000000000000FF00
	maskRank3 uint64 = 0x0000000000FF0000
	maskRank4 uint64 = 0x00000000FF000000
	maskRank5 uint64 = 0x000000FF00000000
	maskRank6 uint64 = 0x0000FF0000000000
	maskRank7 uint64 = 0x00FF000000000000
	maskRank8 uint64 = 0xFF00000000000000
)

const (
	maskFileA = 1<<A1 | 1<<A2 | 1<<A3 | 1<<A4 | 1<<A5 | 1<<A6 | 1<<A7 | 1<<A8
	maskFileH = 1<<H1 | 1<<H2 | 1<<H3 | 1<<H4 | 1<<H5 | 1<<H6 | 1<<H7 | 1<<H8
)
