package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/GeorgeBills/chess/cmd/lichess-bot/internal"
	"github.com/GeorgeBills/chess/lichess"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	upgrade := flag.Bool("upgrade", false, "irreversibly upgrade to a bot account")
	challengeai := flag.Uint("challenge-ai", 0, "challenge the lichess AI at the given level")

	flag.Parse()

	token := os.Getenv("TOKEN")
	if token == "" {
		logger.Fatal(errors.New("TOKEN environment variable is required"))
	}

	transport := lichess.NewAuthorizingTransport(token, &http.Transport{})
	httpClient := &http.Client{
		Transport: transport,
	}
	client := lichess.NewClient(httpClient)

	switch {
	case *upgrade:

		err := client.BotUpgradeToBotAccount()
		if err != nil {
			logger.Fatal(fmt.Errorf("error upgrading to bot account: %w", err))
		}
		logger.Println("upgraded to bot account")

	case *challengeai > 0:

		eventch := make(chan interface{}, 100)

		// FIXME: pass wait groups down the stack

		wg := sync.WaitGroup{}
		wg.Add(2)
		factory := NewGameFactory()
		h := internal.NewEventHandler(client, logger, factory)
		go internal.StreamEvents(client, eventch, logger)
		go lichess.HandleEvents(h, eventch)

		logger.Printf("challenging AI level %d", *challengeai)
		var seconds uint = 600
		err := client.ChallengeAI(
			lichess.ParamsChallengeAI{
				Level: *challengeai,
				ParamsChallenge: lichess.ParamsChallenge{
					Rated:             false,
					ClockLimitSeconds: &seconds,
				},
			},
		)
		if err != nil {
			logger.Fatal(err)
		}

		wg.Wait()

	default:

		eventch := make(chan interface{}, 100)

		// FIXME: pass wait groups down the stack

		wg := sync.WaitGroup{}
		wg.Add(2)
		factory := NewGameFactory()
		h := internal.NewEventHandler(client, logger, factory)
		go internal.StreamEvents(client, eventch, logger)
		go lichess.HandleEvents(h, eventch)
		wg.Wait()
	}
}
