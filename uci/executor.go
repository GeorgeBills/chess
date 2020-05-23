package uci

import (
	"errors"
	"io"
	"log"
	"sort"
)

// NewExecutor returns a new executor.
func NewExecutor(commandch <-chan execer, stopch <-chan struct{}, a Adapter, responsech chan<- Responser, logw io.Writer) *Executor {
	return &Executor{
		commandch,
		stopch,
		log.New(logw, "executor:", log.LstdFlags),
		responsech,
		a,
	}
}

// Executor takes UCI commands from a channel and executes them.
type Executor struct {
	commandch  <-chan execer
	stopch     <-chan struct{}
	logger     *log.Logger
	responsech chan<- Responser
	adapter    Adapter
}

type execer interface {
	Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{})
}

// ExecuteCommands takes commands off commandch, executes them, and sends
// responses to responsech.
func (e *Executor) ExecuteCommands() {
	defer close(e.responsech)
	for cmd := range e.commandch {
		e.logger.Printf("running command: %T; %+v", cmd, cmd)
		// TODO: can we fan-in both stopch and quitch into one ch for exec?
		//       simplifies life for callees
		cmd.Exec(e.adapter, e.responsech, e.stopch)
		e.logger.Printf("finished command: %T", cmd)
	}
	e.logger.Println("finished")
}

type cmdUCI struct{}

func (c cmdUCI) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	name, author, rest := a.Identify()

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

func (c cmdNewGame) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.NewGame()
}

type cmdIsReady struct{}

func (c cmdIsReady) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.IsReady() // block on adapter
	responsech <- responseIsReady{}
}

type cmdSetStartingPosition struct{}

func (c cmdSetStartingPosition) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.SetStartingPosition()
}

type cmdSetPositionFEN struct {
	fen string
}

func (c cmdSetPositionFEN) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.SetPositionFEN(c.fen)
}

type cmdApplyMove struct {
	move Move
}

func (c cmdApplyMove) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.ApplyMove(c.move)
}

type cmdGoNodes struct {
	nodes uint64
}

func (c cmdGoNodes) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	movestr := a.GoNodes(c.nodes)
	responsech <- responseBestMove{movestr}
}

type cmdGoDepth struct {
	plies uint8
}

func (c cmdGoDepth) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	movestr := a.GoDepth(c.plies)
	responsech <- responseBestMove{movestr}
}

type cmdGoInfinite struct{}

func (c cmdGoInfinite) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	a.GoInfinite(stopch, responsech)
}

type cmdGoTime struct {
	tc TimeControl
}

func (c cmdGoTime) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	movestr := a.GoTime(c.tc)
	responsech <- responseBestMove{movestr}
}

type stopCommand struct{}

func (c stopCommand) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	panic(errors.New("stop not implemented"))
}

type quitCommand struct{}

func (c quitCommand) Exec(a Adapter, responsech chan<- Responser, stopch <-chan struct{}) {
	panic(errors.New("quit not implemented"))
}
