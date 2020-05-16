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
