// Package uci implements lexing and parsing for the Universal Chess Interface.
//
// The UCI protocol is canonically documented at
// https://www.shredderchess.com/chess-info/features/uci-universal-chess-interface.html.
//
// A summarised list of UCI commands is documented below:
//
//     uci
//     debug <on|off>
//     isready
//     setoption name <id> [value <x>]
//     register [later|name <x>|code <y>]
//     ucinewgame
//     position [fen <fenstring> | startpos ]  moves <move1> ... <movei>
//     go
//        searchmoves <move1> ... <movei>
//        ponder
//        wtime <x>
//        btime <x>
//        winc <x>
//        binc <x>
//        movestogo <x>
//        depth <x>
//        nodes <x>
//        mate <x>
//        movetime <x>
//        infinite
//     stop
//     ponderhit
//     quit
package uci
