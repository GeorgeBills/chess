package engine_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorWriter struct{}

func (*errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("error writing")
}

func TestNewBoardToFEN(t *testing.T) {
	fen := engine.NewBoard().FEN()
	assert.Equal(t, engine.InitialBoardFEN, fen)
}

func TestNewBoardFromFEN(t *testing.T) {
	expected := engine.NewBoard()
	b, err := engine.NewBoardFromFEN(
		strings.NewReader(engine.InitialBoardFEN),
	)
	require.NoError(t, err)
	assert.Equal(t, expected, b)
}

func TestEnPassant(t *testing.T) {
	fen := []struct {
		fen      string
		expected uint8
	}{
		{"rnbqkbnr/pppppppp/8/8/8/7N/PPPPPPPP/RNBQKB1R b KQkq - 1 1", math.MaxUint8},
		// white pawns
		{"rnbqkbnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQKBNR b KQkq A3 0 1", engine.A3},
		{"rnbqkbnr/pppppppp/8/8/1P6/8/P1PPPPPP/RNBQKBNR b KQkq b3 0 1", engine.B3},
		{"rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR b KQkq C3 0 1", engine.C3},
		{"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1", engine.D3},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq E3 0 1", engine.E3},
		{"rnbqkbnr/pppppppp/8/8/5P2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 1", engine.F3},
		{"rnbqkbnr/pppppppp/8/8/6P1/8/PPPPPP1P/RNBQKBNR b KQkq G3 0 1", engine.G3},
		{"rnbqkbnr/pppppppp/8/8/P7/8/PPPPPPP1/RNBQKBNR b KQkq h3 0 1", engine.H3},
		// black pawns
		{"rnbqkbnr/1ppppppp/8/p7/8/7N/PPPPPPPP/RNBQKB1R w KQkq a6 0 2", engine.A6},
		{"rnbqkbnr/p1pppppp/8/1p6/8/7N/PPPPPPPP/RNBQKB1R w KQkq B6 0 2", engine.B6},
		{"rnbqkbnr/pp1ppppp/8/2p5/8/7N/PPPPPPPP/RNBQKB1R w KQkq c6 0 2", engine.C6},
		{"rnbqkbnr/ppp1pppp/8/3p4/8/7N/PPPPPPPP/RNBQKB1R w KQkq D6 0 2", engine.D6},
		{"rnbqkbnr/pppp1ppp/8/4p3/8/7N/PPPPPPPP/RNBQKB1R w KQkq e6 0 2", engine.E6},
		{"rnbqkbnr/ppppp1pp/8/5p2/8/7N/PPPPPPPP/RNBQKB1R w KQkq F6 0 2", engine.F6},
		{"rnbqkbnr/pppppp1p/8/6p1/8/7N/PPPPPPPP/RNBQKB1R w KQkq g6 0 2", engine.G6},
		{"rnbqkbnr/ppppppp1/8/7p/8/7N/PPPPPPPP/RNBQKB1R w KQkq H6 0 2", engine.H6},
	}
	for i, tt := range fen {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.fen))
			require.NoError(t, err)
			require.NotNil(t, b)
			ep := b.EnPassant()
			assert.Equal(t, tt.expected, ep)
		})
	}
}

func TestNewBoardFromInvalidFEN(t *testing.T) {
	var tests map[string]struct {
		FEN   string
		Error string
	}

	f, err := os.Open("testdata/invalid-fen.json")
	require.NoError(t, err)
	err = json.NewDecoder(f).Decode(&tests)
	require.NoError(t, err)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt.FEN))
			assert.EqualError(t, err, tt.Error)
			assert.Nil(t, b)
		})
	}
}

func TestRoundTripFEN(t *testing.T) {
	fen := []string{
		// randomish FEN strings
		// http://bernd.bplaced.net/fengenerator/fengenerator.html
		"2Q5/8/NP2P1pP/r4BP1/2p2p2/4p3/1P2rk2/3K4 w - - 0 1",
		"1nb3R1/3P2Q1/p5p1/1P5P/1P2K3/1kN4B/1p6/7R w - - 0 1",
		"5k2/Q4p2/1P1Pp3/1p2pPNr/2n3p1/6p1/8/4K2b w - - 0 1",
		"4QK2/R4R1p/6r1/3P1n1N/2P2p2/2P4P/n6k/2N5 w - - 0 1",
		"1b6/1R1N4/1P4b1/p6P/3K1ppp/P3ppkp/8/8 w - - 0 1",
		"1k2n3/4B1r1/3rN2P/8/3qpp2/1n1P1p2/3p1P2/7K b - - 0 1",
		"3q4/1pb4p/8/P5pk/1pP4N/P4p2/2K3P1/1r3B2 b - - 0 1",
		"2k1nq2/8/8/pP2N2Q/3K4/2pp1r1P/3P2p1/B3R3 b - - 0 1",
		"6n1/1r6/k2b3p/P5pp/4PP2/5PKP/1PN1N3/8 b - - 0 1",
		"5N2/8/1b1p2r1/8/2PP4/2kp1p1K/1p2n2B/q3r1N1 b - - 0 1",
		// FEN with En Passant indicated
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
	}
	for i, tt := range fen {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			b, err := engine.NewBoardFromFEN(strings.NewReader(tt))
			require.NoError(t, err)
			require.NotNil(t, b)
			assert.Equal(t, tt, b.FEN())
		})
	}
}

func TestWriteFENError(t *testing.T) {
	err := engine.NewBoard().WriteFEN(&errorWriter{})
	assert.EqualError(t, err, "error writing")
}

func BenchmarkNewBoardFromFEN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		engine.NewBoardFromFEN(
			strings.NewReader(engine.InitialBoardFEN),
		)
	}
}

func BenchmarkWriteFEN(b *testing.B) {
	board := engine.NewBoard()
	for i := 0; i < b.N; i++ {
		board.WriteFEN(ioutil.Discard)
	}
}
