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
)

const (
	id            = "gbcb"
	drawThreshold = -300
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
		client:  client,
		eventch: eventch,
	}
	go h.handleEvents()
	go h.streamEvents()
	wg.Wait()
}

type eventHandler struct {
	client  *lichess.Client
	eventch chan interface{}
}

func (h *eventHandler) streamEvents() {
	logger.Printf("streaming events")
	err := h.client.BotStreamEvents(h.eventch)
	if err != nil {
		logger.Fatal(fmt.Errorf("error streaming events: %w", err))
	}
}

func (h *eventHandler) handleEvents() {
	for event := range h.eventch {
		logger.Printf("event: %#v", event)

		switch v := event.(type) {

		case *lichess.EventChallenge:
			logger.Printf("accepting challenge: %s", v.Challenge.ID)

			if v.Challenge.Challenger.ID == "GeorgeBills" &&
				v.Challenge.Variant.Name == lichess.VariantNameStandard {
				err := h.client.ChallengeAccept(v.Challenge.ID)
				if err != nil {
					logger.Fatal(err)
				}
				// we now expect an incoming "game start" event
			} else {
				err := h.client.ChallengeDecline(v.Challenge.ID)
				if err != nil {
					logger.Fatal(err)
				}
			}

		case *lichess.EventGameStart:
			eventch := make(chan interface{}, 100)
			h := &gameHandler{
				gameID:  v.Game.ID,
				client:  h.client,
				eventch: eventch,
			}
			go h.streamGameEvents()
			go h.handleGameEvents()

		default:
			logger.Printf("ignoring unrecognized event type: %T", v)
		}
	}
}

type gameHandler struct {
	game    *engine.Game
	gameID  string
	client  *lichess.Client
	eventch chan interface{}
}

func (h *gameHandler) streamGameEvents() {
	logger.Printf("streaming game events")
	err := h.client.BotStreamGame(h.gameID, h.eventch)
	if err != nil {
		logger.Fatal(err)
	}
}

func (h *gameHandler) handleGameEvents() {
	for event := range h.eventch {
		logger.Printf("game event: %#v", event) // TODO: include game id

		switch v := event.(type) {
		case *lichess.EventGameFull:
			logger.Printf("new game: %s", v.ID)
			// TODO: refuse to play if the variant isn't standard

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

			logger.Printf("state: %#v", v.State)
			logger.Printf("white: %#v", v.White)
			logger.Printf("black: %#v", v.Black)

			// are we to move?
			if v.White.ID == id {
				stopch := make(chan struct{})
				statusch := make(chan engine.SearchStatus)
				move, score := h.game.BestMoveToDepth(4, stopch, statusch)
				offeringDraw := score <= drawThreshold
				err := h.client.BotMakeMove(h.gameID, move, offeringDraw)
				if err != nil {
					log.Fatal(err)
				}
			}

		case *lichess.EventGameState:
			logger.Printf("game state")

		case *lichess.EventChatLine:
			logger.Printf("chat line")

		default:
			logger.Printf("ignoring unrecognized event type: %T", v)
		}
	}
}
