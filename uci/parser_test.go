package uci_test

import (
	"io"
	"os"
	"testing"
	"time"

	chess "github.com/GeorgeBills/chess/m/v2"
	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	processing = 10 * time.Millisecond
	timeout    = 1 * time.Second
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []uci.Execer
	}{
		{
			"quit before uci",
			[]string{"quit"},
			nil,
		},
		{
			"uci",
			[]string{"uci", "quit"},
			[]uci.Execer{uci.CommandUCI{}},
		},
		{
			"isready",
			[]string{"uci", "isready", "quit"},
			[]uci.Execer{uci.CommandUCI{}, uci.CommandIsReady{}},
		},
		{
			"ucinewgame",
			[]string{"uci", "ucinewgame", "quit"},
			[]uci.Execer{uci.CommandUCI{}, uci.CommandNewGame{}},
		},
		{
			"position fen",
			[]string{
				"uci", "ucinewgame",
				"position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
				"quit",
			},
			[]uci.Execer{
				uci.CommandUCI{},
				uci.CommandNewGame{},
				&uci.CommandSetPositionFEN{
					FEN:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
					Moves: nil,
				},
			},
		},
		{
			"position startpos",
			[]string{"uci", "ucinewgame", "position startpos", "quit"},
			[]uci.Execer{
				uci.CommandUCI{}, uci.CommandNewGame{}, &uci.CommandSetStartingPosition{},
			},
		},
		{
			"position startpos moves",
			[]string{
				"uci", "ucinewgame",
				"position startpos moves e2e4 e7e5 g1f3 b8c6 f1b5", // Ruy López
				"quit",
			},
			[]uci.Execer{
				uci.CommandUCI{}, uci.CommandNewGame{},
				&uci.CommandSetStartingPosition{
					Moves: []chess.FromToPromoter{
						mustParseMove("e2e4"),
						mustParseMove("e7e5"),
						mustParseMove("g1f3"),
						mustParseMove("b8c6"),
						mustParseMove("f1b5"),
					},
				},
			},
		},
		{
			"go depth",
			[]string{
				"uci", "ucinewgame", "position startpos",
				"go depth 123",
				"quit",
			},
			[]uci.Execer{
				uci.CommandUCI{}, uci.CommandNewGame{}, &uci.CommandSetStartingPosition{},
				uci.CommandGoDepth{Plies: 123},
			},
		},
		{
			"go nodes",
			[]string{
				"uci", "ucinewgame", "position startpos",
				"go nodes 456",
				"quit",
			},
			[]uci.Execer{
				uci.CommandUCI{}, uci.CommandNewGame{}, &uci.CommandSetStartingPosition{},
				uci.CommandGoNodes{Nodes: 456},
			},
		},
		{
			"go time",
			[]string{
				"uci", "ucinewgame", "position startpos",
				"go wtime 300000 btime 300000 winc 0 binc 0",
				"quit",
			},
			[]uci.Execer{
				uci.CommandUCI{}, uci.CommandNewGame{}, &uci.CommandSetStartingPosition{},
				uci.CommandGoTime{
					uci.TimeControl{
						WhiteTime: 5 * time.Minute,
						BlackTime: 5 * time.Minute,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			piper, pipew := io.Pipe()

			p, commandch, stopch := uci.NewParser(piper, os.Stdout)
			require.NotNil(t, p)
			require.NotNil(t, commandch)
			require.NotNil(t, stopch)

			go p.ParseInput()

			var commands []uci.Execer
			go func() {
				for cmd := range commandch {
					commands = append(commands, cmd)
				}
			}()

			for _, str := range tt.input {
				pipew.Write([]byte(str + "\n"))
				time.Sleep(processing)
			}

			assert.Equal(t, tt.expected, commands)

			// should be no more commands, and the channel should be closed
			select {
			case cmd, open := <-commandch:
				assert.Equal(t, nil, cmd)
				assert.False(t, open)
			default:
				t.Errorf("channel is still open") // closed channel wouldn't block
			}
		})
	}
}
