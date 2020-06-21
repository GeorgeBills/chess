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
	factory := NewGameFactory()
	h := internal.NewEventHandler(client, logger, factory)
	go internal.StreamEvents(client, eventch, logger)
	go lichess.HandleEvents(h, eventch)
	wg.Wait()
}
