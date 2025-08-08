package main

// KnightAttacks stores precomputed attack bitboards for a knight from each square
var KnightAttacks [64]Bitboard

func init() {
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		bb := EmptyBitboard()

		deltas := [8][2]int{
			{1, 2}, {2, 1}, {2, -1}, {1, -2},
			{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
		}

		for _, d := range deltas {
			nf := file + d[0]
			nr := rank + d[1]
			if nf >= 0 && nf < 8 && nr >= 0 && nr < 8 {
				bb = bb.Set(fileRankToIndex(nf, nr))
			}
		}

		KnightAttacks[index] = bb
	}
}

func GetKnightAttacks(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return KnightAttacks[index]
}

func GetKnightAttacksFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetKnightAttacks(fileRankToIndex(file, rank))
}

func GetKnightAttacksFromSquare(square string) Bitboard {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return EmptyBitboard()
	}
	return GetKnightAttacksFromFileRank(file, rank)
}

// GetValidKnightMoves returns all squares a given knight piece can legally move to,
// excluding squares occupied by own pieces. Enemy-occupied squares are included.
func GetValidKnightMoves(pos *Position, piece *Piece) Bitboard {
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

	return KnightAttacks[index].And(selfOccupancy.Not())
}
