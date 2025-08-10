package main

import (
	"fmt"
	"strings"
)

type RayMovesLookup struct {
	N  [64]Bitboard
	E  [64]Bitboard
	S  [64]Bitboard
	W  [64]Bitboard
	NE [64]Bitboard
	NW [64]Bitboard
	SE [64]Bitboard
	SW [64]Bitboard
}

var Rays RayMovesLookup

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

// firstHitPiece finds the nearest occupied board index along a ray.
//
// Parameters:
//   - ray: bitboard of all reachable indexes along a single direction from the origin index
//   - selfOccupancy: bitboard of own pieces
//   - enemyOccupancy: bitboard of opponent pieces
//   - increasing: when true, choose the lowest set bit on the ray (towards increasing file/rank);
//     when false, choose the highest set bit (towards decreasing file/rank)
//
// Returns:
// - nearestIndex: 0..63 board index of the nearest blocking piece on the ray
// - nearestIsEnemy: true if the blocking piece is an enemy piece, false if it is own piece
// - found: false when no piece lies on the ray; in that case nearestIndex/nearestIsEnemy are undefined
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

func GetRayMoves(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return Rays.N[index].Or(Rays.E[index]).Or(Rays.S[index]).Or(Rays.W[index]).
		Or(Rays.NE[index]).Or(Rays.NW[index]).Or(Rays.SE[index]).Or(Rays.SW[index])
}

func GetRayMovesFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetRayMoves(fileRankToIndex(file, rank))
}

func GetRayMovesFromSquare(square string) Bitboard {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return EmptyBitboard()
	}
	return GetRayMovesFromFileRank(file, rank)
}

func PrintRayMoves() {
	fmt.Println("Ray Moves (All directions):")
	fmt.Println("==============================")
	for index := uint64(0); index < 64; index++ {
		file, rank := indexToFileRank(index)
		square := fmt.Sprintf("%c%d", 'a'+file, rank+1)
		union := GetRayMoves(index)
		fmt.Printf("Square %s (%d):\n", square, index)
		fmt.Printf("  Bitboard:\n%s", union.String())
		fmt.Printf("  Squares: %s\n", strings.Join(union.ToSquares(), " "))
		fmt.Println()
	}
}

// GetValidRayMoves returns valid ray moves for a piece in each direction (N,E,S,W,NE,NW,SE,SW),
// useful for generating queen, rook, and bishop moves.
// Rays are truncated at the first blocker piece encountered. If the blocker is an enemy piece,
// that square is included in the valid moves. If it's a friendly piece, that square is excluded.
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
