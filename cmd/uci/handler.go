package main

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
)

// Name is the name of our engine.
const Name = "github.com/GeorgeBills/chess"

// Author is the author of our engine.
const Author = "George Bills"

type handler struct {
	game *engine.Game
}

func (h *handler) Identify() (name, author string, other map[string]string) {
	return Name, Author, nil
}

func (h *handler) IsReady() {
	// TODO: block on mutex in engine if we're waiting on anything slow?
}

func (h *handler) NewGame() {
	logger.Println("initialised new game")
	g := engine.NewGame(nil) // TODO: return pointer
	h.game = &g
}

func (h *handler) SetStartingPosition() {
	logger.Println("set starting position")
	b := engine.NewBoard()
	h.game.SetBoard(&b)
}

func (h *handler) SetPosition(fen string) {
	panic("SetPosition not implemented")
}

func (h *handler) GoDepth(plies uint8) string {
	m, _ := h.game.BestMoveToDepth(plies * 2)
	return m.SAN()
}

func (h *handler) GoNodes(nodes uint64) string {
	panic("GoNodes not implemented")
}

func (h *handler) GoInfinite() {
	panic("GoInfinite not implemented")
}
