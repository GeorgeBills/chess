package main

import "github.com/GeorgeBills/chess/lichess"

type eventHandler struct {
	client *lichess.Client
}

func (h *eventHandler) Challenge(v *lichess.EventChallenge) {
	logger.Printf(
		"challenge; id: %s; challenger: %s; variant: %s; rated: %t",
		v.Challenge.ID,
		v.Challenge.Challenger.ID,
		v.Challenge.Variant.Key,
		v.Challenge.Rated,
	)

	if v.Challenge.Challenger.ID == "georgebills" &&
		v.Challenge.Rated == false && // require unrated for now to avoid changing my own rating
		v.Challenge.Variant.Key == lichess.VariantKeyStandard {

		logger.Printf("accepting challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeAccept(v.Challenge.ID)
		if err != nil {
			logger.Fatal(err)
		}
		// we now expect an incoming "game start" event
	} else {

		logger.Printf("declining challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeDecline(v.Challenge.ID)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func (h *eventHandler) GameStart(v *lichess.EventGameStart) {
	eventch := make(chan interface{}, 100)
	gh := &gameHandler{
		gameID: v.Game.ID,
		client: h.client,
	}
	go streamGameEvents(v.Game.ID, h.client, eventch)
	go lichess.HandleGameEvents(gh, logger, eventch)
}
