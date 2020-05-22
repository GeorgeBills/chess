package uci

import "github.com/GeorgeBills/chess/m/v2/engine"

//go:generate moq -out mocks/adapter.go -pkg mocks . Adapter

// Adapter handles events generated from parsing UCI.
type Adapter interface {
	Identify() (name, author string, other map[string]string)
	IsReady()
	NewGame()
	SetStartingPosition()
	SetPositionFEN(fen string)
	ApplyMove(move engine.FromToPromote)
	GoDepth(plies uint8) string
	GoNodes(nodes uint64) string
	GoInfinite()
	GoTime(tc TimeControl) string
	Quit()
	// TODO: most adapter methods should return error
	// TODO: return proper type instead of string for moves?
	//       adapter shouldn't need to understand UCI format
}
