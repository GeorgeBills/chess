package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/GeorgeBills/chess/lichess"
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

func NewGameHandler(gameID string, client Lichesser, logger *log.Logger, factory GameFactory) *GameHandler {
	return &GameHandler{
		gameID:  gameID,
		client:  client,
		logger:  logger,
		factory: factory,
	}
}

type GameHandler struct {
	game    Game
	gameID  string
	client  Lichesser
	colour  Colour
	logger  *log.Logger
	factory GameFactory
}

func (h *GameHandler) GameFull(v *lichess.EventGameFull) {
	h.logger.Printf("game full; game: %s; white: %s; black: %s", v.ID, v.White.ID, v.Black.ID)

	switch v.InitialFen {
	case "startpos":
		h.game = h.factory.NewGame()
	default:
		var err error
		h.game, err = h.factory.NewGameFromFEN(v.InitialFen)
		if err != nil {
			h.logger.Fatal(fmt.Errorf("error initialising game: %w", err))
		}
	}

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
		h.game.MakeMove(movestr)
	}

	// logger.Printf("variant: %#v", v.Variant)
	// logger.Printf("clock: %#v", v.Clock)
	// logger.Printf("perf: %#v", v.Perf)
	// logger.Printf("state: %#v", v.State)
	// logger.Printf("white: %#v", v.White)
	// logger.Printf("black: %#v", v.Black)

	// are we to move?
	if tomove == h.colour {
		h.bestMove()
	}
}

func (h *GameHandler) GameState(v *lichess.EventGameState) {
	h.logger.Printf("game state")
	h.logger.Printf("status: %#v", v.Status)

	movestrs, tomove := splitMoves(v.Moves)

	// apply only the last move, the others have been applied already
	h.game.MakeMove(movestrs[len(movestrs)-1])

	if tomove == h.colour {
		h.bestMove()
	}
}

func (h *GameHandler) bestMove() {
	move, score := h.game.BestMove()
	h.logger.Printf("making move %v", move)
	err := h.client.BotMakeMove(h.gameID, move, h.offerDraw(score))
	if err != nil {
		h.logger.Fatal(fmt.Errorf("error making move: %v", err))
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

// offerDraw returns true if we should offer a draw and hope our opponent takes
// mercy on us.
// FIXME: just resign here with a ggwp instead, no sense wasting our opponents time
// TODO: add logic to detect a "probable" draw (e.g. drawish material)
func (h *GameHandler) offerDraw(score int16) bool {
	if h.colour == ColourWhite {
		return score < absDrawThreshold*-1
	}
	return score > absDrawThreshold
}
