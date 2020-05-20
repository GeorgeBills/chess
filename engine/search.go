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
			if _, s := g.bestMoveToDepth(depth-1, mm*-1); s >= best.score {
				best.score = s
				best.move = m
			}
			g.UnmakeMove()
		}
	case minimizing:
		best.score = +1 * infinity
		if len(moves) == 0 && !isCheck {
			return 0, 0 // stalemate
		}
		for _, m := range moves {
			g.MakeMove(m)
			if _, s := g.bestMoveToDepth(depth-1, mm*-1); s <= best.score {
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
