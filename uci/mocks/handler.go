package mocks

import (
	"errors"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

// Handler implements the uci.Handler interface in an easily mockable fashion.
type Handler struct {
	IdentifyFunc            func() (name, author string, other map[string]string)
	IsReadyFunc             func()
	NewGameFunc             func()
	SetStartingPositionFunc func()
	SetPositionFENFunc      func(fen string)
	PlayMoveFunc            func(move engine.FromTo)
	GoDepthFunc             func(plies uint8) string
	GoNodesFunc             func(nodes uint64) string
	GoInfiniteFunc          func()
	QuitFunc                func()
}

// Identify implements uci.Handler.Identify()
// It does so by calling the IdentifyFunc explicitly added to the handler.
func (h *Handler) Identify() (name, author string, other map[string]string) {
	if h.IdentifyFunc == nil {
		panic(errors.New("Identify not implemented"))
	}
	return h.IdentifyFunc()
}

// IsReady implements uci.Handler.IsReady().
// It does so by calling the IsReadyFunc explicitly added to the handler.
func (h *Handler) IsReady() {
	if h.IsReadyFunc == nil {
		panic(errors.New("IsReady not implemented"))
	}
	h.IsReadyFunc()
}

// NewGame implements uci.Handler.NewGame().
// It does so by calling the NewGameFunc explicitly added to the handler.
func (h *Handler) NewGame() {
	if h.NewGameFunc == nil {
		panic(errors.New("NewGame not implemented"))
	}
	h.NewGameFunc()
}

// SetStartingPosition implements uci.Handler.SetStartingPosition().
// It does so by calling the SetStartingPositionFunc explicitly added to the handler.
func (h *Handler) SetStartingPosition() {
	if h.SetStartingPositionFunc == nil {
		panic(errors.New("SetStartingPosition not implemented"))
	}
	h.SetStartingPositionFunc()
}

// SetPositionFEN implements uci.Handler.SetPosition().
// It does so by calling the SetPositionFunc explicitly added to the handler.
func (h *Handler) SetPositionFEN(fen string) {
	if h.SetPositionFENFunc == nil {
		panic(errors.New("SetPosition not implemented"))
	}
	h.SetPositionFENFunc(fen)
}

// PlayMove implements uci.Handler.PlayMove().
// It does so by calling the PlayMoveFunc explicitly added to the handler.
func (h *Handler) PlayMove(move engine.FromTo) {
	if h.PlayMoveFunc == nil {
		panic(errors.New("PlayMove not implemented"))
	}
	h.PlayMoveFunc(move)
}

// GoDepth implements uci.Handler.GoDepth().
// It does so by calling the GoDepthFunc explicitly added to the handler.
func (h *Handler) GoDepth(plies uint8) string {
	if h.GoDepthFunc == nil {
		panic(errors.New("GoDepth not implemented"))
	}
	return h.GoDepthFunc(plies)
}

// GoNodes implements uci.Handler.GoNodes().
// It does so by calling the GoNodesFunc explicitly added to the handler.
func (h *Handler) GoNodes(nodes uint64) string {
	if h.GoNodesFunc == nil {
		panic(errors.New("GoNodes not implemented"))
	}
	return h.GoNodesFunc(nodes)
}

// GoInfinite implements uci.Handler.GoInfinite().
// It does so by calling the GoInfiniteFunc explicitly added to the handler.
func (h *Handler) GoInfinite() {
	if h.GoInfiniteFunc == nil {
		panic(errors.New("GoInfinite not implemented"))
	}
	h.GoInfiniteFunc()
}

// Quit implements uci.Handler.Quit().
// It does so by calling the QuitFunc explicitly added to the handler.
func (h *Handler) Quit() {
	if h.QuitFunc == nil {
		panic(errors.New("Quit not implemented"))
	}
	h.QuitFunc()
}
