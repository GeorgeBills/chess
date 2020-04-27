package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPiece_RunePanic(t *testing.T) {
	p := engine.Piece(0b11111111) // invalid piece
	assert.Panics(t, func() { _ = p.Rune() })
}
