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

func GetValidOrthogonalMoves(pos *Position, piece *Piece) Bitboard {
	if piece == nil || (piece.Kind != Rook && piece.Kind != Queen) {
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

	// Start with all potential orthogonal attacks from this square
	potentialMoves := OrthogonalAttacks[square]
	validMoves := EmptyBitboard()

	// Get file and rank for blocking detection
	file, rank := piece.FileRank()
	if file < 0 || rank < 0 {
		return EmptyBitboard()
	}

	// Check each direction for blocking pieces
	// Left direction
	for f := file - 1; f >= 0; f-- {
		squareIndex := fileRankToIndex(f, rank)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, rank)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}

	// Right direction
	for f := file + 1; f < 8; f++ {
		squareIndex := fileRankToIndex(f, rank)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(f, rank)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}

	// Down direction
	for r := rank - 1; r >= 0; r-- {
		squareIndex := fileRankToIndex(file, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(file, r)
		if target != nil {
			if target.Color != piece.Color {
				validMoves = validMoves.Set(squareIndex)
			}
			break
		}
		validMoves = validMoves.Set(squareIndex)
	}

	// Up direction
	for r := rank + 1; r < 8; r++ {
		squareIndex := fileRankToIndex(file, r)
		if !potentialMoves.IsSet(squareIndex) {
			break
		}
		target := pos.GetPiece(file, r)
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
