package main_test

import (
	"strconv"
	"testing"
)

var stopch chan struct{} = make(chan struct{})

// run a busy loop just for a baseline
func BenchmarkBusyLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

// poll the channel on every iteration
func BenchmarkPollChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		select {
		case <-stopch:
			return
		default:
		}
	}
}

// poll the channel every n iterations, where n is a variable
// this is half the speed of just polling the channel on every iteration!
// mod overhead, register loading, some compiler optimisation we can't do?
// https://godbolt.org/z/nfniCy
func BenchmarkPollChannelEveryN(b *testing.B) {
	tests := []int{100, 1_000, 10_000}
	for _, n := range tests {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if i%n == 0 {
					select {
					case <-stopch:
						return
					default:
					}
				}
			}
		})
	}
}

// hardcode n: this is about 4 times faster than polling on every loop
// it doesn't scale linearly with n; mod overhead must be dominating
// https://godbolt.org/z/3pGZah
func BenchmarkPollChannelEvery1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			select {
			case <-stopch:
				return
			default:
			}
		}
	}
}

// poll the channel every n = 1024 operations: twice as fast as n = 1000
// compiler recognizes that we can just do a bitwise and here: testq, jne
// about half the speed of the busy loop (no polling) version
// https://godbolt.org/z/LkDafK
func BenchmarkPollChannelEvery1024(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i%1024 == 0 {
			select {
			case <-stopch:
				return
			default:
			}
		}
	}
}

// trying to manually replicate the n = 1024 compiler optimisation
// this is much slower! no-idea-what-im-doing.jpg
func BenchmarkPollChannelBitwiseAnd1024(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if i&0b10000000000 != 0 {
			select {
			case <-stopch:
				return
			default:
			}
		}
	}
}
