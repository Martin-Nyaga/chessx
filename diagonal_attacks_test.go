package main

import (
	"testing"
)

func TestDiagonalAttacks(t *testing.T) {
	tests := []struct {
		square      string
		expected    []string
		description string
	}{
		{
			square:      "a1",
			expected:    []string{"b2", "c3", "d4", "e5", "f6", "g7", "h8"},
			description: "Corner square should attack all squares in positive diagonal",
		},
		{
			square:      "e4",
			expected:    []string{"a8", "b7", "c6", "d5", "f3", "g2", "h1", "b1", "c2", "d3", "f5", "g6", "h7"},
			description: "Center square should attack all squares in both diagonals",
		},
		{
			square:      "h8",
			expected:    []string{"a1", "b2", "c3", "d4", "e5", "f6", "g7"},
			description: "Corner square should attack all squares in negative diagonal",
		},
	}

	for _, test := range tests {
		t.Run(test.square, func(t *testing.T) {
			attacks := GetDiagonalAttacksFromSquare(test.square)
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

func TestDiagonalAttacksCount(t *testing.T) {
	expectedCounts := map[uint64]int{
		// Edge squares (files a,h and ranks 1,8) have 7 attacks
		0: 7, 1: 7, 2: 7, 3: 7, 4: 7, 5: 7, 6: 7, 7: 7, // rank 1
		8: 7, 15: 7, 16: 7, 23: 7, 24: 7, 31: 7, 32: 7, 39: 7, // files a,h
		40: 7, 47: 7, 48: 7, 55: 7, 56: 7, 57: 7, 58: 7, 59: 7, 60: 7, 61: 7, 62: 7, 63: 7, // rank 8
		// Second rank/column squares have 9 attacks
		9: 9, 10: 9, 11: 9, 12: 9, 13: 9, 14: 9, // rank 2
		17: 9, 22: 9, 25: 9, 30: 9, 33: 9, 38: 9, 41: 9, 46: 9, 49: 9, 50: 9, 51: 9, 52: 9, 53: 9, 54: 9, // files b,g
		// Third rank/column squares have 11 attacks
		18: 11, 19: 11, 20: 11, 21: 11, // rank 3
		26: 11, 29: 11, 34: 11, 37: 11, 42: 11, 43: 11, 44: 11, 45: 11, // files c,f
		// Center squares have 13 attacks
		27: 13, 28: 13, 35: 13, 36: 13, // d4, e4, d5, e5
	}

	for square := uint64(0); square < 64; square++ {
		attacks := GetDiagonalAttacks(square)
		expectedCount, exists := expectedCounts[square]
		if !exists {
			t.Errorf("Square %d: No expected count defined", square)
			continue
		}

		if attacks.Count() != expectedCount {
			t.Errorf("Square %d: Expected %d attacks, got %d", square, expectedCount, attacks.Count())
		}
	}
}

func TestDiagonalAttacksInvalidInput(t *testing.T) {
	invalidSquares := []string{"", "a", "a9", "i1", "z5", "a0"}

	for _, square := range invalidSquares {
		attacks := GetDiagonalAttacksFromSquare(square)
		if !attacks.IsEmpty() {
			t.Errorf("Expected empty bitboard for invalid square '%s', got %d attacks", square, attacks.Count())
		}
	}
}

func TestGetValidDiagonalMoves(t *testing.T) {
	pos := NewPosition()

	pos.SetPieceAtSquare("a1", Bishop, White)
	pos.SetPieceAtSquare("b2", Pawn, Black)
	pos.SetPieceAtSquare("c3", Pawn, White)

	bishop := pos.GetPieceAtSquare("a1")
	moves := GetValidDiagonalMoves(pos, bishop)
	squares := moves.ToSquares()

	expected := []string{"b2"}
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
