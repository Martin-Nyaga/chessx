## ChessX

Tiny chess engine that can actually play a game. It knows the rules (yes, even castling and promotions), and it’ll happily face you with its very serious opponent: Dumbfish — a proud engine that always picks the first legal move it sees.

### Run
- Play a quick game in your terminal:
  - `go run .`

### Tests (optional fun)
- `go test ./...`
- Have Stockfish installed and set `CHESSX_STOCKFISH=1` to compare our moves against the big fish. Add `CHESSX_VERBOSE=1` if you like seeing lots of move lists.

### Why
Because building chess engines is fun, and sometimes “it works!” is all the victory you need.


