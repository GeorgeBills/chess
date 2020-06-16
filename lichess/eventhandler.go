package lichess

import "log"

//go:generate moq -out mocks/eventhandler.go -pkg mocks . EventHandler

type EventHandler interface {
	Challenge(e *EventChallenge)
	GameStart(e *EventGameStart)
}

func HandleEvents(h EventHandler, logger *log.Logger, eventch <-chan interface{}) {
	for event := range eventch {
		logger.Printf("event: %#v", event)
		switch v := event.(type) {
		case *EventChallenge:
			h.Challenge(v)
		case *EventGameStart:
			h.GameStart(v)
		default:
			logger.Printf("ignoring unrecognized event type: %T", v) // errch?
		}
	}
}
