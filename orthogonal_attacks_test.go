package main

import (
	"testing"
)

func TestOrthogonalAttacks(t *testing.T) {
	tests := []struct {
		square      string
		expected    []string
		description string
	}{
		{
			square:      "a1",
			expected:    []string{"b1", "c1", "d1", "e1", "f1", "g1", "h1", "a2", "a3", "a4", "a5", "a6", "a7", "a8"},
			description: "Corner square should attack all squares in same rank and file",
		},
		{
			square:      "e4",
			expected:    []string{"a4", "b4", "c4", "d4", "f4", "g4", "h4", "e1", "e2", "e3", "e5", "e6", "e7", "e8"},
			description: "Center square should attack all squares in same rank and file",
		},
		{
			square:      "h8",
			expected:    []string{"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h1", "h2", "h3", "h4", "h5", "h6", "h7"},
			description: "Corner square should attack all squares in same rank and file",
		},
	}

	for _, test := range tests {
		t.Run(test.square, func(t *testing.T) {
			attacks := GetOrthogonalAttacksFromSquare(test.square)
			squares := attacks.ToSquares()

			if len(squares) != len(test.expected) {
				t.Errorf("Expected %d squares, got %d", len(test.expected), len(squares))
				return
			}

			for _, expected := range test.expected {
				found := false
				for _, square := range squares {
					if square == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected square %s not found in attacks", expected)
				}
			}
		})
	}
}

func TestOrthogonalAttacksCount(t *testing.T) {
	for square := uint64(0); square < 64; square++ {
		attacks := GetOrthogonalAttacks(square)
		expectedCount := 14

		if attacks.Count() != expectedCount {
			t.Errorf("Square %d: Expected %d attacks, got %d", square, expectedCount, attacks.Count())
		}
	}
}

func TestOrthogonalAttacksInvalidInput(t *testing.T) {
	invalidSquares := []string{"", "a", "a9", "i1", "z5", "a0"}

	for _, square := range invalidSquares {
		attacks := GetOrthogonalAttacksFromSquare(square)
		if !attacks.IsEmpty() {
			t.Errorf("Expected empty bitboard for invalid square '%s', got %d attacks", square, attacks.Count())
		}
	}
}

func TestGetValidOrthogonalMoves(t *testing.T) {
	pos := NewPosition()

	pos.SetPieceAtSquare("a1", Rook, White)
	pos.SetPieceAtSquare("b1", Pawn, Black)
	pos.SetPieceAtSquare("a2", Pawn, White)

	rook := pos.GetPieceAtSquare("a1")
	moves := GetValidOrthogonalMoves(pos, rook)
	squares := moves.ToSquares()

	expected := []string{"b1"}
	for _, square := range expected {
		found := false
		for _, s := range squares {
			if s == square {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected square %s not found in moves", square)
		}
	}

	if len(squares) != len(expected) {
		t.Errorf("Expected %d moves, got %d", len(expected), len(squares))
	}
}
