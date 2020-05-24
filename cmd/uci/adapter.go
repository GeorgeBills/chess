package main

import (
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

func (a *adapter) NewGame() {
	a.logger.Println("initialised new game")
	g := engine.NewGame(nil) // TODO: return pointer
	a.game = &g
}

func (a *adapter) SetStartingPosition() {
	a.logger.Println("set starting position")
	// TODO: nil check game on SetPositionFEN
	//       or just return a new game from SetBoard if none already?
	// if a.game == nil {
	// 	return errors.New("no game")
	// }
	b := engine.NewBoard()
	a.game.SetBoard(&b)
}

func (a *adapter) SetPositionFEN(fen string) {
	a.logger.Println("set position")

	// TODO: nil check game on SetPositionFEN
	//       or just return a new game from SetBoard if none already?

	b, _ := engine.NewBoardFromFEN(strings.NewReader(fen))
	// TODO: return err from SetPositionFEN
	// if err != nil {
	// 	return err
	// }

	a.game.SetBoard(b)
}

func (a *adapter) ApplyMove(move chess.FromToPromoter) {
	a.logger.Printf("playing move: %v", move)
	m, err := a.game.HydrateMove(move)
	if err != nil {
		panic(err) // FIXME: return errors from most adapter methods...
	}
	a.game.MakeMove(m)
}

func (a *adapter) GoDepth(plies uint8) chess.FromToPromoter {
	a.logger.Println("go depth")
	m, _ := a.game.BestMoveToDepth(plies * 2)
	return m
}

func (a *adapter) GoNodes(nodes uint64) chess.FromToPromoter {
	a.logger.Println("go nodes")
	panic("GoNodes not implemented")
}

func (a *adapter) GoInfinite(stopch <-chan struct{}, responsech chan<- uci.Responser) chess.FromToPromoter {
	a.logger.Println("go infinite")

	statusch := make(chan engine.SearchStatus, 100)
	go func() {
		for info := range statusch {
			responsech <- uci.ResponseSearchInformation{Depth: info.Depth}
		}
	}()

	m, _ := a.game.BestMoveInfinite(stopch, statusch)
	return m
}

func (a *adapter) GoTime(tc uci.TimeControl) chess.FromToPromoter {
	a.logger.Println("go time")
	m, _ := a.game.BestMoveToTime(tc.WhiteTime, tc.BlackTime, tc.WhiteIncrement, tc.BlackIncrement)
	return m
}
