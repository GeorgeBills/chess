package uci_test

import (
	"testing"
	"time"

	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	processing = 10 * time.Millisecond
	timeout    = 1 * time.Second
)

func mustParseMove(ucin string) *uci.Move {
	m, err := uci.ParseUCIN(ucin)
	if err != nil {
		panic(err)
	}
	return m
}

func TestUCIN(t *testing.T) {
	tests := []string{
		"0000",
		"a1a2",
		"b3b4",
		"c5c6",
		"d7d8",
		"e7e8q",
		"f7f8r",
		"g7g8n",
		"h7h8b",
	}
	for _, ucin := range tests {
		t.Run(ucin, func(t *testing.T) {
			parsed, err := uci.ParseUCIN(ucin)
			require.NoError(t, err)
			str := uci.ToUCIN(parsed)
			assert.Equal(t, ucin, str)
		})
	}
}
