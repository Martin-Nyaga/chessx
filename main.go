package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Chess Engine Demo")
	fmt.Println("=================")

	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(startingFEN)
	if err != nil {
		fmt.Printf("Error parsing FEN: %v\n", err)
		return
	}

	fmt.Println("Starting Position:")
	fmt.Println(pos.String())
	fmt.Println()

	fmt.Println("Orthogonal Attack Patterns (Rook moves):")
	fmt.Println("========================================")

	exampleSquares := []string{"a1", "e4", "h8"}
	for _, square := range exampleSquares {
		attacks := GetOrthogonalAttacksFromSquare(square)
		fmt.Printf("Rook at %s can attack:\n", square)
		fmt.Printf("  Squares: %s\n", strings.Join(attacks.ToSquares(), " "))
		fmt.Printf("  Bitboard:\n%s", attacks.String())
		fmt.Println()
	}

	fmt.Println("Valid Orthogonal Moves (considering blocking pieces):")
	fmt.Println("=====================================================")

	testPos := NewPosition()
	testPos.SetPiece(0, 0, Rook, White)
	testPos.SetPiece(1, 0, Pawn, Black)
	testPos.SetPiece(0, 1, Pawn, White)

	fmt.Println("Position with White Rook at a1, Black Pawn at b1, White Pawn at a2:")
	fmt.Println(testPos.String())

	moves := GetValidOrthogonalMoves(testPos, 0, 0)
	fmt.Printf("Valid moves for Rook at a1: %s\n", strings.Join(moves.ToSquares(), " "))
	fmt.Printf("  Bitboard:\n%s", moves.String())
	fmt.Println()

	fmt.Println("Diagonal Attack Patterns (Bishop moves):")
	fmt.Println("=======================================")

	exampleSquares = []string{"a1", "e4", "h8"}
	for _, square := range exampleSquares {
		attacks := GetDiagonalAttacksFromSquare(square)
		fmt.Printf("Bishop at %s can attack:\n", square)
		fmt.Printf("  Squares: %s\n", strings.Join(attacks.ToSquares(), " "))
		fmt.Printf("  Bitboard:\n%s", attacks.String())
		fmt.Println()
	}

	fmt.Println("Valid Diagonal Moves (considering blocking pieces):")
	fmt.Println("==================================================")

	testPos2 := NewPosition()
	testPos2.SetPiece(0, 0, Bishop, White)
	testPos2.SetPiece(1, 1, Pawn, Black)
	testPos2.SetPiece(2, 2, Pawn, White)

	fmt.Println("Position with White Bishop at a1, Black Pawn at b2, White Pawn at c3:")
	fmt.Println(testPos2.String())

	moves2 := GetValidDiagonalMoves(testPos2, 0, 0)
	fmt.Printf("Valid moves for Bishop at a1: %s\n", strings.Join(moves2.ToSquares(), " "))
	fmt.Printf("  Bitboard:\n%s", moves2.String())
}
