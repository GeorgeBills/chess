package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fatal := func(code int, v interface{}) {
		fmt.Println(v)
		os.Exit(code)
	}

	if len(os.Args) != 2 {
		fatal(0, fmt.Sprintf("%s <board>", os.Args[0]))
	}

	bitstring := os.Args[1]

	var board uint64
	var err error
	switch {
	case strings.HasPrefix(bitstring, "0x"):
		clean := strings.ReplaceAll(bitstring[2:], "_", "")
		board, err = strconv.ParseUint(clean, 16, 64)
		if err != nil {
			fatal(1, err)
		}
	case strings.HasPrefix(bitstring, "0b"):
		clean := strings.ReplaceAll(bitstring[2:], "_", "")
		board, err = strconv.ParseUint(clean, 2, 64)
		if err != nil {
			fatal(1, err)
		}
	default:
		fatal(0, fmt.Sprintf("invalid board: %s", bitstring))
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	for i := 0; i < 64; i++ {
		if i != 0 && i%8 == 0 {
			f.WriteRune('\n')
		}
		idx := i + 56 - 16*(i/8)
		if board&(1<<idx) != 0 {
			f.WriteRune('■')
		} else {
			f.WriteRune('□')
		}
	}
}
