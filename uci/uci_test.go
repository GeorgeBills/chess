package uci_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/GeorgeBills/chess/m/v2/uci/mocks"
	"github.com/stretchr/testify/assert"
)

const Name = "test-engine"
const Author = "George Bills"

func TestQuitBeforeUCI(t *testing.T) {
	const in = "quit"
	r := strings.NewReader(in)
	h := &mocks.Handler{
		QuitFunc: func() {},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	assert.Equal(t, "", w.String())
}

func TestUCI(t *testing.T) {
	piper, pipew := io.Pipe()

	var calledNewGame bool
	var calledIsReady bool
	var calledSetStartingPosition bool
	var calledGoDepthWith uint8
	var calledGoTimeWith uci.TimeControl
	var calledQuit bool
	h := &mocks.Handler{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return Name, Author, nil
		},
		NewGameFunc: func() {
			calledNewGame = true
		},
		IsReadyFunc: func() {
			calledIsReady = true
		},
		SetStartingPositionFunc: func() {
			calledSetStartingPosition = true
		},
		GoDepthFunc: func(depth uint8) string {
			calledGoDepthWith = depth
			return "a1h8"
		},
		GoTimeFunc: func(tc uci.TimeControl) string {
			calledGoTimeWith = tc
			return "a8h1"
		},
		QuitFunc: func() {
			calledQuit = true
		},
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
			assert.True(t, calledIsReady)
		})

		t.Run("ucinewgame", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("ucinewgame\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.True(t, calledNewGame)
		})

		t.Run("position startpos", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("position startpos\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.True(t, calledSetStartingPosition)
		})

		// TODO: test "position fen"
		// TODO: test "position startpos moves e2e4"

		t.Run("go depth", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go depth 123\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a1h8\n", buf.String())
			assert.EqualValues(t, 123, calledGoDepthWith)
		})

		t.Run("go wtime btime winc binc", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("go wtime 300000 btime 300000 winc 0 binc 0\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "bestmove a8h1\n", buf.String())
			expected := uci.TimeControl{WhiteTime: 5 * time.Minute, BlackTime: 5 * time.Minute}
			assert.EqualValues(t, expected, calledGoTimeWith)
		})

		t.Run("quit", func(t *testing.T) {
			buf.Reset()
			pipew.Write([]byte("quit\n"))
			time.Sleep(1 * time.Millisecond) // GROSS... need to be sure parser has done the work
			assert.Equal(t, "", buf.String())
			assert.True(t, calledQuit)
		})

		pipew.Close()
	}()

	p := uci.NewParser(h, piper, buf, os.Stderr)
	p.Run()
}

func TestExtraInformation(t *testing.T) {
	const in = "uci\nquit"
	r := strings.NewReader(in)
	h := &mocks.Handler{
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
