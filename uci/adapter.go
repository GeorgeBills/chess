package uci

import chess "github.com/GeorgeBills/chess/m/v2"

//go:generate moq -out mocks/adapter.go -pkg mocks . Adapter

// Adapter handles events generated from parsing UCI.
type Adapter interface {
	Identify() (name, author string, other map[string]string)
	NewGame() error
	SetStartingPosition(moves []chess.FromToPromoter) error
	SetPositionFEN(fen string, moves []chess.FromToPromoter) error
	GoDepth(plies uint8, stopch <-chan struct{}, infoch chan<- Response) (chess.FromToPromoter, error)
	GoNodes(nodes uint64, stopch <-chan struct{}, infoch chan<- Response) (chess.FromToPromoter, error)
	GoTime(tc TimeControl, stopch <-chan struct{}, infoch chan<- Response) (chess.FromToPromoter, error)
	GoInfinite(stopch <-chan struct{}, infoch chan<- Response) (chess.FromToPromoter, error)
}
