package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/GeorgeBills/chess/lichess"
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

func streamEvents(client *lichess.Client, eventch chan<- interface{}) {
	logger.Printf("streaming events")
	err := client.BotStreamEvents(eventch)
	if err != nil {
		logger.Fatal(fmt.Errorf("error streaming events: %w", err))
	}
}
