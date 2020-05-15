package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Errorf("%s <board>", os.Args[0]))
	}

	bitstring := os.Args[1]

	var board uint64
	var err error

	// https://golang.org/pkg/strconv/#ParseInt
	//
	// "If the base argument is 0, the true base is implied by the string's
	// prefix: 2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise.
	// Also, for argument base 0 only, underscore characters are permitted as
	// defined by the Go syntax for integer literals."
	board, err = strconv.ParseUint(bitstring, 0, 64)
	if err != nil {
		fatal(err)
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	var i uint8
	for i = 0; i < 64; i++ {
		if i != 0 && i%8 == 0 {
			f.WriteRune('\n')
		}
		idx := engine.PrintOrderedIndex(i) // reverse ranks as we print
		if board&(1<<idx) != 0 {
			f.WriteRune('■')
		} else {
			f.WriteRune('□')
		}
	}
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}
