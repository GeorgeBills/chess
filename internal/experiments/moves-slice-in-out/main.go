package main

import "fmt"

func main() {
	var moves []int
	moves = generateMoves(moves)
	fmt.Printf("%v", moves)
}

func generateMoves(moves []int) []int { // leaking param: moves to result ~r1 level=0
	for i := 0; i < 20; i++ {
		moves = append(moves, i)
	}
	return moves
}
