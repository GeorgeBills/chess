package lichess

import "fmt"

//go:generate moq -out mocks/eventhandler.go -pkg mocks . EventHandler

type EventHandler interface {
	Challenge(e *EventChallenge)
	GameStart(e *EventGameStart)
}

func HandleEvents(h EventHandler, eventch <-chan interface{}) {
	for event := range eventch {
		switch v := event.(type) {
		case *EventChallenge:
			h.Challenge(v)
		case *EventGameStart:
			h.GameStart(v)
		default:
			panic(fmt.Errorf("unrecognized event type: %T", v))
		}
	}
}
