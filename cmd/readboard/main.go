package main

import (
	"fmt"
	"os"
)

func main() {
	fatal := func(code int, v interface{}) {
		fmt.Println(v)
		os.Exit(code)
	}

	if len(os.Args) < 2 {
		fatal(0, fmt.Sprintf("%s <sq1> [sq2] [sq3] ... [sqn]", os.Args[0]))
	}

	var board uint64
	for i := 1; i < len(os.Args); i++ {
		if len(os.Args[i]) != 2 {
			fatal(0, fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		var file uint8
		switch os.Args[i][0] {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			file = uint8(os.Args[i][0]-'a') + 1
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
			file = uint8(os.Args[i][0]-'A') + 1
		default:
			fatal(0, fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		var rank uint8
		switch os.Args[i][1] {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			rank = uint8(os.Args[i][1] - '0')
		default:
			fatal(0, fmt.Sprintf("invalid square: %s", os.Args[i]))
		}

		board |= 1 << (8*(rank-1) + file - 1)
	}
	fmt.Printf("0b%064b", board)
}
