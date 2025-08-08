package main

// fileRankToIndex converts 0..7 file/rank to 0..63 board index
func fileRankToIndex(file, rank int) uint64 {
	return uint64(rank*8 + file)
}

// indexToFileRank converts 0..63 board index to 0..7 file/rank
func indexToFileRank(index uint64) (int, int) {
	file := int(index % 8)
	rank := int(index / 8)
	return file, rank
}

// squareToFileRank converts an algebraic square (e.g., "e4") to file/rank (0..7).
// Returns ok=false if input is invalid.
func squareToFileRank(square string) (file int, rank int, ok bool) {
	if len(square) != 2 {
		return 0, 0, false
	}
	file = int(square[0] - 'a')
	rank = int(square[1] - '1')
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return 0, 0, false
	}
	return file, rank, true
}

// squareToIndex converts an algebraic square (e.g., "e4") to a 0..63 index.
// Returns ok=false if input is invalid.
func squareToIndex(square string) (idx uint64, ok bool) {
	file, rank, valid := squareToFileRank(square)
	if !valid {
		return 0, false
	}
	return fileRankToIndex(file, rank), true
}
