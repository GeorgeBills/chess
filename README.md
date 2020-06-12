## Playing vs the engine

The engine is playable via [UCI][UCI]. It's currently very basic, but easily
beats me (that's a very low bar).

[UCI]: <https://www.chessprogramming.org/UCI> "Universal Chess Interface"

 1. Build the UCI wrapped engine with `go build -o 'uci.exe' ./cmd/uci/`.
 2. Download a UCI compatible GUI e.g. [Arena](http://www.playwitharena.de/).
 3. Follow the GUI instructions to install `uci.exe` as an engine.

## Useful links

 * https://www.chessprogramming.org/
 * https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/
 * http://www.craftychess.com/hyatt/boardrep.html
 * https://lichess.org/editor/
 * https://godbolt.org/
 * https://play.golang.org/
 * https://about.sourcegraph.com/go/gophercon-2019-optimizing-go-code-without-a-blindfold
 * https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/

## Useful commands

 * `go test ./engine -bench 'GenerateLegalMoves'`

   Run unit tests and move generation micro-benchmarks.

 * `go test ./engine -covermode=count -coverprofile='coverage.out'`

   Generate unit test coverage info.
 
 * `go tool cover -html='coverage.out'`

   View coverage report.

 * `go test ./engine  -bench 'GenerateLegalMoves' -cpuprofile cpu.prof -memprofile mem.prof`

   Write out profiling information.

 * `go tool pprof cpu.prof` (followed by e.g. `top` or `list GenerateLegalMoves`)

   View profiling information.

 * `go build -gcflags='-m' ./engine`

   Information on exactly which functions are being inlined and why.
