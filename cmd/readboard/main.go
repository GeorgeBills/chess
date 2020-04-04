package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Sprintf("%s <sq1> [sq2] [sq3] ... [sqn]", os.Args[0]))
	}

	var board uint64
	for i := 1; i < len(os.Args); i++ {
		if len(os.Args[i]) != 2 {
			fatal(fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		var file uint8
		switch os.Args[i][0] {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			file = uint8(os.Args[i][0]-'a') + 1
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
			file = uint8(os.Args[i][0]-'A') + 1
		default:
			fatal(fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		var rank uint8
		switch os.Args[i][1] {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			rank = uint8(os.Args[i][1] - '0')
		default:
			fatal(fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		board |= 1 << (8*(rank-1) + file - 1)
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

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}
