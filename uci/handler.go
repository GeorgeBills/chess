package uci

import "github.com/GeorgeBills/chess/m/v2/engine"

//go:generate moq -out mocks/handler.go -pkg mocks . Handler

// Handler handles events generated from parsing UCI.
type Handler interface {
	Identify() (name, author string, other map[string]string)
	IsReady()
	NewGame()
	SetStartingPosition()
	SetPositionFEN(fen string)
	PlayMove(move engine.FromToPromote)
	GoDepth(plies uint8) string
	GoNodes(nodes uint64) string
	GoInfinite()
	GoTime(tc TimeControl) string
	Quit()
	// TODO: most handler methods should return error
	// TODO: is "adapter" a better name for this?
	// TODO: return proper type instead of string for moves?
	//       handler shouldn't need to understand UCI format
}
