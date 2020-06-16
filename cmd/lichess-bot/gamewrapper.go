package main

import (
	"fmt"

	"github.com/GeorgeBills/chess"
	"github.com/GeorgeBills/chess/engine"
	"github.com/GeorgeBills/chess/uci"
)

type engineGameWrapper struct {
	wrapped *engine.Game
}

func NewGameWrapper(game *engine.Game) engineGameWrapper {
	return engineGameWrapper{game}
}

func (g engineGameWrapper) MakeMove(movestr string) error {
	parsed, err := uci.ParseUCIN(movestr)
	if err != nil {
		return fmt.Errorf("error applying move: %w", err)
	}

	move, err := g.wrapped.HydrateMove(parsed)
	if err != nil {
		return fmt.Errorf("error making move: %w", err)
	}

	g.wrapped.MakeMove(move) // TODO: rename make, unmake to apply, unapply?
	return nil
}

func (g engineGameWrapper) BestMove() (chess.FromToPromoter, int16) {
	stopch := make(chan struct{})
	statusch := make(chan engine.SearchStatus)
	return g.wrapped.BestMoveToDepth(4, stopch, statusch)
}
