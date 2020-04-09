package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Sprintf("%s <board>", os.Args[0]))
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
		fatal(fmt.Sprintf("invalid board: %s; should be a 64 bit hex (0x...) or binary (0b...) variable", bitstring))
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	for i := 0; i < 64; i++ {
		if i != 0 && i%8 == 0 {
			f.WriteRune('\n')
		}
		idx := i + 56 - 16*(i/8) // reverse ranks as we print
		if board&(1<<idx) != 0 {
			f.WriteRune('■')
		} else {
			f.WriteRune('□')
		}
	}
}

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}
