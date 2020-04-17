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
	maskAll  uint64 = 0xFFFFFFFFFFFFFFFF
	maskNone uint64 = 0x0000000000000000
)

const (
	maskRank1 uint64 = 1<<A1 | 1<<B1 | 1<<C1 | 1<<D1 | 1<<E1 | 1<<F1 | 1<<G1 | 1<<H1
	maskRank2 uint64 = 1<<A2 | 1<<B2 | 1<<C2 | 1<<D2 | 1<<E2 | 1<<F2 | 1<<G2 | 1<<H2
	maskRank3 uint64 = 1<<A3 | 1<<B3 | 1<<C3 | 1<<D3 | 1<<E3 | 1<<F3 | 1<<G3 | 1<<H3
	maskRank4 uint64 = 1<<A4 | 1<<B4 | 1<<C4 | 1<<D4 | 1<<E4 | 1<<F4 | 1<<G4 | 1<<H4
	maskRank5 uint64 = 1<<A5 | 1<<B5 | 1<<C5 | 1<<D5 | 1<<E5 | 1<<F5 | 1<<G5 | 1<<H5
	maskRank6 uint64 = 1<<A6 | 1<<B6 | 1<<C6 | 1<<D6 | 1<<E6 | 1<<F6 | 1<<G6 | 1<<H6
	maskRank7 uint64 = 1<<A7 | 1<<B7 | 1<<C7 | 1<<D7 | 1<<E7 | 1<<F7 | 1<<G7 | 1<<H7
	maskRank8 uint64 = 1<<A8 | 1<<B8 | 1<<C8 | 1<<D8 | 1<<E8 | 1<<F8 | 1<<G8 | 1<<H8
)

const (
	maskFileA uint64 = 1<<A1 | 1<<A2 | 1<<A3 | 1<<A4 | 1<<A5 | 1<<A6 | 1<<A7 | 1<<A8
	maskFileB uint64 = 1<<B1 | 1<<B2 | 1<<B3 | 1<<B4 | 1<<B5 | 1<<B6 | 1<<B7 | 1<<B8
	maskFileC uint64 = 1<<C1 | 1<<C2 | 1<<C3 | 1<<C4 | 1<<C5 | 1<<C6 | 1<<C7 | 1<<C8
	maskFileD uint64 = 1<<D1 | 1<<D2 | 1<<D3 | 1<<D4 | 1<<D5 | 1<<D6 | 1<<D7 | 1<<D8
	maskFileE uint64 = 1<<E1 | 1<<E2 | 1<<E3 | 1<<E4 | 1<<E5 | 1<<E6 | 1<<E7 | 1<<E8
	maskFileF uint64 = 1<<F1 | 1<<F2 | 1<<F3 | 1<<F4 | 1<<F5 | 1<<F6 | 1<<F7 | 1<<F8
	maskFileG uint64 = 1<<G1 | 1<<G2 | 1<<G3 | 1<<G4 | 1<<G5 | 1<<G6 | 1<<G7 | 1<<G8
	maskFileH uint64 = 1<<H1 | 1<<H2 | 1<<H3 | 1<<H4 | 1<<H5 | 1<<H6 | 1<<H7 | 1<<H8
)
