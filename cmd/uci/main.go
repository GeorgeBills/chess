package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/GeorgeBills/chess/uci"
)

func main() {
	logf, err := os.Create("uci.log")
	if err != nil {
		fatal(err)
	}

	adapter := newAdapter(logf)
	parser, commandch, stopch := uci.NewParser(os.Stdin, logf)
	executor, responsech := uci.NewExecutor(commandch, stopch, adapter, logf)
	responder := uci.NewResponder(responsech, os.Stdout, logf)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		parser.ParseInput()
		wg.Done()
	}()

	go func() {
		executor.ExecuteCommands()
		wg.Done()
	}()

	go func() {
		responder.WriteResponses()
		wg.Done()
	}()

	wg.Wait()
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
