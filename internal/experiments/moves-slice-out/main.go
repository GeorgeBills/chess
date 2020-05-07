package main

import "fmt"

func main() {
	moves := generateMoves()
	fmt.Printf("%v", moves)
}

func generateMoves() []int {
	var moves []int
	for i := 0; i < 20; i++ {
		moves = append(moves, i)
	}
	return moves
}
