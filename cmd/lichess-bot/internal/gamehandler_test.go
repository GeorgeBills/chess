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

func TestGameFullToMove(t *testing.T) {
	m := &mocks.LichesserMock{
		BotMakeMoveFunc: func(gameID string, move chess.FromToPromoter, offeringDraw bool) error {
			return nil
		},
	}
	h := internal.NewGameHandler("QSdKvR", m, logger)
	h.GameFull(&lichess.EventGameFull{
		ID:         "QSdKvR",
		White:      lichess.Player{ID: "gbcb"},
		Black:      lichess.Player{ID: "georgebills"},
		InitialFen: "startpos",
		State:      lichess.EventGameState{Moves: ""},
	})

	calls := m.BotMakeMoveCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "QSdKvR", calls[0].GameID)
		assert.Equal(t, "e2e4", uci.ToUCIN(calls[0].Move))
		assert.Equal(t, false, calls[0].OfferingDraw)
	}
}
