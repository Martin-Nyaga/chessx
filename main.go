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
		attacks := GetRayAttacksFromSquare(square).And(FromSquare("a1").Not()) // placeholder mask removal unused; will filter below
		// Filter to orthogonal rays only
		file := int(square[0] - 'a')
		rank := int(square[1] - '1')
		index := fileRankToIndex(file, rank)
		attacks = Rays.N[index].Or(Rays.E[index]).Or(Rays.S[index]).Or(Rays.W[index])
		fmt.Printf("Rook at %s can attack:\n", square)
		fmt.Printf("  Squares: %s\n", strings.Join(attacks.ToSquares(), " "))
		fmt.Printf("  Bitboard:\n%s", attacks.String())
		fmt.Println()
	}

	fmt.Println("Valid Orthogonal Moves (considering blocking pieces):")
	fmt.Println("=====================================================")

	testPos := NewPosition()
	testPos.SetPieceAtSquare("a1", Rook, White)
	testPos.SetPieceAtSquare("b1", Pawn, Black)
	testPos.SetPieceAtSquare("a2", Pawn, White)

	fmt.Println("Position with White Rook at a1, Black Pawn at b1, White Pawn at a2:")
	fmt.Println(testPos.String())

	rook := testPos.GetPieceAtSquare("a1")
	moves := GetValidRayMoves(testPos, rook)
    orthogonal := moves.Orthogonal()
	fmt.Printf("Valid moves for Rook at a1: %s\n", strings.Join(orthogonal.ToSquares(), " "))
	fmt.Printf("  Bitboard:\n%s", orthogonal.String())
	fmt.Println()

	fmt.Println("Diagonal Attack Patterns (Bishop moves):")
	fmt.Println("=======================================")

	exampleSquares = []string{"a1", "e4", "h8"}
	for _, square := range exampleSquares {
		file := int(square[0] - 'a')
		rank := int(square[1] - '1')
		index := fileRankToIndex(file, rank)
		attacks := Rays.NE[index].Or(Rays.NW[index]).Or(Rays.SE[index]).Or(Rays.SW[index])
		fmt.Printf("Bishop at %s can attack:\n", square)
		fmt.Printf("  Squares: %s\n", strings.Join(attacks.ToSquares(), " "))
		fmt.Printf("  Bitboard:\n%s", attacks.String())
		fmt.Println()
	}

	fmt.Println("Valid Diagonal Moves (considering blocking pieces):")
	fmt.Println("==================================================")

	testPos2 := NewPosition()
	testPos2.SetPieceAtSquare("a1", Bishop, White)
	testPos2.SetPieceAtSquare("b2", Pawn, Black)
	testPos2.SetPieceAtSquare("c3", Pawn, White)

	fmt.Println("Position with White Bishop at a1, Black Pawn at b2, White Pawn at c3:")
	fmt.Println(testPos2.String())

	bishop := testPos2.GetPieceAtSquare("a1")
	moves2 := GetValidRayMoves(testPos2, bishop)
    diagonal := moves2.Diagonal()
	fmt.Printf("Valid moves for Bishop at a1: %s\n", strings.Join(diagonal.ToSquares(), " "))
	fmt.Printf("  Bitboard:\n%s", diagonal.String())
}
