package main

import "testing"

func TestRayAttacks_Basic(t *testing.T) {
	t.Run("diagonal_basic", func(t *testing.T) {
		tests := []struct {
			square      string
			expected    []string
			description string
		}{
			{square: "a1", expected: []string{"b2", "c3", "d4", "e5", "f6", "g7", "h8"}, description: "Corner square should attack all squares in positive diagonal"},
			{square: "e4", expected: []string{"a8", "b7", "c6", "d5", "f3", "g2", "h1", "b1", "c2", "d3", "f5", "g6", "h7"}, description: "Center square should attack all squares in both diagonals"},
			{square: "h8", expected: []string{"a1", "b2", "c3", "d4", "e5", "f6", "g7"}, description: "Corner square should attack all squares in negative diagonal"},
		}
		for _, test := range tests {
			t.Run(test.square, func(t *testing.T) {
				if len(test.square) != 2 {
					t.Fatalf("bad test square")
				}
				file := int(test.square[0] - 'a')
				rank := int(test.square[1] - '1')
				index := fileRankToIndex(file, rank)
				attacks := Rays.NE[index].Or(Rays.NW[index]).Or(Rays.SE[index]).Or(Rays.SW[index])
				squares := attacks.ToSquares()
				if len(squares) != len(test.expected) {
					t.Fatalf("Expected %d squares, got %d", len(test.expected), len(squares))
				}
				for _, expected := range test.expected {
					found := false
					for _, s := range squares {
						if s == expected {
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
	})

	t.Run("orthogonal_basic", func(t *testing.T) {
		tests := []struct {
			square      string
			expected    []string
			description string
		}{
			{square: "a1", expected: []string{"b1", "c1", "d1", "e1", "f1", "g1", "h1", "a2", "a3", "a4", "a5", "a6", "a7", "a8"}, description: "Corner square should attack all squares in same rank and file"},
			{square: "e4", expected: []string{"a4", "b4", "c4", "d4", "f4", "g4", "h4", "e1", "e2", "e3", "e5", "e6", "e7", "e8"}, description: "Center square should attack all squares in same rank and file"},
			{square: "h8", expected: []string{"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h1", "h2", "h3", "h4", "h5", "h6", "h7"}, description: "Corner square should attack all squares in same rank and file"},
		}
		for _, test := range tests {
			t.Run(test.square, func(t *testing.T) {
				if len(test.square) != 2 {
					t.Fatalf("bad test square")
				}
				file := int(test.square[0] - 'a')
				rank := int(test.square[1] - '1')
				index := fileRankToIndex(file, rank)
				attacks := Rays.N[index].Or(Rays.E[index]).Or(Rays.S[index]).Or(Rays.W[index])
				squares := attacks.ToSquares()
				if len(squares) != len(test.expected) {
					t.Fatalf("Expected %d squares, got %d", len(test.expected), len(squares))
				}
				for _, expected := range test.expected {
					found := false
					for _, s := range squares {
						if s == expected {
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
	})
}

func TestRayAttacks_Counts(t *testing.T) {
	t.Run("diagonal_counts_subset", func(t *testing.T) {
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
			attacks := Rays.NE[square].Or(Rays.NW[square]).Or(Rays.SE[square]).Or(Rays.SW[square])
			expectedCount, exists := expectedCounts[square]
			if !exists {
				continue
			}
			if attacks.Count() != expectedCount {
				t.Errorf("Square %d: Expected %d attacks, got %d", square, expectedCount, attacks.Count())
			}
		}
	})

	t.Run("orthogonal_counts_all", func(t *testing.T) {
		for square := uint64(0); square < 64; square++ {
			attacks := Rays.N[square].Or(Rays.E[square]).Or(Rays.S[square]).Or(Rays.W[square])
			expectedCount := 14
			if attacks.Count() != expectedCount {
				t.Errorf("Square %d: Expected %d attacks, got %d", square, expectedCount, attacks.Count())
			}
		}
	})
}

func TestRayAttacks_InvalidInput(t *testing.T) {
	invalidSquares := []string{"", "a", "a9", "i1", "z5", "a0"}
	for _, square := range invalidSquares {
		if attacks := GetRayAttacksFromSquare(square); !attacks.IsEmpty() {
			t.Errorf("Expected empty bitboard for invalid square '%s', got %d attacks", square, attacks.Count())
		}
	}
}

func TestGetValidMoves(t *testing.T) {
	t.Run("diagonal_valid_moves", func(t *testing.T) {
		pos := NewPosition()
		pos.SetPieceAtSquare("a1", Bishop, White)
		pos.SetPieceAtSquare("b2", Pawn, Black)
		pos.SetPieceAtSquare("c3", Pawn, White)

		bishop := pos.GetPieceAtSquare("a1")
		moves := GetValidRayMoves(pos, bishop)
		squares := moves.NE.Or(moves.NW).Or(moves.SE).Or(moves.SW).ToSquares()

		expected := []string{"b2"}
		for _, sq := range expected {
			found := false
			for _, s := range squares {
				if s == sq {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected square %s not found in moves", sq)
			}
		}
		if len(squares) != len(expected) {
			t.Errorf("Expected %d moves, got %d", len(expected), len(squares))
		}
	})

	t.Run("orthogonal_valid_moves", func(t *testing.T) {
		pos := NewPosition()
		pos.SetPieceAtSquare("a1", Rook, White)
		pos.SetPieceAtSquare("b1", Pawn, Black)
		pos.SetPieceAtSquare("a2", Pawn, White)

		rook := pos.GetPieceAtSquare("a1")
		moves := GetValidRayMoves(pos, rook)
		squares := moves.N.Or(moves.E).Or(moves.S).Or(moves.W).ToSquares()

		expected := []string{"b1"}
		for _, sq := range expected {
			found := false
			for _, s := range squares {
				if s == sq {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected square %s not found in moves", sq)
			}
		}
		if len(squares) != len(expected) {
			t.Errorf("Expected %d moves, got %d", len(expected), len(squares))
		}
	})
}
