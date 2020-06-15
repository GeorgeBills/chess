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

func emptyok(uri string, data url.Values) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

type move struct {
	from, to  uint8
	promoteTo chess.PromoteTo
}

func (m move) From() uint8                { return m.from }
func (m move) To() uint8                  { return m.to }
func (m move) PromoteTo() chess.PromoteTo { return m.promoteTo }

func TestBotUpgradeToBotAccount(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
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

	eventch := make(chan interface{}, 100)

	err = c.BotStreamEvents(eventch)
	require.NoError(t, err)

	event1 := <-eventch // TODO: timeout
	expected1 := &lichess.EventChallenge{
		Challenge: &lichess.EventChallengeChallenge{
			ID:     "7pGLxJ4F",
			Status: "created",
			Rated:  true,
			Color:  "random",
			Challenger: &lichess.Player{
				ID:          "lovlas",
				Name:        "Lovlas",
				Provisional: false,
				Rating:      2506,
				Title:       "IM",
				Online:      true,
				Lag:         24,
			},
			DestinationUser: &lichess.Player{
				ID:          "thibot",
				Name:        "thibot",
				Provisional: true,
				Rating:      1500,
				Title:       "",
				Online:      true,
				Lag:         45,
			},
			Variant: &lichess.Variant{
				Key:   "standard",
				Name:  "Standard",
				Short: "Std",
			},
			Perf: &lichess.Perf{
				Icon: "#",
				Name: "Rapid",
			},
			TimeControl: &lichess.TimeControl{
				Type:      "clock",
				Limit:     300,
				Increment: 25,
				Show:      "5+25",
			},
		},
	}
	assert.Equal(t, expected1, event1)

	event2 := <-eventch // TODO: timeout
	expected2 := &lichess.EventGameStart{Game: &lichess.EventGameStartGame{ID: "1lsvP62l"}}
	assert.Equal(t, expected2, event2)

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

	eventch := make(chan interface{}, 100)

	err = c.BotStreamGame("9cttHZ", eventch)
	require.NoError(t, err)

	event1 := <-eventch // TODO: timeout
	expected1 := &lichess.EventGameFull{
		ID:         "5IrD6Gzz",
		Rated:      true,
		Variant:    &lichess.Variant{Key: "standard", Name: "Standard", Short: "Std"},
		Clock:      &lichess.EventGameFullClock{Initial: 1200000, Increment: 10000},
		Speed:      "classical",
		Perf:       &lichess.Perf{Name: "Classical"},
		CreatedAt:  1523825103562,
		InitialFen: "startpos",
		White:      &lichess.Player{ID: "lovlas", Name: "lovlas", Provisional: false, Rating: 2500, Title: "IM"},
		Black:      &lichess.Player{ID: "leela", Name: "leela", Provisional: false, Rating: 2390, Title: ""},
		State: &lichess.EventGameState{
			Moves:          "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7",
			WhiteTime:      7598040,
			BlackTime:      8395220,
			WhiteIncrement: 10000,
			BlackIncrement: 10000,
			Status:         "started",
		},
	}
	assert.Equal(t, expected1, event1)

	event2 := <-eventch // TODO: timeout
	expected2 := &lichess.EventGameState{
		Status:         "started",
		Moves:          "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7 b1c3",
		WhiteTime:      7598040,
		BlackTime:      8395220,
		WhiteIncrement: 10000,
		BlackIncrement: 10000,
	}
	assert.Equal(t, expected2, event2)

	event3 := <-eventch // TODO: timeout
	expected3 := &lichess.EventChatLine{Username: "thibault", Text: "Good luck, have fun", Room: "player"}
	assert.Equal(t, expected3, event3)

	event4 := <-eventch // TODO: timeout
	expected4 := &lichess.EventChatLine{Username: "lovlas", Text: "!eval", Room: "spectator"}
	assert.Equal(t, expected4, event4)

	event5 := <-eventch // TODO: timeout
	expected5 := &lichess.EventGameState{
		Status:         "resign",
		Moves:          "e2e4 c7c5 f2f4 d7d6 g1f3 b8c6 f1c4 g8f6 d2d3 g7g6 e1g1 f8g7 b1c3",
		WhiteTime:      7598040,
		BlackTime:      8395220,
		WhiteIncrement: 10000,
		BlackIncrement: 10000,
	}
	assert.Equal(t, expected5, event5)
}

func TestBotMakeMove(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.BotMakeMove("v8vDhD", move{12, 34, chess.PromoteToNone}, false)
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/v8vDhD/move/e2c5?offeringDraw=false", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestBotWriteChat(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.BotWriteChat("gLQEsv", lichess.ChatRoomPlayer, "ggwp!")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/gLQEsv/chat", calls[0].URI)
		assert.Equal(t, "room=player&text=ggwp%21", calls[0].Data.Encode())
	}
}

func TestBotAbortGame(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.BotAbortGame("DRBgmL")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/DRBgmL/abort", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestBotResignGame(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.BotResignGame("fN9LTy")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/bot/game/fN9LTy/resign", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestChallengeCreate(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	var limit, increment uint = 900, 10
	params := lichess.ChallengeCreateParams{
		Username:              "GeorgeBills",
		ClockLimitSeconds:     &limit,
		ClockIncrementSeconds: &increment,
	}
	err := c.ChallengeCreate(params)
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/challenge/GeorgeBills", calls[0].URI)
		assert.Equal(t, "clock.increment=10&clock.limit=900&rated=false", calls[0].Data.Encode())
	}
}

func TestChallengeAccept(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.ChallengeAccept("zbcLEG")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/challenge/zbcLEG/accept", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}

func TestChallengeDecline(t *testing.T) {
	m := &mocks.GetPosterMock{PostFormFunc: emptyok}
	c := lichess.NewClient(m)
	err := c.ChallengeDecline("PDCtHG")
	require.NoError(t, err)
	calls := m.PostFormCalls()
	if assert.Len(t, calls, 1) {
		assert.Equal(t, "https://lichess.org/api/challenge/PDCtHG/decline", calls[0].URI)
		assert.Equal(t, "", calls[0].Data.Encode())
	}
}
