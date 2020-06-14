package lichess

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/GeorgeBills/chess"
	"github.com/GeorgeBills/chess/uci"
)

const (
	endpoint = "https://lichess.org"

	// "If you receive an HTTP response with a 429 status, please wait a full
	// minute before resuming API usage."
	backoff = 1 * time.Minute
)

//go:generate moq -out mocks/getposter.go -pkg mocks . GetPoster

type GetPoster interface {
	Get(uri string) (*http.Response, error)
	PostForm(uri string, data url.Values) (*http.Response, error)
}

func NewClient(httpClient GetPoster) *Client {
	return &Client{
		httpClient,
	}
}

type Client struct {
	httpClient GetPoster
}

// https://lichess.org/api#operation/botAccountUpgrade
func (c *Client) BotUpgradeToBotAccount() error {
	const path = "/api/bot/account/upgrade"

	uri := endpoint + path
	resp, err := c.httpClient.PostForm(uri, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return newBadRequestError(resp.Body)
	default:
		return newUnexpectedStatusCodeError(resp.StatusCode)
	}
}

type Event struct {
	EventType string          `json:"type"`
	Game      *EventGame      `json:"game"`
	Challenge *EventChallenge `json:"challenge"`
}

type EventGame struct {
	ID string `json:"id"`
}

type EventChallenge struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// https://lichess.org/api#operation/apiStreamEvent
func (c *Client) BotStreamEvents(eventch chan<- *Event) error {
	const path = "/api/stream/event"

	uri := endpoint + path
	resp, err := c.httpClient.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" { // tickle
			continue
		}

		var event Event
		err := json.Unmarshal([]byte(text), &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling '%s': %w", text, err)
		}

		eventch <- &event
	}

	return scanner.Err()
}

type GameEvent struct {
	GameEventType string `json:"type"`
	ID            string `json:"id"`
	Rated         bool   `json:"rated"`
}

// https://lichess.org/api#operation/botGameStream
func (c *Client) BotStreamGame(eventch chan<- *GameEvent) error {
	const path = "/api/bot/game/stream/{gameId}"

	uri := endpoint + path
	resp, err := c.httpClient.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "" { // tickle
			continue
		}

		var event GameEvent
		err := json.Unmarshal([]byte(text), &event)
		if err != nil {
			return fmt.Errorf("error unmarshalling '%s': %w", text, err)
		}

		eventch <- &event
	}

	return scanner.Err()
}

// https://lichess.org/api#operation/botGameMove
func (c *Client) BotMakeMove(gameID string, move chess.FromToPromoter, offeringDraw bool) error {
	const path = "/api/bot/game/%s/move/%s?offeringDraw=%t"

	uri := endpoint + fmt.Sprintf(path, gameID, uci.ToUCIN(move), offeringDraw)
	resp, err := c.httpClient.PostForm(uri, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return newBadRequestError(resp.Body)
	default:
		return newUnexpectedStatusCodeError(resp.StatusCode)
	}
}

type ChatRoom string

const (
	ChatRoomPlayer    = "player"
	ChatRoomSpectator = "spectator"
)

// https://lichess.org/api#operation/botGameChat
func (c *Client) BotWriteChat(gameID string, room ChatRoom, text string) error {
	const path = "/api/bot/game/%s/chat"

	uri := endpoint + fmt.Sprintf(path, gameID)
	resp, err := c.httpClient.PostForm(
		uri,
		url.Values{
			"room": {string(room)},
			"text": {text},
		},
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return newBadRequestError(resp.Body)
	default:
		return newUnexpectedStatusCodeError(resp.StatusCode)
	}
}

// https://lichess.org/api#operation/botGameAbort
func (c *Client) BotAbortGame(gameID string) error {
	const path = "/api/bot/game/%s/abort"

	uri := endpoint + fmt.Sprintf(path, gameID)
	resp, err := c.httpClient.PostForm(uri, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return newBadRequestError(resp.Body)
	default:
		return newUnexpectedStatusCodeError(resp.StatusCode)
	}
}

// https://lichess.org/api#operation/botGameResign
func (c *Client) BotResignGame(gameID string) error {
	const path = "/api/bot/game/%s/resign"

	uri := endpoint + fmt.Sprintf(path, gameID)
	resp, err := c.httpClient.PostForm(uri, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return newBadRequestError(resp.Body)
	default:
		return newUnexpectedStatusCodeError(resp.StatusCode)
	}
}

func newBadRequestError(body io.Reader) error {
	return errors.New("bad request") // TODO: parse body, should indicate exact error
}

func newUnexpectedStatusCodeError(code int) error {
	return fmt.Errorf("unexpected status code: %d", code)
}
