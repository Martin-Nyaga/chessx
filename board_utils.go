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
