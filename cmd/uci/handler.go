package main

import (
	"fmt"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

type handler struct {
	game *engine.Game
}

func (h *handler) Identify() {
	fmt.Println(etgID, etgIDName, name)
	fmt.Println(etgID, etgIDAuthor, author)
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

func (h *handler) GoDepth(plies uint8) {
	m, _ := h.game.BestMoveToDepth(plies * 2)
	fmt.Println(etgBestMove, m.SAN())
}

func (h *handler) GoNodes(nodes uint64) {
	panic("GoNodes not implemented")
}

func (h *handler) GoInfinite() {
	panic("GoInfinite not implemented")
}
