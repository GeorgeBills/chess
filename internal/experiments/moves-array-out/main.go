package main

import (
	"fmt"
	"math"
)

// theoretically ~218 max moves: https://www.chessprogramming.org/Chess_Position
// 255 isn't much bigger and hopefully triggers bounds check elimination for uint8 indexes
const maxMoves = math.MaxUint8

func main() {
	moves, n := generateMoves()
	fmt.Printf("%v", moves[0:n])
}

func generateMoves() ([math.MaxUint8]int, uint8) {
	var moves [math.MaxUint8]int
	var i uint8
	for i = 0; i < 20; i++ {
		moves[i] = 20
	}
	return moves, i
}
