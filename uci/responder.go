package uci

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/GeorgeBills/chess"
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

type Response interface {
	Response() string
}

// NewResponder returns a new responder.
func NewResponder(responsech <-chan Response, out io.Writer, logw io.Writer) *Responder {
	return &Responder{
		responsech,
		out,
		log.New(logw, "responder:", log.LstdFlags),
	}
}

// Responder pulls responses off a channel and writes them to the writer.
type Responder struct {
	responsech <-chan Response
	out        io.Writer
	logger     *log.Logger
}

// WriteResponses pulls responses off the responsech and writes them to the
// writer.
func (r Responder) WriteResponses() {
	r.logger.Println("starting")
	for response := range r.responsech {
		fmt.Fprintln(r.out, response.Response())
	}
	r.logger.Println("finished")
}

type ResponseID struct{ Key, Value string }

func (r ResponseID) Response() string {
	return strings.Join([]string{etgID, r.Key, r.Value}, " ")
}

type ResponseOK struct{}

func (r ResponseOK) Response() string { return etgUCIOK }

type ResponseIsReady struct{}

func (r ResponseIsReady) Response() string { return etgReadyOK }

type ResponseBestMove struct {
	Move chess.FromToPromoter
}

func (r ResponseBestMove) Response() string {
	movestr := ToUCIN(r.Move)
	return strings.Join([]string{etgBestMove, movestr}, " ")
}

type ResponseSearchInformation struct {
	Depth uint8
}

func (r ResponseSearchInformation) Response() string {
	return strings.Join([]string{etgInfo, "depth", strconv.Itoa(int(r.Depth))}, " ")
}
