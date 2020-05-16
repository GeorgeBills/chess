package mocks

import (
	"errors"
)

type Handler struct {
	IdentifyFunc            func() (name, author string, other map[string]string)
	IsReadyFunc             func()
	NewGameFunc             func()
	SetStartingPositionFunc func()
	SetPositionFunc         func(fen string)
	GoDepthFunc             func(plies uint8) string
	GoNodesFunc             func(nodes uint64) string
	GoInfiniteFunc          func()
}

func (h *Handler) Identify() (name, author string, other map[string]string) {
	if h.IdentifyFunc == nil {
		panic(errors.New("Identify not implemented"))
	}
	return h.IdentifyFunc()
}

func (h *Handler) IsReady() {
	if h.IsReadyFunc == nil {
		panic(errors.New("IsReady not implemented"))
	}
	h.IsReadyFunc()
}

func (h *Handler) NewGame() {
	if h.NewGameFunc == nil {
		panic(errors.New("NewGame not implemented"))
	}
	h.NewGameFunc()
}

func (h *Handler) SetStartingPosition() {
	if h.SetStartingPositionFunc == nil {
		panic(errors.New("SetStartingPosition not implemented"))
	}
	h.SetStartingPositionFunc()
}

func (h *Handler) SetPosition(fen string) {
	if h.SetPositionFunc == nil {
		panic(errors.New("SetPosition not implemented"))
	}
	h.SetPositionFunc(fen)
}

func (h *Handler) GoDepth(plies uint8) string {
	if h.GoDepthFunc == nil {
		panic(errors.New("GoDepth not implemented"))
	}
	return h.GoDepthFunc(plies)
}

func (h *Handler) GoNodes(nodes uint64) string {
	if h.GoNodesFunc == nil {
		panic(errors.New("GoNodes not implemented"))
	}
	return h.GoNodesFunc(nodes)
}

func (h *Handler) GoInfinite() {
	if h.GoInfiniteFunc == nil {
		panic(errors.New("GoInfinite not implemented"))
	}
	h.GoInfiniteFunc()
}
