package engine_test

import (
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"

	"github.com/stretchr/testify/assert"
)

func TestFEN(t *testing.T) {
	fen := engine.NewBoard().FEN()
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	assert.Equal(t, expected, fen)
}
