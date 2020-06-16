package internal

import (
	"fmt"
	"log"

	"github.com/GeorgeBills/chess/lichess"
)

func NewEventHandler(client Lichesser, logger *log.Logger) *EventHandler {
	return &EventHandler{
		client: client,
		logger: logger,
	}
}

type EventHandler struct {
	client Lichesser
	logger *log.Logger
}

func (h *EventHandler) Challenge(v *lichess.EventChallenge) {
	h.logger.Printf(
		"challenge; id: %s; challenger: %s; variant: %s; rated: %t",
		v.Challenge.ID,
		v.Challenge.Challenger.ID,
		v.Challenge.Variant.Key,
		v.Challenge.Rated,
	)

	if v.Challenge.Challenger.ID == "georgebills" &&
		v.Challenge.Rated == false && // require unrated for now to avoid changing my own rating
		v.Challenge.Variant.Key == lichess.VariantKeyStandard {

		h.logger.Printf("accepting challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeAccept(v.Challenge.ID)
		if err != nil {
			h.logger.Fatal(err)
		}
		// we now expect an incoming "game start" event
	} else {

		h.logger.Printf("declining challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeDecline(v.Challenge.ID)
		if err != nil {
			h.logger.Fatal(err)
		}
	}
}

func (h *EventHandler) GameStart(v *lichess.EventGameStart) {
	eventch := make(chan interface{}, 100)
	gh := NewGameHandler(v.Game.ID, h.client, h.logger)
	go StreamGameEvents(v.Game.ID, h.client, eventch, h.logger)
	go lichess.HandleGameEvents(gh, h.logger, eventch)
}

func StreamGameEvents(gameID string, client Lichesser, eventch chan<- interface{}, logger *log.Logger) {
	logger.Printf("streaming game: %s", gameID)
	err := client.BotStreamGame(gameID, eventch)
	if err != nil {
		logger.Fatal(fmt.Errorf("error streaming game events: %w", err))
	}
}

func StreamEvents(client Lichesser, eventch chan<- interface{}, logger *log.Logger) {
	logger.Printf("streaming events")
	err := client.BotStreamEvents(eventch)
	if err != nil {
		logger.Fatal(fmt.Errorf("error streaming events: %w", err))
	}
}
