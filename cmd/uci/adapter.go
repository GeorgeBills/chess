package main

import (
	"errors"
	"io"
	"log"
	"strings"

	chess "github.com/GeorgeBills/chess/m/v2"
	"github.com/GeorgeBills/chess/m/v2/engine"
	"github.com/GeorgeBills/chess/m/v2/uci"
)

// Name is the name of our engine.
const Name = "github.com/GeorgeBills/chess"

// Author is the author of our engine.
const Author = "George Bills"

var errNoGame = errors.New("must initialise new game first")

// newAdapter returns a new adapter.
func newAdapter(logw io.Writer) *adapter {
	return &adapter{
		logger: log.New(logw, "adapter: ", log.LstdFlags),
	}
}

type adapter struct {
	logger *log.Logger
	game   *engine.Game
}

func (a *adapter) Identify() (name, author string, other map[string]string) {
	a.logger.Println("identify")
	return Name, Author, nil
}

func (a *adapter) NewGame() error {
	a.logger.Println("initialised new game")
	a.game = engine.NewGame(nil)
	return nil
}

func (a *adapter) SetStartingPosition(moves []chess.FromToPromoter) error {
	a.logger.Println("set starting position")

	if a.game == nil {
		return errNoGame
	}

	for _, move := range moves {
		m, err := a.game.HydrateMove(move)
		if err != nil {
			return err
		}
		a.game.MakeMove(m)
	}

	a.game.SetBoard(engine.NewBoard())
	return nil
}

func (a *adapter) SetPositionFEN(fen string, moves []chess.FromToPromoter) error {
	a.logger.Println("set position")

	if a.game == nil {
		return errNoGame
	}

	b, err := engine.NewBoardFromFEN(strings.NewReader(fen))
	if err != nil {
		return err
	}

	for _, move := range moves {
		m, err := a.game.HydrateMove(move)
		if err != nil {
			return err
		}
		a.game.MakeMove(m)
	}

	a.game.SetBoard(b)
	return nil
}

func (a *adapter) GoDepth(plies uint8, stopch <-chan struct{}, responsech chan<- uci.Responser) (chess.FromToPromoter, error) {
	a.logger.Println("go depth")

	if a.game == nil {
		return nil, errNoGame
	}

	statusch := make(chan engine.SearchStatus, 100)
	go forward(statusch, responsech)

	depth := 2 * plies // convert from full moves to half moves
	m, _ := a.game.BestMoveToDepth(depth, stopch, statusch)
	return m, nil
}

func (a *adapter) GoNodes(nodes uint64, stopch <-chan struct{}, responsech chan<- uci.Responser) (chess.FromToPromoter, error) {
	a.logger.Println("go nodes")

	if a.game == nil {
		return nil, errNoGame
	}

	return nil, errors.New("GoNodes not implemented")
}

func (a *adapter) GoInfinite(stopch <-chan struct{}, responsech chan<- uci.Responser) (chess.FromToPromoter, error) {
	a.logger.Println("go infinite")

	if a.game == nil {
		return nil, errNoGame
	}

	statusch := make(chan engine.SearchStatus, 100)
	go forward(statusch, responsech)

	m, _ := a.game.BestMoveInfinite(stopch, statusch)
	return m, nil
}

func (a *adapter) GoTime(tc uci.TimeControl, stopch <-chan struct{}, responsech chan<- uci.Responser) (chess.FromToPromoter, error) {
	a.logger.Println("go time")

	if a.game == nil {
		return nil, errNoGame
	}

	statusch := make(chan engine.SearchStatus, 100)
	go forward(statusch, responsech)

	m, _ := a.game.BestMoveToTime(tc.WhiteTime, tc.BlackTime, tc.WhiteIncrement, tc.BlackIncrement, stopch, statusch)
	return m, nil
}

// forward takes messages off statusch, converts them to uci responses and sends
// them off on responsech.
func forward(statusch <-chan engine.SearchStatus, responsech chan<- uci.Responser) {
	for info := range statusch {
		responsech <- uci.ResponseSearchInformation{Depth: info.Depth}
	}
}
