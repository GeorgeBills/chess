package engine_test

import (
	"github.com/GeorgeBills/chess/m/v2/engine"
	. "github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToAlgebraicNotation(t *testing.T) {
	tests := []struct {
		i        uint8
		expected string
	}{
		{A1, "a1"},
		{B1, "b1"},
		{A2, "a2"},
		{H8, "h8"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			an := engine.ToAlgebraicNotation(tt.i)
			assert.Equal(t, tt.expected, an)
		})
	}
}
