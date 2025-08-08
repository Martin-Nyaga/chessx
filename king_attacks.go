package main

var KingAttacks [64]Bitboard

func init() {
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		attacksMask := EmptyBitboard()

		moveDeltas := [8][2]int{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1},
			{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		}
		for _, delta := range moveDeltas {
			nextFile := file + delta[0]
			nextRank := rank + delta[1]
			if nextFile >= 0 && nextFile < 8 && nextRank >= 0 && nextRank < 8 {
				attacksMask = attacksMask.Set(fileRankToIndex(nextFile, nextRank))
			}
		}
		KingAttacks[index] = attacksMask
	}
}

func GetKingAttacks(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return KingAttacks[index]
}

func GetKingAttacksFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetKingAttacks(fileRankToIndex(file, rank))
}

func GetKingAttacksFromSquare(square string) Bitboard {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return EmptyBitboard()
	}
	return GetKingAttacksFromFileRank(file, rank)
}

func GetValidKingMoves(pos *Position, piece *Piece) Bitboard {
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
	return KingAttacks[index].And(selfOccupancy.Not())
}
