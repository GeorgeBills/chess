## Useful links

 * https://www.chessprogramming.org/
 * https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/
 * http://www.craftychess.com/hyatt/boardrep.html
 * https://lichess.org/editor/
 * https://godbolt.org/
 * https://play.golang.org/

## Useful commands

 * `go test ./engine -bench 'MoveGen'`

   Run unit tests and move generation micro-benchmarks.

 * `go test ./engine -covermode=count -coverprofile='coverage.out'`

   Generate unit test coverage info.
 
 * `go tool cover -html='coverage.out'`

   View coverage report.

 * `go test ./engine  -bench 'MoveGen' -cpuprofile cpu.prof -memprofile mem.prof`

   Write out profiling information.

 * `go tool pprof cpu.prof` (followed by e.g. `top`)

   View profiling information.

 * `go build -gcflags='-m' ./engine`

   Information on exactly which functions are being inlined and why.
