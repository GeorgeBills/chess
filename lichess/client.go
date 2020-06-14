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
	EventType string `json:"type"`
	Game      struct {
		ID string `json:"id"`
	} `json:"game"`
	Challenge struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"challenge"`
}

// https://lichess.org/api#operation/apiStreamEvent
func (c *Client) BotStreamEvents() error {
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
		var event Event
		json.Unmarshal([]byte(text), &event)
	}

	return scanner.Err()
}

// [
//   {
//     "type": "gameFull",
//     "id": "5IrD6Gzz",
//     "rated": true,
//     "variant": {
//       "key": "standard",
//       "name": "Standard",
//       "short": "Std"
//     },
//     "clock": {
//       "initial": 1200000,
//       "increment": 10000
//     },
//     "speed": "classical",
//     "perf": {
//       "name": "Classical"
//     },
//     "createdAt": 1523825103562,
//     "white": {
//       "id": "lovlas",
//       "name": "lovlas",
//       "provisional": false,
//       "rating": 2500,
//       "title": "IM"
//     },
//     "black": {
//       "id": "leela",
//       "name": "leela",
//       "rating": 2390,
//       "title": null
//     },
//     "initialFen": "startpos",
//     "state": {
//       "type": "gameState",
//       "moves": "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7",
//       "wtime": 7598040,
//       "btime": 8395220,
//       "winc": 10000,
//       "binc": 10000,
//       "status": "started"
//     }
//   },

//   {
//     "type": "gameState",
//     "moves": "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7 b1c3",
//     "wtime": 7598040,
//     "btime": 8395220,
//     "winc": 10000,
//     "binc": 10000,
//     "status": "started"
//   },

//   {
//     "type": "chatLine",
//     "username": "thibault",
//     "text": "Good luck, have fun",
//     "room": "player"
//   },

type Game struct {
	GameType string `json:"type"`
	ID       string `json:"id"`
	Rated    bool   `json:"rated"`
}

// https://lichess.org/api#operation/apiStreamEvent
func (c *Client) BotStreamGame() error {
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
		var event Event
		json.Unmarshal([]byte(text), &event)
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
