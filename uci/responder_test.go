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
	responsech := make(chan uci.Responser)
	var buf bytes.Buffer
	r := uci.NewResponder(responsech, &buf, os.Stdout)

	go r.WriteResponses()

	t.Run("id", func(t *testing.T) {
		responsech <- uci.ResponseID{Key: "author", Value: Author}
		time.Sleep(processing)
		assert.Equal(t, "id author George Bills\n", buf.String())
		buf.Reset()
	})

	t.Run("ok", func(t *testing.T) {
		responsech <- uci.ResponseOK{}
		time.Sleep(processing)
		assert.Equal(t, "uciok\n", buf.String())
		buf.Reset()
	})

	t.Run("readyok", func(t *testing.T) {
		responsech <- uci.ResponseIsReady{}
		time.Sleep(processing)
		assert.Equal(t, "readyok\n", buf.String())
		buf.Reset()
	})

	t.Run("bestmove", func(t *testing.T) {
		responsech <- uci.ResponseBestMove{Move: mustParseMove("a1h8")}
		time.Sleep(processing)
		assert.Equal(t, "bestmove a1h8\n", buf.String())
		buf.Reset()
	})

	t.Run("info depth", func(t *testing.T) {
		responsech <- uci.ResponseSearchInformation{Depth: 123}
		time.Sleep(processing)
		assert.Equal(t, "info depth 123\n", buf.String())
		buf.Reset()
	})

	close(responsech)
}
