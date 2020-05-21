package main

import (
	"log"
	"time"
)

func main() {
	stopch := make(chan struct{})
	resultch := make(chan searchResult, 100)
	go bestMoveInfinite(stopch, resultch)

	go func() {
		for {
			result := <-resultch
			log.Println("result:", result)
		}
	}()

	time.Sleep(5 * time.Second)
	stopch <- struct{}{}
}

type move interface {
	From() uint8
	To() uint8
}

type searchResult struct {
	move  move
	score int16
}

func bestMoveInfinite(stopch <-chan struct{}, resultch chan<- searchResult) {
	for {
		select {
		case <-stopch:
			break
		default:
			// do some work (iterative deepening?) to calculate a search result.
			// the problem with iterative deepening is that we can't cancel out
			// of search 20 plies in until it's done, and a search that deep is
			// going to take a very long time. do we need to pass stopch down
			// into minimax and check it every nodes%x?
			time.Sleep(1 * time.Second)
			result := searchResult{}

			// return the result
			select {
			case resultch <- result:
			default:
				// we could write to an errorch, but what do we do if that's
				// full too? otoh it's likely we'll need to return other errs.
				// do we just drop the result on the floor, and document that
				// callers must drain resultch faster than we fill it?
				panic("resultch full")
			}
		}
	}
}
