package uci

import (
	"io"
	"log"
	"sort"

	chess "github.com/GeorgeBills/chess/m/v2"
)

// NewExecutor returns a new executor.
func NewExecutor(commandch <-chan Command, stopch <-chan struct{}, a Adapter, logw io.Writer) (*Executor, <-chan Response) {
	responsech := make(chan Response, 100)
	executor := &Executor{
		commandch,
		stopch,
		log.New(logw, "executor:", log.LstdFlags),
		responsech,
		a,
	}
	return executor, responsech
}

// Executor takes UCI commands from a channel and executes them.
type Executor struct {
	commandch  <-chan Command
	stopch     <-chan struct{}
	logger     *log.Logger
	responsech chan<- Response
	adapter    Adapter
}

type Command interface {
	Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error
}

type AppendMoveCommand interface {
	Command
	AppendMove(m *Move)
}

// ExecuteCommands takes commands off commandch, executes them, and sends
// responses to responsech.
func (e *Executor) ExecuteCommands() {
	defer close(e.responsech)

	e.logger.Println("starting")
	for cmd := range e.commandch {
		e.logger.Printf("running command: %T; %+v", cmd, cmd)
		// TODO: can we fan-in both stopch and quitch into one ch for exec?
		//       simplifies life for callees
		err := cmd.Execute(e.adapter, e.responsech, e.stopch)
		if err != nil {
			e.logger.Printf("error running command: %v", err)
		}
		e.logger.Printf("finished command: %T", cmd)
	}
	e.logger.Println("finished")
}

type CommandUCI struct{}

func (c CommandUCI) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	name, author, rest := a.Identify()

	// respond with required name and author
	responsech <- ResponseID{etgIDName, name}
	responsech <- ResponseID{etgIDAuthor, author}

	// respond with rest of our id information in sorted order
	keys := make([]string, 0, len(rest))
	for k := range rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		responsech <- ResponseID{k, rest[k]}
	}

	responsech <- ResponseOK{}

	return nil
}

type CommandNewGame struct{}

func (c CommandNewGame) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	return a.NewGame()
}

type CommandIsReady struct{}

func (c CommandIsReady) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	responsech <- ResponseIsReady{}
	return nil
}

type CommandSetStartingPosition struct {
	Moves []chess.FromToPromoter
}

func (c *CommandSetStartingPosition) AppendMove(m *Move) {
	c.Moves = append(c.Moves, m)
}

func (c CommandSetStartingPosition) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	return a.SetStartingPosition(c.Moves)
}

type CommandSetPositionFEN struct {
	FEN   string
	Moves []chess.FromToPromoter
}

func (c CommandSetPositionFEN) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	return a.SetPositionFEN(c.FEN, c.Moves)
}

func (c *CommandSetPositionFEN) AppendMove(m *Move) {
	c.Moves = append(c.Moves, m)
}

type CommandGoNodes struct {
	Nodes uint64
}

func (c CommandGoNodes) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	movestr, err := a.GoNodes(c.Nodes, stopch, responsech)
	if err != nil {
		return err
	}
	responsech <- ResponseBestMove{movestr}
	return nil
}

type CommandGoDepth struct {
	Plies uint8
}

func (c CommandGoDepth) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	movestr, err := a.GoDepth(c.Plies, stopch, responsech)
	if err != nil {
		return err
	}
	responsech <- ResponseBestMove{movestr}
	return nil
}

type CommandGoInfinite struct{}

func (c CommandGoInfinite) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	move, err := a.GoInfinite(stopch, responsech)
	if err != nil {
		return err
	}
	responsech <- ResponseBestMove{move}
	return nil
}

type CommandGoTime struct {
	TimeControl
}

func (c CommandGoTime) Execute(a Adapter, responsech chan<- Response, stopch <-chan struct{}) error {
	move, err := a.GoTime(c.TimeControl, stopch, responsech)
	if err != nil {
		return err
	}
	responsech <- ResponseBestMove{move}
	return nil
}
