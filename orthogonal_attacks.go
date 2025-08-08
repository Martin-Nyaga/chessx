package main

import (
	"fmt"
	"strings"
)

var OrthogonalAttacks [64]Bitboard

func generateOrthogonalAttacks(square uint64) Bitboard {
	file, rank := indexToFileRank(square)
	attacks := EmptyBitboard()

	for f := 0; f < 8; f++ {
		if f != file {
			attacks = attacks.Set(fileRankToIndex(f, rank))
		}
	}

	for r := 0; r < 8; r++ {
		if r != rank {
			attacks = attacks.Set(fileRankToIndex(file, r))
		}
	}

	return attacks
}

func generateAllOrthogonalAttacks() [64]Bitboard {
	var attacks [64]Bitboard
	for square := uint64(0); square < 64; square++ {
		attacks[square] = generateOrthogonalAttacks(square)
	}
	return attacks
}

func init() {
	OrthogonalAttacks = generateAllOrthogonalAttacks()
}

func PrintOrthogonalAttacks() {
	fmt.Println("Orthogonal Attacks (Rook moves):")
	fmt.Println("=================================")

	for square := uint64(0); square < 64; square++ {
		file, rank := indexToFileRank(square)
		squareName := fmt.Sprintf("%c%d", 'a'+file, rank+1)

		fmt.Printf("Square %s (%d):\n", squareName, square)
		fmt.Printf("  Bitboard:\n%s", OrthogonalAttacks[square].String())
		fmt.Printf("  Squares: %s\n", strings.Join(OrthogonalAttacks[square].ToSquares(), " "))
		fmt.Println()
	}
}

func GetOrthogonalAttacks(square uint64) Bitboard {
	if square >= 64 {
		return EmptyBitboard()
	}
	return OrthogonalAttacks[square]
}

func GetOrthogonalAttacksFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetOrthogonalAttacks(fileRankToIndex(file, rank))
}

func GetOrthogonalAttacksFromSquare(square string) Bitboard {
	if len(square) != 2 {
		return EmptyBitboard()
	}

	file := int(square[0] - 'a')
	rank := int(square[1] - '1')

	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}

	return GetOrthogonalAttacksFromFileRank(file, rank)
}

func GetValidOrthogonalMoves(pos *Position, file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}

	piece := pos.GetPiece(file, rank)
	if piece == nil || (piece.Kind != Rook && piece.Kind != Queen) {
		return EmptyBitboard()
	}

	validMoves := EmptyBitboard()

	for f := file - 1; f >= 0; f-- {
		target := pos.GetPiece(f, rank)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(fileRankToIndex(f, rank))
			}
			break
		}
		validMoves = validMoves.Set(fileRankToIndex(f, rank))
	}

	for f := file + 1; f < 8; f++ {
		target := pos.GetPiece(f, rank)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(fileRankToIndex(f, rank))
			}
			break
		}
		validMoves = validMoves.Set(fileRankToIndex(f, rank))
	}

	for r := rank - 1; r >= 0; r-- {
		target := pos.GetPiece(file, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(fileRankToIndex(file, r))
			}
			break
		}
		validMoves = validMoves.Set(fileRankToIndex(file, r))
	}

	for r := rank + 1; r < 8; r++ {
		target := pos.GetPiece(file, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(fileRankToIndex(file, r))
			}
			break
		}
		validMoves = validMoves.Set(fileRankToIndex(file, r))
	}

	return validMoves
}
