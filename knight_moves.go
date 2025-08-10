package main

// KnightMoves stores precomputed move bitboards for a knight from each square
var KnightMoves [64]Bitboard

func init() {
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		attacksMask := EmptyBitboard()

		moveDeltas := [8][2]int{
			{1, 2}, {2, 1}, {2, -1}, {1, -2},
			{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
		}

		for _, delta := range moveDeltas {
			nextFile := file + delta[0]
			nextRank := rank + delta[1]
			if nextFile >= 0 && nextFile < 8 && nextRank >= 0 && nextRank < 8 {
				attacksMask = attacksMask.Set(fileRankToIndex(nextFile, nextRank))
			}
		}

		KnightMoves[index] = attacksMask
	}
}

func GetKnightMoves(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return KnightMoves[index]
}

func GetKnightMovesFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetKnightMoves(fileRankToIndex(file, rank))
}

func GetKnightMovesFromSquare(square string) Bitboard {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return EmptyBitboard()
	}
	return GetKnightMovesFromFileRank(file, rank)
}

// GetValidKnightMoves returns all squares a given knight piece can legally move to,
// excluding squares occupied by own pieces. Enemy-occupied squares are included.
func GetPossibleKnightMoves(pos *Position, piece *Piece) Bitboard {
	if piece == nil || piece.Location.IsEmpty() {
		return EmptyBitboard()
	}
	index := piece.Location.FirstSet()
	if index >= 64 {
		return EmptyBitboard()
	}

	var selfOccupancy Bitboard
	if piece.Color == White {
		selfOccupancy = pos.GetWhiteOccupancy()
	} else {
		selfOccupancy = pos.GetBlackOccupancy()
	}

	return KnightMoves[index].And(selfOccupancy.Not())
}
