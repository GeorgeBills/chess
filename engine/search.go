package engine

import (
	"fmt"
	"math"
	"time"
)

// https://www.chessprogramming.org/Search

const (
	maximizing = +1
	minimizing = -1
	infinity   = math.MaxInt16
)

func (g *Game) getMaximizingMinimizing() int8 {
	switch g.ToMove() {
	case White:
		return maximizing
	case Black:
		return minimizing
	default:
		panic(fmt.Errorf("invalid to move; %#v", g))
	}
}

type moveScore struct {
	move  Move
	score int16
}

type SearchInformation struct {
	depth              uint8
	time               time.Duration
	principalVariation []Move
	nodesPerSecond     uint64
}

func (g *Game) BestMoveInfinite(stopch <-chan struct{}, infoch chan<- SearchInformation) (Move, int16) {
	mm := g.getMaximizingMinimizing()
	var best moveScore
	var depth uint8
DEEPEN:
	for depth = 0; ; depth++ {
		select {
		case <-stopch:
			break DEEPEN
		default:
			// spiral out, keep going.
			best = g.bestMoveToDepth(depth, mm)
		}
	}
	return best.move, best.score
}

func (g *Game) BestMoveToTime(whiteTime, blackTime, whiteIncrement, blackIncrement time.Duration) (Move, int16) {
	return g.BestMoveToDepth(4)
	// TODO: properly implement basic time controls
	//
	// http://www.chessgames.com/chessstats.html
	// average game is 40 full moves (80 half moves) long
	//
	// asymptote from some max towards 0 for flat time controls
	// asymptote from some max towards increment if incremented
	//
	// expected average game length, based on number of moves played
	// games that go for 10 moves, on average go for another X moves...
	// games that go for 20 moves, on average go for another Y moves...
	// games that go for 30 moves, on average go for another Z moves...
	// https://chess.stackexchange.com/a/4899:
	//     59.3 + (72830 - 2330 k)/(2644 + k (10 + k))
	//
	// bump time if evaluations are unstable, the opposite if they're stable
}

// BestMoveToDepth returns the best move (with its score) to the given depth.
func (g *Game) BestMoveToDepth(depth uint8) (Move, int16) {
	mm := g.getMaximizingMinimizing()
	best := g.bestMoveToDepth(depth, mm)
	return best.move, best.score
}

func (g *Game) bestMoveToDepth(depth uint8, mm int8) moveScore {
	if depth == 0 {
		score := g.Evaluate()
		return moveScore{score: score}
	}

	moves, isCheck := g.GenerateLegalMoves(nil)

	var best moveScore

	switch mm {
	case maximizing:
		best.score = -1 * infinity
		if len(moves) == 0 && !isCheck {
			return moveScore{score: 0} // stalemate
		}
		for _, m := range moves {
			g.MakeMove(m)
			if child := g.bestMoveToDepth(depth-1, mm*-1); child.score >= best.score {
				best = moveScore{m, child.score}
			}
			g.UnmakeMove()
		}
	case minimizing:
		best.score = +1 * infinity
		if len(moves) == 0 && !isCheck {
			return moveScore{score: 0} // stalemate
		}
		for _, m := range moves {
			g.MakeMove(m)
			if child := g.bestMoveToDepth(depth-1, mm*-1); child.score <= best.score {
				best = moveScore{m, child.score}
			}
			g.UnmakeMove()
		}
	default:
		panic("mm neither minimizing nor maximizing")
	}

	return best
}
