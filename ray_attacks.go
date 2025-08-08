package main

import (
	"fmt"
	"strings"
)

type RayAttacks struct {
	N  [64]Bitboard
	E  [64]Bitboard
	S  [64]Bitboard
	W  [64]Bitboard
	NE [64]Bitboard
	NW [64]Bitboard
	SE [64]Bitboard
	SW [64]Bitboard
}

var Rays RayAttacks

func init() {
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		ray := EmptyBitboard()
		for f := file + 1; f < 8; f++ {
			ray = ray.Set(fileRankToIndex(f, rank))
		}
		Rays.E[index] = ray
		ray = EmptyBitboard()
		for f := file - 1; f >= 0; f-- {
			ray = ray.Set(fileRankToIndex(f, rank))
		}
		Rays.W[index] = ray
		ray = EmptyBitboard()
		for r := rank + 1; r < 8; r++ {
			ray = ray.Set(fileRankToIndex(file, r))
		}
		Rays.N[index] = ray
		ray = EmptyBitboard()
		for r := rank - 1; r >= 0; r-- {
			ray = ray.Set(fileRankToIndex(file, r))
		}
		Rays.S[index] = ray
		ray = EmptyBitboard()
		for f, r := file+1, rank+1; f < 8 && r < 8; f, r = f+1, r+1 {
			ray = ray.Set(fileRankToIndex(f, r))
		}
		Rays.NE[index] = ray
		ray = EmptyBitboard()
		for f, r := file-1, rank+1; f >= 0 && r < 8; f, r = f-1, r+1 {
			ray = ray.Set(fileRankToIndex(f, r))
		}
		Rays.NW[index] = ray
		ray = EmptyBitboard()
		for f, r := file+1, rank-1; f < 8 && r >= 0; f, r = f+1, r-1 {
			ray = ray.Set(fileRankToIndex(f, r))
		}
		Rays.SE[index] = ray
		ray = EmptyBitboard()
		for f, r := file-1, rank-1; f >= 0 && r >= 0; f, r = f-1, r-1 {
			ray = ray.Set(fileRankToIndex(f, r))
		}
		Rays.SW[index] = ray
	}
}

// firstHitPiece returns the nearest occupied index on the ray and whether it is an enemy.
func firstHitPiece(ray Bitboard, selfOccupancy, enemyOccupancy Bitboard, increasing bool) (uint64, bool, bool) {
	selfOccupancyMask := ray.And(selfOccupancy)
	enemyOccupancyMask := ray.And(enemyOccupancy)

	const noIndex = ^uint64(0)
	nearestIndex := noIndex
	nearestIsEnemy := false

	if !selfOccupancyMask.IsEmpty() {
		if increasing {
			nearestIndex = selfOccupancyMask.FirstSet()
		} else {
			nearestIndex = selfOccupancyMask.LastSet()
		}
		nearestIsEnemy = false
	}

	if !enemyOccupancyMask.IsEmpty() {
		var candidateIndex uint64
		if increasing {
			candidateIndex = enemyOccupancyMask.FirstSet()
		} else {
			candidateIndex = enemyOccupancyMask.LastSet()
		}
		if nearestIndex == noIndex {
			nearestIndex = candidateIndex
			nearestIsEnemy = true
		} else {
			if increasing {
				if candidateIndex < nearestIndex {
					nearestIndex = candidateIndex
					nearestIsEnemy = true
				}
			} else {
				if candidateIndex > nearestIndex {
					nearestIndex = candidateIndex
					nearestIsEnemy = true
				}
			}
		}
	}

	if nearestIndex == noIndex {
		return 0, false, false
	}
	return nearestIndex, nearestIsEnemy, true
}

type RayMoves struct {
	N  Bitboard
	E  Bitboard
	S  Bitboard
	W  Bitboard
	NE Bitboard
	NW Bitboard
	SE Bitboard
	SW Bitboard
}

func (moves RayMoves) Union() Bitboard {
	return moves.N.Or(moves.E).Or(moves.S).Or(moves.W).Or(moves.NE).Or(moves.NW).Or(moves.SE).Or(moves.SW)
}

func (moves RayMoves) Orthogonal() Bitboard {
	return moves.N.Or(moves.E).Or(moves.S).Or(moves.W)
}

func (moves RayMoves) Diagonal() Bitboard {
	return moves.NE.Or(moves.NW).Or(moves.SE).Or(moves.SW)
}

func (moves RayMoves) All() Bitboard {
	return moves.Union()
}

func GetRayAttacks(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return Rays.N[index].Or(Rays.E[index]).Or(Rays.S[index]).Or(Rays.W[index]).
		Or(Rays.NE[index]).Or(Rays.NW[index]).Or(Rays.SE[index]).Or(Rays.SW[index])
}

func GetRayAttacksFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetRayAttacks(fileRankToIndex(file, rank))
}

func GetRayAttacksFromSquare(square string) Bitboard {
	if len(square) != 2 {
		return EmptyBitboard()
	}
	file := int(square[0] - 'a')
	rank := int(square[1] - '1')
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetRayAttacksFromFileRank(file, rank)
}

func PrintRayAttacks() {
	fmt.Println("Ray Attacks (All directions):")
	fmt.Println("==============================")
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		square := fmt.Sprintf("%c%d", 'a'+file, rank+1)
		union := GetRayAttacks(index)
		fmt.Printf("Square %s (%d):\n", square, index)
		fmt.Printf("  Bitboard:\n%s", union.String())
		fmt.Printf("  Squares: %s\n", strings.Join(union.ToSquares(), " "))
		fmt.Println()
	}
}

func GetValidRayMoves(pos *Position, piece *Piece) RayMoves {
	if piece == nil || piece.Location.IsEmpty() {
		return RayMoves{}
	}
	index := piece.Location.FirstSet()
	if index >= 64 {
		return RayMoves{}
	}

	var selfOccupancy, enemyOccupancy Bitboard
	if piece.Color == White {
		selfOccupancy = pos.GetWhiteOccupancy()
		enemyOccupancy = pos.GetBlackOccupancy()
	} else {
		selfOccupancy = pos.GetBlackOccupancy()
		enemyOccupancy = pos.GetWhiteOccupancy()
	}

	// Truncate a direction's ray at the nearest blocker using XOR with the ray from the blocker.
	computeDirection := func(directionTable [64]Bitboard, increasing bool) Bitboard {
		ray := directionTable[index]
		if ray.IsEmpty() {
			return EmptyBitboard()
		}
		nearestIndex, nearestIsEnemy, found := firstHitPiece(ray, selfOccupancy, enemyOccupancy, increasing)
		if !found {
			return ray
		}
		between := ray.Xor(directionTable[nearestIndex])
		if nearestIsEnemy {
			return between
		}
		return between.And(FromIndex(nearestIndex).Not())
	}

	return RayMoves{
		N:  computeDirection(Rays.N, true),
		E:  computeDirection(Rays.E, true),
		S:  computeDirection(Rays.S, false),
		W:  computeDirection(Rays.W, false),
		NE: computeDirection(Rays.NE, true),
		NW: computeDirection(Rays.NW, true),
		SE: computeDirection(Rays.SE, false),
		SW: computeDirection(Rays.SW, false),
	}
}
