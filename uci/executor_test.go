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

func mustParseMove(ucin string) uci.Move {
	m, err := uci.ParseUCIN(ucin)
	if err != nil {
		panic(err)
	}
	return m
}

func TestExecuteCommands(t *testing.T) {
	a := &mocks.AdapterMock{
		IdentifyFunc: func() (name, author string, other map[string]string) {
			return Name, Author, nil
		},
		NewGameFunc:             func() error { return nil },
		SetStartingPositionFunc: func([]chess.FromToPromoter) error { return nil },
		SetPositionFENFunc:      func(string, []chess.FromToPromoter) error { return nil },
	}

	commandch := make(chan uci.Execer)
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
		ok := timeoutReadResponse(t, responsech)
		assert.Equal(t, uci.ResponseOK{}, ok)
	})

	t.Run("ucinewgame", func(t *testing.T) {
		commandch <- uci.CommandNewGame{}
		time.Sleep(processing)
		assert.Len(t, a.NewGameCalls(), 1)
	})

	t.Run("position startpos", func(t *testing.T) {
		commandch <- uci.CommandSetStartingPosition{}
		time.Sleep(processing)
		assert.Len(t, a.SetStartingPositionCalls(), 1)
	})

	close(stopch)
	close(commandch)
}

func timeoutReadResponse(t *testing.T, responsech <-chan uci.Responser) uci.Responser {
	select {
	case response := <-responsech:
		return response
	case <-time.After(timeout):
		t.Errorf("timeout trying to read response off responsech")
		return nil
	}
}
