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
	"github.com/mitchellh/mapstructure"
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

type EventGameStart struct {
	Game *EventGameStartGame `json:"game" mapstructure:"game"`
}

type EventGameStartGame struct {
	ID string `json:"id" mapstructure:"id"`
}

type EventChallenge struct {
	Challenge *EventChallengeChallenge `json:"challenge" mapstructure:"challenge"`
}

type EventChallengeChallenge struct {
	ID     string `json:"id" mapstructure:"id"`
	Status string `json:"status" mapstructure:"status"`
}

type EventGameFull struct {
	ID         string                   `json:"id" mapstructure:"id"`
	Rated      bool                     `json:"rated" mapstructure:"rated"`
	Variant    *EventGameFullVariant    `json:"variant" mapstructure:"variant"`
	Clock      *EventGameFullClock      `json:"clock" mapstructure:"clock"`
	Speed      string                   `json:"speed" mapstructure:"speed"`
	Perf       *EventGameFullPerf       `json:"perf" mapstructure:"perf"`
	CreatedAt  int                      `json:"createdAt" mapstructure:"createdAt"` // TODO: should be time.Time
	White      *EventGameFullPlayerInfo `json:"white" mapstructure:"white"`
	Black      *EventGameFullPlayerInfo `json:"black" mapstructure:"black"`
	InitialFen string                   `json:"initialFen" mapstructure:"initialFen"`
	State      *EventGameState          `json:"state" mapstructure:"state"`
}

type EventGameFullVariant struct {
	Key   string `json:"key" mapstructure:"key"`
	Name  string `json:"name" mapstructure:"name"`
	Short string `json:"short" mapstructure:"short"`
}

type EventGameFullClock struct {
	Initial   int `json:"initial" mapstructure:"initial"`
	Increment int `json:"increment" mapstructure:"increment"`
}

type EventGameFullPerf struct {
	Name string `json:"name" mapstructure:"name"`
}

type EventGameFullPlayerInfo struct {
	ID          string `json:"id" mapstructure:"id"`
	Name        string `json:"name" mapstructure:"name"`
	Provisional bool   `json:"provisional" mapstructure:"provisional"`
	Rating      int    `json:"rating" mapstructure:"rating"`
	Title       string `json:"title" mapstructure:"title"`
}

type EventGameState struct {
	Moves          string `json:"moves" mapstructure:"moves"`
	WhiteTime      int    `json:"wtime" mapstructure:"wtime"`
	BlackTime      int    `json:"btime" mapstructure:"btime"`
	WhiteIncrement int    `json:"winc" mapstructure:"winc"`
	BlackIncrement int    `json:"binc" mapstructure:"binc"`
	Status         string `json:"status" mapstructure:"status"`
}

type EventChatLine struct {
	Username string `json:"username" mapstructure:"username"`
	Text     string `json:"text" mapstructure:"text"`
	Room     string `json:"room" mapstructure:"room"`
}

// https://lichess.org/api#operation/apiStreamEvent
func (c *Client) BotStreamEvents(eventch chan<- interface{}) error {
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

		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(text), &parsed)
		if err != nil {
			return fmt.Errorf("error unmarshalling '%s': %w", text, err)
		}

		switch parsed["type"] {
		case "challenge":
			var ec EventChallenge
			err := mapstructure.Decode(parsed, &ec)
			if err != nil {
				return fmt.Errorf("error mapping struct: %w", err)
			}
			eventch <- &ec
		case "gameStart":
			var egs EventGameStart
			err := mapstructure.Decode(parsed, &egs)
			if err != nil {
				return fmt.Errorf("error mapping struct: %w", err)
			}
			eventch <- &egs
		default:
			return fmt.Errorf("unrecognized event type: %s", parsed["type"])
		}
	}

	return scanner.Err()
}

// https://lichess.org/api#operation/botGameStream
func (c *Client) BotStreamGame(eventch chan<- interface{}) error {
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

		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(text), &parsed)
		if err != nil {
			return fmt.Errorf("error unmarshalling '%s': %w", text, err)
		}

		switch parsed["type"] {
		case "gameFull":
			var egf EventGameFull
			err := mapstructure.Decode(parsed, &egf)
			if err != nil {
				return fmt.Errorf("error mapping struct: %w", err)
			}
			eventch <- &egf
		case "gameState":
			var egs EventGameState
			err := mapstructure.Decode(parsed, &egs)
			if err != nil {
				return fmt.Errorf("error mapping struct: %w", err)
			}
			eventch <- &egs
		case "chatLine":
			var ecl EventChatLine
			err := mapstructure.Decode(parsed, &ecl)
			if err != nil {
				return fmt.Errorf("error mapping struct: %w", err)
			}
			eventch <- &ecl
		default:
			return fmt.Errorf("unrecognized event type: %s", parsed["type"])
		}
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
