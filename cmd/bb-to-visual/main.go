package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Errorf("%s <board>", os.Args[0]))
	}

	bitstring := os.Args[1]

	var board uint64
	var err error
	switch {
	case strings.HasPrefix(bitstring, "0x"):
		clean := strings.ReplaceAll(bitstring[2:], "_", "")
		board, err = strconv.ParseUint(clean, 16, 64)
		if err != nil {
			fatal(err)
		}
	case strings.HasPrefix(bitstring, "0b"):
		clean := strings.ReplaceAll(bitstring[2:], "_", "")
		board, err = strconv.ParseUint(clean, 2, 64)
		if err != nil {
			fatal(err)
		}
	case strings.HasPrefix(bitstring, "0d"):
		board, err = strconv.ParseUint(bitstring[2:], 10, 64)
		if err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("invalid board: %s; should be a 64 bit decimal (0d...), hex (0x...), or binary (0b...) variable", bitstring))
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
