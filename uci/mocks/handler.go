package mocks

import (
	"testing"
)

type Handler struct {
	TB                      testing.TB
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
	if h.IdentifyFunc != nil {
		return h.IdentifyFunc()
	}
	h.TB.Fatal("Identify not implemented")
	return "", "", nil
}

func (h *Handler) IsReady() {
	if h.IsReadyFunc != nil {
		h.IsReadyFunc()
	}
	h.TB.Fatal("IsReady not implemented")
	return
}

func (h *Handler) NewGame() {
	if h.NewGameFunc != nil {
		h.NewGameFunc()
	}
	h.TB.Fatal("NewGame not implemented")
	return
}

func (h *Handler) SetStartingPosition() {
	if h.SetStartingPositionFunc != nil {
		h.SetStartingPositionFunc()
	}
	h.TB.Fatal("SetStartingPosition not implemented")
	return
}

func (h *Handler) SetPosition(fen string) {
	if h.SetPositionFunc != nil {
		h.SetPositionFunc(fen)
	}
	h.TB.Fatal("SetPosition not implemented")
	return
}

func (h *Handler) GoDepth(plies uint8) string {
	if h.GoDepthFunc != nil {
		return h.GoDepthFunc(plies)
	}
	h.TB.Fatal("GoDepth not implemented")
	return ""
}

func (h *Handler) GoNodes(nodes uint64) string {
	if h.GoNodesFunc != nil {
		return h.GoNodesFunc(nodes)
	}
	h.TB.Fatal("GoNodes not implemented")
	return ""
}

func (h *Handler) GoInfinite() {
	if h.GoInfiniteFunc != nil {
		h.GoInfiniteFunc()
	}
	h.TB.Fatal("GoInfinite not implemented")
	return
}
