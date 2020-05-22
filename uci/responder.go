package uci

import (
	"fmt"
	"io"
	"strings"

	chess "github.com/GeorgeBills/chess/m/v2"
)

// Engine-to-GUI constants are tokens sent from the engine to the GUI.
const (
	etgID       = "id"       // sent to identify the engine
	etgIDName   = "name"     // e.g. "id name Shredder X.Y\n"
	etgIDAuthor = "author"   // e.g. "id author Stefan MK\n"
	etgUCIOK    = "uciok"    // the engine has sent all infos and is ready
	etgReadyOK  = "readyok"  // the engine is ready to accept new commands
	etgBestMove = "bestmove" // engine has stopped searching and found the best move
	etgInfo     = "info"     // engine wants to send information to the GUI
)

type Responser interface {
	Response() string
}

// NewResponder returns a new responder.
func NewResponder(responsech <-chan Responser, out io.Writer) *Responder {
	return &Responder{
		responsech,
		out,
	}
}

// Responder pulls responses off a channel and writes them to the writer.
type Responder struct {
	responsech <-chan Responser
	out        io.Writer
}

// WriteResponses pulls responses off the responsech and writes them to the
// writer.
func (r Responder) WriteResponses() {
	for {
		response := <-r.responsech
		fmt.Fprintln(r.out, response.Response())
	}
}

type responseID struct{ key, value string }

func (r responseID) Response() string {
	return strings.Join([]string{etgID, r.key, r.value}, " ")
}

type responseOK struct{}

func (r responseOK) Response() string { return etgUCIOK }

type responseIsReady struct{}

func (r responseIsReady) Response() string { return etgReadyOK }

type responseBestMove struct {
	move chess.FromToPromoter
}

func (r responseBestMove) Response() string {
	movestr := ToUCIN(r.move)
	return strings.Join([]string{etgBestMove, movestr}, " ")
}
