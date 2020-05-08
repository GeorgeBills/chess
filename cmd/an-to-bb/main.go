package main

import (
	"fmt"
	"os"
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

		var file uint8
		switch os.Args[i][0] {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			file = uint8(os.Args[i][0] - 'a')
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
			file = uint8(os.Args[i][0] - 'A')
		default:
			fatalsq(os.Args[i])
		}

		var rank uint8
		switch os.Args[i][1] {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			rank = uint8(os.Args[i][1] - '1')
		default:
			fatalsq(os.Args[i])
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
