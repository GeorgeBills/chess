package uci

import chess "github.com/GeorgeBills/chess/m/v2"

//go:generate moq -out mocks/adapter.go -pkg mocks . Adapter

// Adapter handles events generated from parsing UCI.
type Adapter interface {
	Identify() (name, author string, other map[string]string)
	NewGame() error
	SetStartingPosition() error
	SetPositionFEN(fen string) error
	ApplyMove(move chess.FromToPromoter) error
	GoDepth(plies uint8) (chess.FromToPromoter, error)
	GoNodes(nodes uint64) (chess.FromToPromoter, error)
	GoInfinite(stopch <-chan struct{}, infoch chan<- Responser) (chess.FromToPromoter, error)
	GoTime(tc TimeControl) (chess.FromToPromoter, error)
}
