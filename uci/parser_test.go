package uci_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	processing = 10 * time.Millisecond
	timeout    = 1 * time.Second
)

func TestParse(t *testing.T) {
	piper, pipew := io.Pipe()
	p, commandch, stopch := uci.NewParser(piper, os.Stdout)
	require.NotNil(t, p)
	require.NotNil(t, commandch)
	require.NotNil(t, stopch)

	go p.Parse()

	// The command channel is unbuffered; the parser will error out of writing
	// to it if a goroutine isn't currently blocking on reading it. That makes
	// writing tests tricky: we'd need to fiddle with goroutines and timeouts to
	// make sure that our select is currently blocking before we write to the
	// pipe. So here we busy read off the channel onto a buffered one.
	buffered := make(chan uci.Execer, 10)
	go func() {
		for cmd := range commandch {
			t.Logf("took command off commandch: %#v", cmd)
			buffered <- cmd
		}
		close(buffered)
	}()

	t.Run("uci", func(t *testing.T) {
		pipew.Write([]byte("uci\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			assert.Equal(t, uci.CommandUCI{}, cmd)
		}
	})

	t.Run("isready", func(t *testing.T) {
		pipew.Write([]byte("isready\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			assert.Equal(t, uci.CommandIsReady{}, cmd)
		}
	})

	t.Run("ucinewgame", func(t *testing.T) {
		pipew.Write([]byte("ucinewgame\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			assert.Equal(t, uci.CommandNewGame{}, cmd)
		}
	})

	// t.Run("position startpos", func(t *testing.T) { // TODO: test "moves" here as well
	// 	pipew.Write([]byte("position startpos\n"))
	// 	time.Sleep(processing)
	// 	select {
	// 	case cmd := <-buffered:
	// 		assert.Equal(t, uci.CommandSetStartingPosition{}, cmd)
	// 	case <-time.After(timeout):
	// 		t.Errorf("timeout")
	// 	}
	// })

	t.Run("position fen", func(t *testing.T) {
		pipew.Write([]byte("position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			expected := uci.CommandSetPositionFEN{
				FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			}
			assert.Equal(t, expected, cmd)
		}
	})

	t.Run("go depth", func(t *testing.T) {
		pipew.Write([]byte("go depth 123\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			assert.Equal(t, uci.CommandGoDepth{Plies: 123}, cmd)
		}
	})

	t.Run("go nodes", func(t *testing.T) {
		pipew.Write([]byte("go nodes 456\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			assert.Equal(t, uci.CommandGoNodes{Nodes: 456}, cmd)
		}
	})

	t.Run("go wtime btime winc binc", func(t *testing.T) {
		pipew.Write([]byte("go wtime 300000 btime 300000 winc 0 binc 0\n"))
		time.Sleep(processing)
		if cmd := timeoutReadCommand(t, buffered); cmd != nil {
			expected := uci.CommandGoTime{
				uci.TimeControl{
					WhiteTime: 5 * time.Minute,
					BlackTime: 5 * time.Minute,
				},
			}
			assert.Equal(t, expected, cmd)
		}
	})

	t.Run("quit", func(t *testing.T) {
		pipew.Write([]byte("quit\n"))
		time.Sleep(processing)
		select {
		case cmd, open := <-buffered:
			// should be no more commands, and the channel should be closed
			assert.Equal(t, nil, cmd)
			assert.False(t, open)
		case <-time.After(timeout):
			t.Errorf("timeout")
		}
	})
}

func TestQuitBeforeUCI(t *testing.T) {
	piper, pipew := io.Pipe()
	p, commandch, stopch := uci.NewParser(piper, os.Stdout)
	require.NotNil(t, p)
	require.NotNil(t, commandch)
	require.NotNil(t, stopch)

	go p.Parse()

	pipew.Write([]byte("quit\n"))
	time.Sleep(processing)

	select {
	case cmd, open := <-commandch:
		// should be no commands, and the channel should be closed
		assert.Equal(t, nil, cmd)
		assert.False(t, open)
	case <-time.After(timeout):
		t.Errorf("timeout")
	}
}

func timeoutReadCommand(t *testing.T, commandch <-chan uci.Execer) uci.Execer {
	select {
	case cmd := <-commandch:
		return cmd
	case <-time.After(timeout):
		t.Errorf("timeout trying to read command off commandch")
		return nil
	}
}
