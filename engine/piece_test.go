package engine_test

import (
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
)

func TestPiece_RunePanic(t *testing.T) {
	const p = engine.Piece(0b11111111) // invalid piece
	assert.Panics(t, func() { _ = p.Rune() })
}
