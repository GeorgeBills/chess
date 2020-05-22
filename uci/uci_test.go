package uci_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	chess "github.com/GeorgeBills/chess/m/v2"
	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/GeorgeBills/chess/m/v2/uci/mocks"
	"github.com/stretchr/testify/assert"
)

const Name = "test-engine"
const Author = "George Bills"

func TestQuitBeforeUCI(t *testing.T) {
	const in = "quit"
	r := strings.NewReader(in)
	a := &mocks.AdapterMock{
		QuitFunc: func() {},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(a, r, w, ioutil.Discard)
	p.Run()
	assert.Equal(t, "", w.String())
}

func TestUCI(t *testing.T) {
	piper, pipew := io.Pipe()

	a := &mocks.AdapterMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return Name, Author, nil
		},
		NewGameFunc:             func() {},
		IsReadyFunc:             func() {},
		SetStartingPositionFunc: func() {},
		SetPositionFENFunc:      func(fen string) {},
		ApplyMoveFunc:           func(ft chess.FromToPromoter) {},
		GoDepthFunc:             func(depth uint8) string { return "a1h8" },
		GoTimeFunc:              func(tc uci.TimeControl) string { return "a8h1" },
		GoNodesFunc:             func(nodes uint64) string { return "a1h1" },
		QuitFunc:                func() {},
	}

	buf := &bytes.Buffer{}

	go func() {
		t.Run("uci", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("uci\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "id name test-engine\nid author George Bills\nuciok\n", buf.String())
		})

		t.Run("isready", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("isready\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "readyok\n", buf.String())
			assert.Len(t, a.IsReadyCalls(), 1)
		})

		t.Run("ucinewgame", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("ucinewgame\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, a.NewGameCalls(), 1)
		})

		t.Run("position startpos moves", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("position startpos moves e2e4 e7e5 f1c4\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, a.SetStartingPositionCalls(), 1)
			calls := a.ApplyMoveCalls()
			if assert.Len(t, calls, 3) {
				assert.Equal(t, "e2e4", uci.ToUCIN(calls[0].Move))
				assert.Equal(t, "e7e5", uci.ToUCIN(calls[1].Move))
				assert.Equal(t, "f1c4", uci.ToUCIN(calls[2].Move))
			}
		})

		t.Run("position fen", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			const expected = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			if assert.Len(t, a.SetPositionFENCalls(), 1) {
				assert.Equal(t, expected, a.SetPositionFENCalls()[0].Fen)
			}
		})

		t.Run("go depth", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go depth 123\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a1h8\n", buf.String())
			if assert.Len(t, a.GoDepthCalls(), 1) {
				assert.EqualValues(t, 123, a.GoDepthCalls()[0].Plies)
			}
		})

		t.Run("go nodes", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go nodes 456\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a1h1\n", buf.String())
			if assert.Len(t, a.GoNodesCalls(), 1) {
				assert.EqualValues(t, 456, a.GoNodesCalls()[0].Nodes)
			}
		})

		t.Run("go wtime btime winc binc", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go wtime 300000 btime 300000 winc 0 binc 0\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a8h1\n", buf.String())
			expected := uci.TimeControl{WhiteTime: 5 * time.Minute, BlackTime: 5 * time.Minute}
			if assert.Len(t, a.GoTimeCalls(), 1) {
				assert.Equal(t, expected, a.GoTimeCalls()[0].Tc)
			}
		})

		t.Run("quit", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("quit\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, a.QuitCalls(), 1)
		})

		pipew.Close()
	}()

	p := uci.NewParser(a, piper, buf, os.Stderr)
	p.Run()
}

func TestExtraInformation(t *testing.T) {
	const in = "uci\nquit"
	r := strings.NewReader(in)
	a := &mocks.AdapterMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return "super-chess", "Jane Smith", map[string]string{
				"version":      "v1.2.3",
				"release-date": "2020-05-16",
			}
		},
		QuitFunc: func() {},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(a, r, w, ioutil.Discard)
	p.Run()
	const expected = "id name super-chess\nid author Jane Smith\nid release-date 2020-05-16\nid version v1.2.3\nuciok\n"
	assert.Equal(t, expected, w.String())
}
