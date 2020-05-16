package uci_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GeorgeBills/chess/m/v2/uci"
	"github.com/GeorgeBills/chess/m/v2/uci/mocks"
	"github.com/stretchr/testify/assert"
)

const Name = "test-engine"
const Author = "George Bills"

func TestQuitBeforeUCI(t *testing.T) {
	const in = "quit"
	r := strings.NewReader(in)
	h := &mocks.Handler{TB: t}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	assert.Equal(t, "", w.String())
}

func TestUCIOK(t *testing.T) {
	const in = "uci\nquit"
	r := strings.NewReader(in)
	h := &mocks.Handler{
		TB: t,
		IdentifyFunc: func() (name, author string, rest map[string]string) {
			return Name, Author, nil
		},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	const expected = "id name test-engine\nid author George Bills\nuciok\n"
	assert.Equal(t, expected, w.String())
}

func TestExtraInformation(t *testing.T) {
	const in = "uci\nquit"
	r := strings.NewReader(in)
	h := &mocks.Handler{
		TB: t,
		IdentifyFunc: func() (name, author string, rest map[string]string) {
			return "super-chess", "Jane Smith", map[string]string{
				"version":      "v1.2.3",
				"release-date": "2020-05-16",
			}
		},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	const expected = "id name super-chess\nid author Jane Smith\nid release-date 2020-05-16\nid version v1.2.3\nuciok\n"
	assert.Equal(t, expected, w.String())
}

func TestNewGame(t *testing.T) {
	const in = "uci\nucinewgame\nquit"
	r := strings.NewReader(in)
	var calledNewGame bool
	h := &mocks.Handler{
		TB: t,
		IdentifyFunc: func() (name, author string, rest map[string]string) {
			return Name, Author, nil
		},
		NewGameFunc: func() {
			calledNewGame = true
		},
	}
	w := &bytes.Buffer{}
	p := uci.NewParser(h, r, w, ioutil.Discard)
	p.Run()
	const expected = "id name test-engine\nid author George Bills\nuciok\n"
	assert.Equal(t, expected, w.String())
	assert.True(t, calledNewGame)
}
