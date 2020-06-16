package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/GeorgeBills/chess/engine"
	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/uci"
)

const (
	id               = "gbcb"
	absDrawThreshold = 500 // relative to side to move!
)

type Colour uint8

const (
	ColourUnknown Colour = iota
	ColourWhite
	ColourBlack
)

func NewGameHandler(gameID string, client Lichesser, logger *log.Logger) *GameHandler {
	return &GameHandler{
		gameID: gameID,
		client: client,
		logger: logger,
		// TODO: inject engine.NewGame, engine.NewBoard and engine.NewBoardFromFEN
	}
}

type GameHandler struct {
	game   *engine.Game
	gameID string
	client Lichesser
	colour Colour
	logger *log.Logger
}

func (h *GameHandler) GameFull(v *lichess.EventGameFull) {
	h.logger.Printf("game full; game: %s; white: %s; black: %s", v.ID, v.White.ID, v.Black.ID)

	var b *engine.Board
	switch v.InitialFen {
	case "startpos":
		b = engine.NewBoard()
	default:
		var err error
		b, err = engine.NewBoardFromFEN(strings.NewReader(v.InitialFen))
		if err != nil {
			log.Fatal(fmt.Errorf("error while parsing FEN: %w", err))
		}
	}

	h.game = engine.NewGame(b)

	switch id {
	case v.White.ID:
		h.colour = ColourWhite
	case v.Black.ID:
		h.colour = ColourBlack
	default:
		h.logger.Fatal("unknown colour to play")
	}
	h.logger.Println(h.colour)

	movestrs, tomove := splitMoves(v.State.Moves)
	for _, movestr := range movestrs {
		h.makeMove(movestr)
	}

	// logger.Printf("variant: %#v", v.Variant)
	// logger.Printf("clock: %#v", v.Clock)
	// logger.Printf("perf: %#v", v.Perf)
	// logger.Printf("state: %#v", v.State)
	// logger.Printf("white: %#v", v.White)
	// logger.Printf("black: %#v", v.Black)

	// are we to move?
	if tomove == h.colour {
		if err := h.searchMove(); err != nil {
			log.Fatal(err)
		}
	}
}

func (h *GameHandler) GameState(v *lichess.EventGameState) {
	h.logger.Printf("game state")
	h.logger.Printf("status: %#v", v.Status)

	movestrs, tomove := splitMoves(v.Moves)

	// apply only the last move, the others have been applied already
	h.makeMove(movestrs[len(movestrs)-1])

	if tomove == h.colour {
		if err := h.searchMove(); err != nil {
			log.Fatal(err)
		}
	}
}

func (h *GameHandler) ChatLine(v *lichess.EventChatLine) {
	h.logger.Printf("chat line; '%s'", v.Text)
}

func splitMoves(moves string) ([]string, Colour) {
	if moves == "" {
		return []string{}, ColourWhite
	}

	movestrs := strings.Split(moves, " ")
	if len(movestrs)%2 == 0 {
		return movestrs, ColourWhite
	}
	return movestrs, ColourBlack
}

func (h *GameHandler) makeMove(movestr string) error {
	parsed, err := uci.ParseUCIN(movestr)
	if err != nil {
		return fmt.Errorf("error applying move: %w", err)
	}

	move, err := h.game.HydrateMove(parsed)
	if err != nil {
		return fmt.Errorf("error making move: %w", err)
	}

	h.game.MakeMove(move)
	return nil
}

func (h *GameHandler) searchMove() error {
	stopch := make(chan struct{})
	statusch := make(chan engine.SearchStatus)
	move, score := h.game.BestMoveToDepth(4, stopch, statusch)

	h.logger.Printf("making move %s", move.SAN())

	offerDraw := h.offerDraw(score)
	return h.client.BotMakeMove(h.gameID, move, offerDraw)
}

// offerDraw returns true if we should offer a draw and hope our opponent takes
// mercy on us.
// FIXME: just resign here instead, no sense wasting our opponents time
// TODO: add logic to detect a "probable" draw (e.g. drawish material)
func (h *GameHandler) offerDraw(score int16) bool {
	if h.colour == ColourWhite {
		return score < absDrawThreshold*-1
	}
	return score > absDrawThreshold
}
