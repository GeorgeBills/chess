package main

import (
	"io"
	"log"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

// Name is the name of our engine.
const Name = "github.com/GeorgeBills/chess"

// Author is the author of our engine.
const Author = "George Bills"

// NewHandler returns a new handler.
func NewHandler(logw io.Writer) *handler {
	return &handler{
		logger: log.New(logw, "handler: ", log.LstdFlags),
	}
}

type handler struct {
	logger *log.Logger
	game   *engine.Game
}

func (h *handler) Identify() (name, author string, other map[string]string) {
	h.logger.Println("identify")
	return Name, Author, nil
}

func (h *handler) IsReady() {
	h.logger.Println("is ready")
	// TODO: block on mutex in engine if we're waiting on anything slow?
}

func (h *handler) NewGame() {
	h.logger.Println("initialised new game")
	g := engine.NewGame(nil) // TODO: return pointer
	h.game = &g
}

func (h *handler) SetStartingPosition() {
	h.logger.Println("set starting position")
	b := engine.NewBoard()
	h.game.SetBoard(&b)
}

func (h *handler) SetPosition(fen string) {
	h.logger.Println("set position")
	panic("SetPosition not implemented")
}

func (h *handler) GoDepth(plies uint8) string {
	h.logger.Println("go depth")
	m, _ := h.game.BestMoveToDepth(plies * 2)
	return m.SAN()
}

func (h *handler) GoNodes(nodes uint64) string {
	h.logger.Println("go nodes")
	panic("GoNodes not implemented")
}

func (h *handler) GoInfinite() {
	h.logger.Println("go infinite")
	panic("GoInfinite not implemented")
}
