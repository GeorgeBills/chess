package uci

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

// NewExecutor returns a new executor.
func NewExecutor(commandch <-chan execer, h Handler, responsech chan<- fmt.Stringer, logw io.Writer) *Executor {
	return &Executor{
		commandch,
		log.New(logw, "executor:", log.LstdFlags),
		responsech,
		h,
	}
}

// Executor takes UCI commands from a channel and executes them.
type Executor struct {
	commandch  <-chan execer
	logger     *log.Logger
	responsech chan<- fmt.Stringer
	handler    Handler
}

type execer interface {
	Exec(Handler, chan<- fmt.Stringer)
}

// ExecuteCommands takes commands off commandch, executes them, and sends
// responses to responsech.
func (e *Executor) ExecuteCommands() {
	for {
		cmd := <-e.commandch
		e.logger.Printf("running command %v", cmd)
		cmd.Exec(e.handler, e.responsech)
		e.logger.Printf("finished command %v", cmd)
	}
}

type cmdUCI struct{}

func (c cmdUCI) Exec(h Handler, responsech chan<- fmt.Stringer) {
	name, author, rest := h.Identify()

	// respond with required name and author
	responsech <- responseID{etgIDName, name}
	responsech <- responseID{etgIDAuthor, author}

	// respond with rest of our id information in sorted order
	keys := make([]string, 0, len(rest))
	for k := range rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		responsech <- responseID{k, rest[k]}
	}

	responsech <- responseOK{}
}

type cmdNewGame struct{}

func (c cmdNewGame) Exec(h Handler, responsech chan<- fmt.Stringer) {
	h.NewGame()
}

type cmdIsReady struct{}

func (c cmdIsReady) Exec(h Handler, responsech chan<- fmt.Stringer) {
	h.IsReady() // block on handler
	responsech <- responseIsReady{}
}

type cmdSetStartingPosition struct{}

func (c cmdSetStartingPosition) Exec(h Handler, responsech chan<- fmt.Stringer) {
	h.SetStartingPosition()
}

type cmdSetPositionFEN struct {
	fen string
}

func (c cmdSetPositionFEN) Exec(h Handler, responsech chan<- fmt.Stringer) {
	h.SetPositionFEN(c.fen)
}

type cmdApplyMove struct {
	move engine.FromToPromote
}

func (c cmdApplyMove) Exec(h Handler, responsech chan<- fmt.Stringer) {
	h.ApplyMove(c.move)
}

type cmdGoNodes struct {
	nodes uint64
}

func (c cmdGoNodes) Exec(h Handler, responsech chan<- fmt.Stringer) {
	movestr := h.GoNodes(c.nodes)
	responsech <- responseBestMove{movestr}
}

type cmdGoDepth struct {
	plies uint8
}

func (c cmdGoDepth) Exec(h Handler, responsech chan<- fmt.Stringer) {
	movestr := h.GoDepth(c.plies)
	responsech <- responseBestMove{movestr}
}

type cmdGoInfinite struct{}

func (c cmdGoInfinite) Exec(h Handler, responsech chan<- fmt.Stringer) {
	// TODO: obviously going to block forever; need to do channel plumbing
	h.GoInfinite()
}

type cmdGoTime struct {
	tc TimeControl
}

func (c cmdGoTime) Exec(h Handler, responsech chan<- fmt.Stringer) {
	movestr := h.GoTime(c.tc)
	responsech <- responseBestMove{movestr}
}

type stopCommand struct{}

func (c stopCommand) Exec(h Handler, responsech chan<- fmt.Stringer) {
	panic(errors.New("stop not implemented"))
}

type quitCommand struct{}

func (c quitCommand) Exec(h Handler, responsech chan<- fmt.Stringer) {
	panic(errors.New("quit not implemented"))
}
