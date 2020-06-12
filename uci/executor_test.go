package uci_test

import (
	"os"
	"testing"
	"time"

	chess "github.com/GeorgeBills/chess/m/v2"
	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/GeorgeBills/chess/m/v2/uci/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	Name   = "test-engine"
	Author = "George Bills"
)

func TestExecuteCommands(t *testing.T) {
	a := &mocks.AdapterMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return Name, Author, map[string]string{
				"version":      "1.2.3",
				"release-date": "2020-05-26",
			}
		},
		NewGameFunc:             func() error { return nil },
		SetStartingPositionFunc: func([]chess.FromToPromoter) error { return nil },
		SetPositionFENFunc:      func(string, []chess.FromToPromoter) error { return nil },
		GoDepthFunc: func(plies uint8, stopch <-chan struct{}, infoch chan<- uci.Response) (chess.FromToPromoter, error) {
			return mustParseMove("a1a2"), nil
		},
		GoNodesFunc: func(nodes uint64, stopch <-chan struct{}, infoch chan<- uci.Response) (chess.FromToPromoter, error) {
			return mustParseMove("b3b4"), nil
		},
		GoTimeFunc: func(tc uci.TimeControl, stopch <-chan struct{}, infoch chan<- uci.Response) (chess.FromToPromoter, error) {
			return mustParseMove("c5c6"), nil
		},
		GoInfiniteFunc: func(stopch <-chan struct{}, infoch chan<- uci.Response) (chess.FromToPromoter, error) {
			return mustParseMove("d7d8"), nil
		},
	}

	commandch := make(chan uci.Command)
	stopch := make(chan struct{})
	e, responsech := uci.NewExecutor(commandch, stopch, a, os.Stdout)

	go e.ExecuteCommands()

	t.Run("uci", func(t *testing.T) {
		commandch <- uci.CommandUCI{}
		time.Sleep(processing)
		assert.Len(t, a.IdentifyCalls(), 1)

		name := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseID{Key: "name", Value: "test-engine"}, name)

		author := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseID{Key: "author", Value: "George Bills"}, author)

		date := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseID{Key: "release-date", Value: "2020-05-26"}, date)

		version := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseID{Key: "version", Value: "1.2.3"}, version)

		ok := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseOK{}, ok)
	})

	t.Run("ucinewgame", func(t *testing.T) {
		commandch <- uci.CommandNewGame{}
		time.Sleep(processing)
		assert.Len(t, a.NewGameCalls(), 1)
	})

	t.Run("isready", func(t *testing.T) {
		commandch <- uci.CommandIsReady{}
		time.Sleep(processing)

		response := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseIsReady{}, response)
	})

	t.Run("position startpos", func(t *testing.T) {
		cmd := &uci.CommandSetStartingPosition{}
		cmd.AppendMove(mustParseMove("f1f2"))
		cmd.AppendMove(mustParseMove("f3f4"))
		commandch <- cmd
		time.Sleep(processing)

		calls := a.SetStartingPositionCalls()
		if assert.Len(t, calls, 1) {
			expected := []chess.FromToPromoter{
				mustParseMove("f1f2"),
				mustParseMove("f3f4"),
			}
			assert.Equal(t, expected, calls[0].Moves)
		}
	})

	t.Run("position fen", func(t *testing.T) {
		cmd := &uci.CommandSetPositionFEN{
			FEN: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		}
		cmd.AppendMove(mustParseMove("e1e2"))
		cmd.AppendMove(mustParseMove("e3e4"))
		commandch <- cmd
		time.Sleep(processing)

		calls := a.SetPositionFENCalls()
		if assert.Len(t, calls, 1) {
			assert.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", calls[0].Fen)
			expected := []chess.FromToPromoter{
				mustParseMove("e1e2"),
				mustParseMove("e3e4"),
			}
			assert.Equal(t, expected, calls[0].Moves)
		}
	})

	t.Run("go depth", func(t *testing.T) {
		commandch <- uci.CommandGoDepth{Plies: 123}
		time.Sleep(processing)

		calls := a.GoDepthCalls()
		if assert.Len(t, calls, 1) {
			assert.NotNil(t, calls[0].Infoch)
			assert.NotNil(t, calls[0].Stopch)
			assert.EqualValues(t, 123, calls[0].Plies)
		}

		response := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseBestMove{mustParseMove("a1a2")}, response)
	})

	t.Run("go nodes", func(t *testing.T) {
		commandch <- uci.CommandGoNodes{Nodes: 456}
		time.Sleep(processing)

		calls := a.GoNodesCalls()
		if assert.Len(t, calls, 1) {
			assert.NotNil(t, calls[0].Infoch)
			assert.NotNil(t, calls[0].Stopch)
			assert.EqualValues(t, 456, calls[0].Nodes)
		}

		response := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseBestMove{mustParseMove("b3b4")}, response)
	})

	t.Run("go time", func(t *testing.T) {
		commandch <- uci.CommandGoTime{
			uci.TimeControl{
				WhiteTime:      1 * time.Minute,
				BlackTime:      2 * time.Minute,
				WhiteIncrement: 3 * time.Second,
				BlackIncrement: 4 * time.Second,
			},
		}
		time.Sleep(processing)

		calls := a.GoTimeCalls()
		if assert.Len(t, calls, 1) {
			assert.NotNil(t, calls[0].Infoch)
			assert.NotNil(t, calls[0].Stopch)
			assert.Equal(t, "wtime 60000 btime 120000 winc 3000 binc 4000", calls[0].Tc.String())
		}

		response := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseBestMove{mustParseMove("c5c6")}, response)
	})

	t.Run("go infinite", func(t *testing.T) {
		commandch <- uci.CommandGoInfinite{}
		time.Sleep(processing)

		calls := a.GoInfiniteCalls()
		if assert.Len(t, calls, 1) {
			assert.NotNil(t, calls[0].Infoch)
			assert.NotNil(t, calls[0].Stopch)
		}

		response := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseBestMove{mustParseMove("d7d8")}, response)
	})

	close(stopch)
	close(commandch)
}

func timeoutReadResponse(t *testing.T, responsech <-chan uci.Response) uci.Response {
	select {
	case response := <-responsech:
		return response
	case <-time.After(timeout):
		t.Errorf("timeout trying to read response off responsech")
		return nil
	}
}
