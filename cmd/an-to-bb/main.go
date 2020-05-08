package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GeorgeBills/chess/m/v2/engine"
)

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("%s <sq1> [sq2] [sq3] ... [sqn]", os.Args[0]))
	}

	var board uint64
	for i := 1; i < len(os.Args); i++ {
		if len(os.Args[i]) != 2 {
			fatalsq(os.Args[i])
		}

		rank, file, err := engine.ParseAlgebraicNotation(strings.NewReader(os.Args[i]))
		if err != nil {
			fatal(fmt.Errorf("error parsing '%s' as algebraic notation: %w", os.Args[i], err))
		}

		board |= 1 << (8*rank + file)
	}

	bitstr := fmt.Sprintf("%064b", board)
	fmt.Printf(
		"0b%s_%s_%s_%s_%s_%s_%s_%s",
		bitstr[0:8],
		bitstr[8:16],
		bitstr[16:24],
		bitstr[24:32],
		bitstr[32:40],
		bitstr[40:48],
		bitstr[48:56],
		bitstr[56:64],
	)
}

func fatal(v error) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func fatalsq(sq string) {
	fatal(fmt.Errorf("invalid square: %s; must match ^[a-hA-H][1-8]$", sq))
}
