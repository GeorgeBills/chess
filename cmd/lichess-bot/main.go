package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/GeorgeBills/chess/engine"
	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/uci"
)

const (
	id               = "gbcb"
	absDrawThreshold = 500 // relative to side to move!
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	token := flag.String("token", "", "personal API access token from https://lichess.org/account/oauth/token")
	upgrade := flag.Bool("upgrade", false, "irreversibly upgrade to a bot account")

	flag.Parse()

	if *token == "" {
		logger.Fatal(errors.New("token argument is required"))
	}

	transport := lichess.NewAuthorizingTransport(*token, &http.Transport{})
	httpClient := &http.Client{
		Transport: transport,
	}
	client := lichess.NewClient(httpClient)

	if *upgrade {
		err := client.BotUpgradeToBotAccount()
		if err != nil {
			logger.Fatal(fmt.Errorf("error upgrading to bot account: %w", err))
		}
		logger.Println("upgraded to bot account")
		return
	}

	eventch := make(chan interface{}, 100)

	// FIXME: pass wait groups down the stack

	wg := sync.WaitGroup{}
	wg.Add(2)
	h := &eventHandler{
		client: client,
	}
	go streamEvents(client, eventch)
	go lichess.HandleEvents(h, logger, eventch)
	wg.Wait()
}

type eventHandler struct {
	client *lichess.Client
}

func streamEvents(client *lichess.Client, eventch chan<- interface{}) {
	logger.Printf("streaming events")
	err := client.BotStreamEvents(eventch)
	if err != nil {
		logger.Fatal(fmt.Errorf("error streaming events: %w", err))
	}
}

func (h *eventHandler) Challenge(v *lichess.EventChallenge) {
	logger.Printf(
		"challenge; id: %s; challenger: %s; variant: %s; rated: %t",
		v.Challenge.ID,
		v.Challenge.Challenger.ID,
		v.Challenge.Variant.Key,
		v.Challenge.Rated,
	)

	if v.Challenge.Challenger.ID == "georgebills" &&
		v.Challenge.Rated == false && // require unrated for now to avoid changing my own rating
		v.Challenge.Variant.Key == lichess.VariantKeyStandard {

		logger.Printf("accepting challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeAccept(v.Challenge.ID)
		if err != nil {
			logger.Fatal(err)
		}
		// we now expect an incoming "game start" event
	} else {

		logger.Printf("declining challenge: %s", v.Challenge.ID)
		err := h.client.ChallengeDecline(v.Challenge.ID)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func (h *eventHandler) GameStart(v *lichess.EventGameStart) {
	eventch := make(chan interface{}, 100)
	gh := &gameHandler{
		gameID: v.Game.ID,
		client: h.client,
	}
	go streamGameEvents(v.Game.ID, h.client, eventch)
	go lichess.HandleGameEvents(gh, logger, eventch)
}

type Colour uint8

const (
	ColourUnknown Colour = iota
	ColourWhite
	ColourBlack
)

type gameHandler struct {
	game   *engine.Game
	gameID string
	client *lichess.Client
	colour Colour
}

func streamGameEvents(gameID string, client *lichess.Client, eventch chan<- interface{}) {
	logger.Printf("streaming game: %s", gameID)
	err := client.BotStreamGame(gameID, eventch)
	if err != nil {
		logger.Fatal(err)
	}
}

func (h *gameHandler) GameFull(v *lichess.EventGameFull) {
	logger.Printf("game full; game: %s; white: %s; black: %s", v.ID, v.White.ID, v.Black.ID)

	var b *engine.Board
	switch v.InitialFen {
	case "startpos":
		b = engine.NewBoard()
	default:
		var err error
		b, err = engine.NewBoardFromFEN(strings.NewReader(v.InitialFen))
		if err != nil {
			log.Fatal(err)
		}
	}

	h.game = engine.NewGame(b)

	switch id {
	case v.White.ID:
		h.colour = ColourWhite
	case v.Black.ID:
		h.colour = ColourBlack
	default:
		logger.Fatal("unknown colour to play")
	}
	logger.Println(h.colour)

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

func (h *gameHandler) GameState(v *lichess.EventGameState) {
	logger.Printf("game state")
	logger.Printf("status: %#v", v.Status)

	movestrs, tomove := splitMoves(v.Moves)

	// apply only the last move, the others have been applied already
	h.makeMove(movestrs[len(movestrs)-1])

	if tomove == h.colour {
		if err := h.searchMove(); err != nil {
			log.Fatal(err)
		}
	}
}

func (h *gameHandler) ChatLine(v *lichess.EventChatLine) {
	logger.Printf("chat line; '%s'", v.Text)
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

func (h *gameHandler) makeMove(movestr string) error {
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

func (h *gameHandler) searchMove() error {
	stopch := make(chan struct{})
	statusch := make(chan engine.SearchStatus)
	move, score := h.game.BestMoveToDepth(4, stopch, statusch)

	logger.Printf("making move %s", move.SAN())

	offerDraw := h.offerDraw(score)
	return h.client.BotMakeMove(h.gameID, move, offerDraw)
}

// offerDraw returns true if we should offer a draw and hope our opponent takes
// mercy on us.
func (h *gameHandler) offerDraw(score int16) bool {
	if h.colour == ColourWhite {
		return score < absDrawThreshold*-1
	}
	return score > absDrawThreshold
}
