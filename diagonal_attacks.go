package main

import (
	"fmt"
	"strings"
)

var DiagonalAttacks [64]Bitboard

func generateDiagonalAttacks(square uint64) Bitboard {
	file, rank := indexToFileRank(square)
	attacks := EmptyBitboard()

	// Positive diagonal (top-left to bottom-right)
	for i := 1; i < 8; i++ {
		f, r := file+i, rank+i
		if f >= 8 || r >= 8 {
			break
		}
		attacks = attacks.Set(fileRankToIndex(f, r))
	}
	for i := 1; i < 8; i++ {
		f, r := file-i, rank-i
		if f < 0 || r < 0 {
			break
		}
		attacks = attacks.Set(fileRankToIndex(f, r))
	}

	// Negative diagonal (top-right to bottom-left)
	for i := 1; i < 8; i++ {
		f, r := file+i, rank-i
		if f >= 8 || r < 0 {
			break
		}
		attacks = attacks.Set(fileRankToIndex(f, r))
	}
	for i := 1; i < 8; i++ {
		f, r := file-i, rank+i
		if f < 0 || r >= 8 {
			break
		}
		attacks = attacks.Set(fileRankToIndex(f, r))
	}

	return attacks
}

func generateAllDiagonalAttacks() [64]Bitboard {
	var attacks [64]Bitboard
	for square := uint64(0); square < 64; square++ {
		attacks[square] = generateDiagonalAttacks(square)
	}
	return attacks
}

func init() {
	DiagonalAttacks = generateAllDiagonalAttacks()
}

func PrintDiagonalAttacks() {
	fmt.Println("Diagonal Attacks (Bishop moves):")
	fmt.Println("================================")

	for square := uint64(0); square < 64; square++ {
		file, rank := indexToFileRank(square)
		squareName := fmt.Sprintf("%c%d", 'a'+file, rank+1)

		fmt.Printf("Square %s (%d):\n", squareName, square)
		fmt.Printf("  Bitboard:\n%s", DiagonalAttacks[square].String())
		fmt.Printf("  Squares: %s\n", strings.Join(DiagonalAttacks[square].ToSquares(), " "))
		fmt.Println()
	}
}

func GetDiagonalAttacks(square uint64) Bitboard {
	if square >= 64 {
		return EmptyBitboard()
	}
	return DiagonalAttacks[square]
}

func GetDiagonalAttacksFromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	return GetDiagonalAttacks(fileRankToIndex(file, rank))
}

func GetDiagonalAttacksFromSquare(square string) Bitboard {
	if len(square) != 2 {
		return EmptyBitboard()
	}

	file := int(square[0] - 'a')
	rank := int(square[1] - '1')

	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}

	return GetDiagonalAttacksFromFileRank(file, rank)
}

func GetValidDiagonalMoves(pos *Position, piece *Piece) Bitboard {
	if piece == nil || (piece.Kind != Bishop && piece.Kind != Queen) {
		return EmptyBitboard()
	}

	// Get the piece's square index
	if piece.Location.IsEmpty() {
		return EmptyBitboard()
	}
	square := piece.Location.FirstSet()
	if square >= 64 {
		return EmptyBitboard()
	}

	// Start with all potential diagonal attacks from this square
	potentialMoves := DiagonalAttacks[square]
	validMoves := EmptyBitboard()

	// Get file and rank for blocking detection
	file, rank := piece.FileRank()
	if file < 0 || rank < 0 {
		return EmptyBitboard()
	}

	// Check each diagonal direction for blocking pieces
	// Positive diagonal (top-left to bottom-right)
	for i := 1; i < 8; i++ {
		f, r := file+i, rank+i
		if f >= 8 || r >= 8 {
			break
		}
		squareIndex := fileRankToIndex(f, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}
	for i := 1; i < 8; i++ {
		f, r := file-i, rank-i
		if f < 0 || r < 0 {
			break
		}
		squareIndex := fileRankToIndex(f, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}

	// Negative diagonal (top-right to bottom-left)
	for i := 1; i < 8; i++ {
		f, r := file+i, rank-i
		if f >= 8 || r < 0 {
			break
		}
		squareIndex := fileRankToIndex(f, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}
	for i := 1; i < 8; i++ {
		f, r := file-i, rank+i
		if f < 0 || r >= 8 {
			break
		}
		squareIndex := fileRankToIndex(f, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}

	return validMoves
}
