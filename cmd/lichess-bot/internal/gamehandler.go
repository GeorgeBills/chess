package internal

import (
	"fmt"
	"log"
	"math"
	"strings"

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
	h.logger.Printf(
		"game %s: full; %s; %s vs %s; t%s; %s",
		v.ID,
		strings.ToLower(v.Variant.Name),
		v.White.ID, v.Black.ID,
		timestr(v.Clock.Initial, v.Clock.Increment),
		ratedstr(v.Rated),
	)

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
		h.logger.Fatal("unknown colour to play") // FIXME: errch?
	}

	movestrs, tomove := splitMoves(v.State.Moves)
	for _, movestr := range movestrs {
		h.game.MakeMove(movestr)
	}

	// are we to move?
	if tomove == h.colour {
		h.bestMove()
	}
}

func (h *GameHandler) GameState(v *lichess.EventGameState) {
	h.logger.Printf(
		"game %s: %s; wt: %s; bt: %s; moves: %s",
		h.gameID,
		v.Status,
		timestr(v.WhiteTime, v.WhiteIncrement),
		timestr(v.BlackTime, v.BlackIncrement),
		last(20, v.Moves),
	)

	movestrs, tomove := splitMoves(v.Moves)

	// apply only the last move, the others have been applied already
	h.game.MakeMove(movestrs[len(movestrs)-1])

	if tomove == h.colour {
		h.bestMove()
	}
}

func (h *GameHandler) bestMove() {
	move, score := h.game.BestMove()
	h.logger.Printf("game %s: playing %s", h.gameID, uci.ToUCIN(move))
	err := h.client.BotMakeMove(h.gameID, move, h.offerDraw(score))
	if err != nil {
		h.logger.Fatal(fmt.Errorf("error playing move: %v", err))
	}
}

func (h *GameHandler) ChatLine(v *lichess.EventChatLine) {
	h.logger.Printf("game %s: chat line; %s", h.gameID, v.Text)
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

func timestr(initial, increment int) string {
	if initial == math.MaxInt32 {
		return "∞"
	}
	return fmt.Sprintf("%d+%d", initial, increment)
}

func ratedstr(rated bool) string {
	if rated {
		return "rated"
	}
	return "unrated"
}

func last(n int, s string) string {
	if n >= len(s) {
		return s
	}
	return "..." + s[len(s)-n:]
}
