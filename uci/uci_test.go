package uci_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/GeorgeBills/chess/m/v2/uci/mocks"
	"github.com/stretchr/testify/assert"
)

const Name = "test-engine"
const Author = "George Bills"

func TestQuitBeforeUCI(t *testing.T) {
	const in = "quit"
	r := strings.NewReader(in)
	h := &mocks.HandlerMock{
		QuitFunc: func() {},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	assert.Equal(t, "", w.String())
}

func TestUCI(t *testing.T) {
	piper, pipew := io.Pipe()

	h := &mocks.HandlerMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return Name, Author, nil
		},
		NewGameFunc:             func() {},
		IsReadyFunc:             func() {},
		SetStartingPositionFunc: func() {},
		SetPositionFENFunc:      func(fen string) {},
		PlayMoveFunc:            func(ft engine.FromToPromote) {},
		GoDepthFunc:             func(depth uint8) string { return "a1h8" },
		GoTimeFunc:              func(tc uci.TimeControl) string { return "a8h1" },
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
			assert.Len(t, h.IsReadyCalls(), 1)
		})

		t.Run("ucinewgame", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("ucinewgame\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, h.NewGameCalls(), 1)
		})

		t.Run("position startpos moves", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("position startpos moves e2e4 e7e5 f1c4\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, h.SetStartingPositionCalls(), 1)
			calls := h.PlayMoveCalls()
			if assert.Len(t, calls, 3) {
				assert.Equal(t, "e2e4", engine.UCIN(calls[0].Move))
				assert.Equal(t, "e7e5", engine.UCIN(calls[1].Move))
				assert.Equal(t, "f1c4", engine.UCIN(calls[2].Move))
			}
		})

		// TODO: test "position fen"

		t.Run("go depth", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go depth 123\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a1h8\n", buf.String())
			if assert.Len(t, h.GoDepthCalls(), 1) {
				assert.EqualValues(t, 123, h.GoDepthCalls()[0].Plies)
			}
		})

		t.Run("go wtime btime winc binc", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go wtime 300000 btime 300000 winc 0 binc 0\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a8h1\n", buf.String())
			expected := uci.TimeControl{WhiteTime: 5 * time.Minute, BlackTime: 5 * time.Minute}
			if assert.Len(t, h.GoTimeCalls(), 1) {
				assert.Equal(t, expected, h.GoTimeCalls()[0].Tc)
			}
		})

		t.Run("quit", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("quit\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.Len(t, h.QuitCalls(), 1)
		})

		pipew.Close()
	}()

	p := uci.NewParser(h, piper, buf, os.Stderr)
	p.Run()
}

func TestExtraInformation(t *testing.T) {
	const in = "uci\nquit"
	r := strings.NewReader(in)
	h := &mocks.HandlerMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return "super-chess", "Jane Smith", map[string]string{
				"version":      "v1.2.3",
				"release-date": "2020-05-16",
			}
		},
		QuitFunc: func() {},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	const expected = "id name super-chess\nid author Jane Smith\nid release-date 2020-05-16\nid version v1.2.3\nuciok\n"
	assert.Equal(t, expected, w.String())
}
