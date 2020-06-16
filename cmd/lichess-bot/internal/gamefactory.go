package internal

import "github.com/GeorgeBills/chess"

//go:generate moq -out mocks/gamefactory.go -pkg mocks . GameFactory

type GameFactory interface {
	NewGame() Game
	NewGameFromFEN(fen string) (Game, error)
}

//go:generate moq -out mocks/game.go -pkg mocks . Game

type Game interface {
	MakeMove(move string) error
	BestMove() (move chess.FromToPromoter, score int16)
}
