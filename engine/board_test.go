package engine_test

import (
	"testing"

	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardToString(t *testing.T) {
	str := engine.NewBoard().String()
	expected := `♜♞♝♛♚♝♞♜
♟♟♟♟♟♟♟♟
□□□□□□□□
□□□□□□□□
□□□□□□□□
□□□□□□□□
♙♙♙♙♙♙♙♙
♖♘♗♕♔♗♘♖`
	assert.Equal(t, expected, str)
}
