package main

import (
	"fmt"
	"os"

	"github.com/GeorgeBills/chess/m/v2/uci"
)

func main() {
	logf, err := os.Create("uci.log")
	if err != nil {
		fatal(err)
	}

	a := newAdapter(logf)
	u := uci.New(a, os.Stdin, os.Stdout, logf)
	wg := u.ParseExecuteRespond()
	wg.Wait()
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
