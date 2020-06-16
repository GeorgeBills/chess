package internal_test

import (
	"log"
	"os"
	"testing"

	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal"
	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal/mocks"
	"github.com/GeorgeBills/chess/lichess"
	"github.com/stretchr/testify/assert"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func TestEventHandlerChallengeAccept(t *testing.T) {
	mockClient := &mocks.LichesserMock{
		ChallengeAcceptFunc: func(challengeID string) error { return nil },
	}
	mockFactory := &mocks.GameFactoryMock{}
	h := internal.NewEventHandler(mockClient, logger, mockFactory)
	h.Challenge(
		&lichess.EventChallenge{
			Challenge: lichess.EventChallengeChallenge{
				ID:         "CHWmd4",
				Rated:      false,
				Challenger: lichess.Player{ID: "georgebills"},
				Variant:    lichess.Variant{Key: "standard"},
			},
		},
	)

	calls := mockClient.ChallengeAcceptCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "CHWmd4", calls[0].ChallengeID)
	}
}

func TestEventHandlerChallengeDecline(t *testing.T) {
	tests := []struct {
		name        string
		challenge   *lichess.EventChallenge
		challengeID string
	}{
		{
			"rated",
			&lichess.EventChallenge{
				lichess.EventChallengeChallenge{
					ID:         "Jp7EUq",
					Rated:      true,
					Challenger: lichess.Player{ID: "georgebills"},
					Variant:    lichess.Variant{Key: "standard"},
				},
			},
			"Jp7EUq",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.LichesserMock{
				ChallengeDeclineFunc: func(challengeID string) error { return nil },
			}
			mockFactory := &mocks.GameFactoryMock{}
			h := internal.NewEventHandler(mockClient, logger, mockFactory)

			h.Challenge(tt.challenge)

			calls := mockClient.ChallengeDeclineCalls()
			if assert.Len(t, calls, 1) {
				assert.Equal(t, "Jp7EUq", calls[0].ChallengeID)
			}
		})
	}
}
