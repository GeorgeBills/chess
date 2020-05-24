package uci

import chess "github.com/GeorgeBills/chess/m/v2"

//go:generate moq -out mocks/adapter.go -pkg mocks . Adapter

// Adapter handles events generated from parsing UCI.
type Adapter interface {
	Identify() (name, author string, other map[string]string)
	NewGame()
	SetStartingPosition()
	SetPositionFEN(fen string)
	ApplyMove(move chess.FromToPromoter)
	GoDepth(plies uint8) chess.FromToPromoter
	GoNodes(nodes uint64) chess.FromToPromoter
	GoInfinite(stopch <-chan struct{}, infoch chan<- Responser) chess.FromToPromoter
	GoTime(tc TimeControl) chess.FromToPromoter
	// TODO: most adapter methods should return error
}
