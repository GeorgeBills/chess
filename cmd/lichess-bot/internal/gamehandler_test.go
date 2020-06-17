package internal_test

import (
	"testing"

	"github.com/GeorgeBills/chess"
	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal"
	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal/mocks"
	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/uci"
	"github.com/stretchr/testify/assert"
)

type move struct {
	from, to  uint8
	promoteTo chess.PromoteTo
}

func (m move) From() uint8                { return m.from }
func (m move) To() uint8                  { return m.to }
func (m move) PromoteTo() chess.PromoteTo { return m.promoteTo }

func TestGameFullToMove(t *testing.T) {
	mockClient := &mocks.LichesserMock{
		BotMakeMoveFunc: func(gameID string, move chess.FromToPromoter, offeringDraw bool) error {
			return nil
		},
	}
	mockFactory := &mocks.GameFactoryMock{
		NewGameFunc: func() internal.Game {
			return &mocks.GameMock{
				BestMoveFunc: func() (chess.FromToPromoter, int16) {
					return move{from: 12, to: 28}, 10 // e2e4
					// TODO: copy square constants up to chess pkg
				},
			}
		},
	}
	h := internal.NewGameHandler("QSdKvR", mockClient, logger, mockFactory)
	h.GameFull(&lichess.EventGameFull{
		ID:         "QSdKvR",
		White:      lichess.Player{ID: "gbcb"},
		Black:      lichess.Player{ID: "georgebills"},
		InitialFen: "startpos",
		State:      lichess.EventGameState{Moves: ""},
	})

	calls := mockClient.BotMakeMoveCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "QSdKvR", calls[0].GameID)
		assert.Equal(t, "e2e4", uci.ToUCIN(calls[0].Move))
		assert.Equal(t, false, calls[0].OfferingDraw)
	}
}

func TestStreamGameEvents(t *testing.T) {
	mockClient := &mocks.LichesserMock{
		BotStreamGameFunc: func(gameID string, eventch chan<- interface{}) error {
			eventch <- &lichess.EventGameFull{}
			eventch <- &lichess.EventGameState{}
			eventch <- &lichess.EventChatLine{}
			close(eventch)
			return nil
		},
	}
	eventch := make(chan interface{}, 100)
	internal.StreamGameEvents("NdHWLn", mockClient, eventch, logger)

	received := []interface{}{}
	for event := range eventch {
		received = append(received, event)
	}

	assert.Equal(t, &lichess.EventGameFull{}, received[0])
	assert.Equal(t, &lichess.EventGameState{}, received[1])
	assert.Equal(t, &lichess.EventChatLine{}, received[2])
}
