package uci_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponses(t *testing.T) {

	tests := []struct {
		name     string
		response uci.Responser
		expected string
	}{
		{
			"id",
			uci.ResponseID{Key: "author", Value: Author},
			"id author George Bills\n",
		},
		{
			"ok",
			uci.ResponseOK{},
			"uciok\n",
		},
		{
			"readyok",
			uci.ResponseIsReady{},
			"readyok\n",
		},
		{
			"bestmove",
			uci.ResponseBestMove{Move: mustParseMove("a1h8")},
			"bestmove a1h8\n",
		},
		{
			"info",
			uci.ResponseSearchInformation{Depth: 123},
			"info depth 123\n",
		},
	}

	responsech := make(chan uci.Responser)
	var buf bytes.Buffer
	r := uci.NewResponder(responsech, &buf, os.Stdout)

	go r.WriteResponses()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responsech <- tt.response
			time.Sleep(processing)
			assert.Equal(t, tt.expected, buf.String())
			buf.Reset()
		})
	}

	close(responsech)
}
