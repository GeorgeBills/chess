package lichess

import "log"

//go:generate moq -out mocks/gamehandler.go -pkg mocks . GameHandler

type GameHandler interface {
	GameFull(e *EventGameFull)
	GameState(e *EventGameState)
	ChatLine(e *EventChatLine)
}

func HandleGameEvents(h GameHandler, logger *log.Logger, eventch <-chan interface{}) {
	for event := range eventch {
		logger.Printf("game event: %#v", event) // TODO: include game id
		switch v := event.(type) {
		case *EventGameFull:
			h.GameFull(v)
		case *EventGameState:
			h.GameState(v)
		case *EventChatLine:
			h.ChatLine(v)
		default:
			logger.Printf("ignoring unrecognized game event type: %T", v) // errch
		}
	}
}
