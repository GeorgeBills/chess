package engine

import (
	"math"
)

const (
	maximizing = +1
	minimizing = -1
	infinity   = math.MaxInt16
)

func (g *Game) BestMoveToDepth(depth uint8, mm int8) (Move, int16) {
	if depth == 0 {
		return 0, g.Evaluate()
	}

	moves, _ := g.GenerateLegalMoves(nil)

	var best struct {
		score int16
		move  Move
	}

	switch mm {
	case maximizing:
		best.score = -1 * infinity
		for _, m := range moves {
			g.MakeMove(m)
			if _, s := g.BestMoveToDepth(depth-1, mm*-1); s > best.score {
				best.score = s
				best.move = m
			}
			g.UnmakeMove()
		}
	case minimizing:
		best.score = +1 * infinity
		for _, m := range moves {
			g.MakeMove(m)
			if _, s := g.BestMoveToDepth(depth-1, mm*-1); s < best.score {
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
