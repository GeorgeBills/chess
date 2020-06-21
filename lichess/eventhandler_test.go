package lichess_test

import (
	"testing"
	"time"

	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/lichess/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleEvents(t *testing.T) {
	h := &mocks.EventHandlerMock{
		ChallengeFunc: func(e *lichess.EventChallenge) {},
		GameStartFunc: func(e *lichess.EventGameStart) {},
	}
	eventch := make(chan interface{}, 10)
	go lichess.HandleEvents(h, eventch)

	t.Run("challenge", func(t *testing.T) {
		e := &lichess.EventChallenge{}
		eventch <- e
		time.Sleep(10 * time.Millisecond)
		calls := h.ChallengeCalls()
		if assert.Len(t, calls, 1) {
			assert.Same(t, e, calls[0].E)
		}
	})

	t.Run("game start", func(t *testing.T) {
		e := &lichess.EventGameStart{}
		eventch <- e
		time.Sleep(10 * time.Millisecond)
		calls := h.GameStartCalls()
		if assert.Len(t, calls, 1) {
			assert.Same(t, e, calls[0].E)
		}
	})

	close(eventch)
}
