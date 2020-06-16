package main

import (
	"fmt"
	"strings"

	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal"
	"github.com/GeorgeBills/chess/engine"
)

func NewGameFactory() gameFactory {
	return gameFactory{}
}

type gameFactory struct{}

func (f gameFactory) NewGame() internal.Game {
	b := engine.NewBoard()
	g := engine.NewGame(b)
	return NewGameWrapper(g)
}

func (f gameFactory) NewGameFromFEN(fen string) (internal.Game, error) {
	b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
	if err != nil {
		return nil, fmt.Errorf("error while parsing FEN: %w", err)
	}
	g := engine.NewGame(b)
	return NewGameWrapper(g), nil
}
