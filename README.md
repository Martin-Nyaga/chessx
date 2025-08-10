## ChessX

Tiny, zero-dependency chess engine (in pure go) that can actually play a game.

For now Dumbfish always picks the first legal move it finds.

### Run

- Play a game in your terminal (it's not pretty):
  - `go run .`

### Tests

- `go test ./...`
- If you have `stockfish` installed and in your `PATH`, you can set `CHESSX_STOCKFISH=1` to compare this engine's generated list of legal moves against stockfish.
- Add `CHESSX_VERBOSE=1` to print debug positions along the way.
