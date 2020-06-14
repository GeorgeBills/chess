package lichess_test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess"
	"github.com/GeorgeBills/chess/lichess"
	"github.com/GeorgeBills/chess/lichess/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type move struct {
	from, to  uint8
	promoteTo chess.PromoteTo
}

func (m move) From() uint8                { return m.from }
func (m move) To() uint8                  { return m.to }
func (m move) PromoteTo() chess.PromoteTo { return m.promoteTo }

func TestBotUpgradeToBotAccount(t *testing.T) {
	m := &mocks.GetPosterMock{
		PostFormFunc: func(uri string, data url.Values) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	c := lichess.NewClient(m)
	err := c.BotUpgradeToBotAccount()
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/account/upgrade", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestBotStreamEvents(t *testing.T) {
	f, err := os.Open("testdata/event-stream.ndjson")
	require.NoError(t, err)

	m := &mocks.GetPosterMock{
		GetFunc: func(uri string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       f,
			}, nil
		},
	}
	c := lichess.NewClient(m)

	eventch := make(chan *lichess.Event, 100)

	err = c.BotStreamEvents(eventch)
	require.NoError(t, err)

	event := <-eventch // TODO: timeout
	expected := &lichess.Event{
		EventType: "challenge",
		Challenge: &lichess.EventChallenge{
			ID:     "7pGLxJ4F",
			Status: "created",
		},
	}
	assert.Equal(t, expected, event)
}

func TestBotStreamGame(t *testing.T) {
	f, err := os.Open("testdata/game-stream.ndjson")
	require.NoError(t, err)

	m := &mocks.GetPosterMock{
		GetFunc: func(uri string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       f,
			}, nil
		},
	}
	c := lichess.NewClient(m)

	eventch := make(chan *lichess.GameEvent, 100)

	err = c.BotStreamGame(eventch)
	require.NoError(t, err)

	event := <-eventch // TODO: timeout
	expected := &lichess.GameEvent{
		GameEventType: "gameFull",
		ID:            "5IrD6Gzz",
		Rated:         true,
	}
	assert.Equal(t, expected, event)
}

func TestBotMakeMove(t *testing.T) {
	m := &mocks.GetPosterMock{
		PostFormFunc: func(uri string, data url.Values) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	c := lichess.NewClient(m)
	err := c.BotMakeMove("abc123", move{12, 34, chess.PromoteToNone}, false)
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/abc123/move/e2c5?offeringDraw=false", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestBotWriteChat(t *testing.T) {
	m := &mocks.GetPosterMock{
		PostFormFunc: func(uri string, data url.Values) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	c := lichess.NewClient(m)
	err := c.BotWriteChat("abc123", lichess.ChatRoomPlayer, "ggwp!")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/abc123/chat", calls[0].URI)
		assert.Equal(t, "room=player&text=ggwp%21", calls[0].Data.Encode())
	}
}

func TestBotAbortGame(t *testing.T) {
	m := &mocks.GetPosterMock{
		PostFormFunc: func(uri string, data url.Values) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	c := lichess.NewClient(m)
	err := c.BotAbortGame("abc123")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/abc123/abort", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestBotResignGame(t *testing.T) {
	m := &mocks.GetPosterMock{
		PostFormFunc: func(uri string, data url.Values) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
	c := lichess.NewClient(m)
	err := c.BotResignGame("abc123")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/abc123/resign", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}
