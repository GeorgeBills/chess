package lichess_test

import (
	"testing"
	"time"

	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/lichess/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleGameEvents(t *testing.T) {
	h := &mocks.GameHandlerMock{
		ChatLineFunc:  func(e *lichess.EventChatLine) {},
		GameFullFunc:  func(e *lichess.EventGameFull) {},
		GameStateFunc: func(e *lichess.EventGameState) {},
	}
	eventch := make(chan interface{}, 10)
	go lichess.HandleGameEvents(h, eventch)

	t.Run("chat line", func(t *testing.T) {
		e := &lichess.EventChatLine{}
		eventch <- e
		time.Sleep(10 * time.Millisecond)
		calls := h.ChatLineCalls()
		if assert.Len(t, calls, 1) {
			assert.Same(t, e, calls[0].E)
		}
	})

	t.Run("game full", func(t *testing.T) {
		e := &lichess.EventGameFull{}
		eventch <- e
		time.Sleep(10 * time.Millisecond)
		calls := h.GameFullCalls()
		if assert.Len(t, calls, 1) {
			assert.Same(t, e, calls[0].E)
		}
	})

	t.Run("game state", func(t *testing.T) {
		e := &lichess.EventGameState{}
		eventch <- e
		time.Sleep(10 * time.Millisecond)
		calls := h.GameStateCalls()
		if assert.Len(t, calls, 1) {
			assert.Same(t, e, calls[0].E)
		}
	})

	close(eventch)
}
