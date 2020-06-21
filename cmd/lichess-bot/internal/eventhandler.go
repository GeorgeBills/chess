package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/GeorgeBills/chess/lichess"
)

func NewEventHandler(client Lichesser, logger *log.Logger, factory GameFactory) *EventHandler {
	return &EventHandler{
		client:  client,
		logger:  logger,
		factory: factory,
	}
}

type EventHandler struct {
	client  Lichesser
	logger  *log.Logger
	factory GameFactory
}

func (h *EventHandler) Challenge(v *lichess.EventChallenge) {
	h.logger.Printf(
		"challenge %s: challenger: %s; variant: %s; %s",
		v.Challenge.ID,
		v.Challenge.Challenger.ID,
		strings.ToLower(v.Challenge.Variant.Name),
		ratedstr(v.Challenge.Rated),
	)

	if v.Challenge.Challenger.ID == "georgebills" &&
		v.Challenge.Rated == false && // require unrated for now to avoid changing my own rating
		v.Challenge.Variant.Key == lichess.VariantKeyStandard {

		h.logger.Printf("challenge %s: accepting", v.Challenge.ID)
		err := h.client.ChallengeAccept(v.Challenge.ID)
		if err != nil {
			h.logger.Fatal(err)
		}
		// we now expect an incoming "game start" event
	} else {

		h.logger.Printf("challenge %s: declining", v.Challenge.ID)
		err := h.client.ChallengeDecline(v.Challenge.ID)
		if err != nil {
			h.logger.Fatal(err)
		}
	}
}

func (h *EventHandler) GameStart(v *lichess.EventGameStart) {
	eventch := make(chan interface{}, 100)
	gh := NewGameHandler(v.Game.ID, h.client, h.logger, h.factory)
	go StreamGameEvents(v.Game.ID, h.client, eventch, h.logger)
	go lichess.HandleGameEvents(gh, eventch)
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
