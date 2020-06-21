package lichess

import "fmt"

//go:generate moq -out mocks/gamehandler.go -pkg mocks . GameHandler

type GameHandler interface {
	GameFull(e *EventGameFull)
	GameState(e *EventGameState)
	ChatLine(e *EventChatLine)
}

func HandleGameEvents(h GameHandler, eventch <-chan interface{}) {
	for event := range eventch {
		switch v := event.(type) {
		case *EventGameFull:
			h.GameFull(v)
		case *EventGameState:
			h.GameState(v)
		case *EventChatLine:
			h.ChatLine(v)
		default:
			panic(fmt.Errorf("unrecognized game event type: %T", v))
		}
	}
}
