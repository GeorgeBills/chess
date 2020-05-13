package engine

import (
	"fmt"
	"math"
)

// https://www.chessprogramming.org/Search

const (
	maximizing = +1
	minimizing = -1
	infinity   = math.MaxInt16
)

// BestMoveToDepth returns the best move (with its score) to the given depth.
func (g *Game) BestMoveToDepth(depth uint8) (Move, int16) {
	switch g.ToMove() {
	case White:
		return g.bestMoveToDepth(depth, maximizing)
	case Black:
		return g.bestMoveToDepth(depth, minimizing)
	default:
		panic(fmt.Errorf("invalid to move; %#v", g))
	}
}

func (g *Game) bestMoveToDepth(depth uint8, mm int8) (Move, int16) {
	if depth == 0 {
		return 0, g.Evaluate()
	}

	moves, isCheck := g.GenerateLegalMoves(nil)

	var best struct {
		score int16
		move  Move
	}

	switch mm {
	case maximizing:
		best.score = -1 * infinity
		if len(moves) == 0 && !isCheck {
			return 0, 0 // stalemate
		}
		for _, m := range moves {
			g.MakeMove(m)
			if _, s := g.bestMoveToDepth(depth-1, mm*-1); s > best.score {
				best.score = s
				best.move = m
			}
			g.UnmakeMove()
		}
	case minimizing:
		best.score = +1 * infinity
		for _, m := range moves {
			if len(moves) == 0 && !isCheck {
				return 0, 0 // stalemate
			}
			g.MakeMove(m)
			if _, s := g.bestMoveToDepth(depth-1, mm*-1); s < best.score {
				best.score = s
				best.move = m
			}
			g.UnmakeMove()
		}
	default:
		panic("mm neither minimizing nor maximizing")
	}

	return best.move, best.score
}
