package internal

import "github.com/GeorgeBills/chess"

//go:generate moq -out mocks/lichesser.go -pkg mocks . Lichesser

type Lichesser interface {
	ChallengeAccept(challengeID string) error
	ChallengeDecline(challengeID string) error
	BotStreamEvents(eventch chan<- interface{}) error
	BotStreamGame(gameID string, eventch chan<- interface{}) error
	BotMakeMove(gameID string, move chess.FromToPromoter, offeringDraw bool) error
}
