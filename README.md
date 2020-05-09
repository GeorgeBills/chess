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
