package main

import (
	"io"
	"log"
	"strings"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/GeorgeBills/chess/m/v2/uci"
)

// Name is the name of our engine.
const Name = "github.com/GeorgeBills/chess"

// Author is the author of our engine.
const Author = "George Bills"

// newHandler returns a new handler.
func newHandler(logw io.Writer) *handler {
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
	// TODO: nil check game on SetPositionFEN
	//       or just return a new game from SetBoard if none already?
	// if h.game == nil {
	// 	return errors.New("no game")
	// }
	b := engine.NewBoard()
	h.game.SetBoard(&b)
}

func (h *handler) SetPositionFEN(fen string) {
	h.logger.Println("set position")

	// TODO: nil check game on SetPositionFEN
	//       or just return a new game from SetBoard if none already?

	b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	// TODO: return err from SetPositionFEN
	// if err != nil {
	// 	return err
	// }

	h.game.SetBoard(b)
}

func (h *handler) PlayMove(move engine.FromTo) {
	h.logger.Printf("playing move: %v", move)
	m, err := h.game.HydrateMove(move)
	if err != nil {
		panic(err) // FIXME: return errors from most handler methods...
	}
	h.game.MakeMove(m)
}

func (h *handler) GoDepth(plies uint8) string {
	h.logger.Println("go depth")
	m, _ := h.game.BestMoveToDepth(plies * 2)
	return m.UCIN()
}

func (h *handler) GoNodes(nodes uint64) string {
	h.logger.Println("go nodes")
	panic("GoNodes not implemented")
}

func (h *handler) GoInfinite() {
	h.logger.Println("go infinite")
	panic("GoInfinite not implemented")
}

func (h *handler) GoTime(tc uci.TimeControl) string {
	h.logger.Println("go time")
	m, _ := h.game.BestMoveToTime(tc.WhiteTime, tc.BlackTime, tc.WhiteIncrement, tc.BlackIncrement)
	return m.UCIN()
}

func (h *handler) Quit() { h.logger.Println("quit") } // nothing to cleanup
