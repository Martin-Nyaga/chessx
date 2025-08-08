package main

import (
	"sort"
	"testing"
)

func TestKingAttacks_Basic(t *testing.T) {
	tests := []struct {
		square   string
		expected []string
	}{
		{square: "a1", expected: []string{"a2", "b1", "b2"}},
		{square: "e4", expected: []string{"d3", "e3", "f3", "d4", "f4", "d5", "e5", "f5"}},
		{square: "h8", expected: []string{"g7", "g8", "h7"}},
	}

	for _, tt := range tests {
		t.Run(tt.square, func(t *testing.T) {
			att := GetKingAttacksFromSquare(tt.square)
			got := att.ToSquares()
			sort.Strings(got)
			sort.Strings(tt.expected)
			if len(got) != len(tt.expected) {
				t.Fatalf("Expected %d squares, got %d: %v", len(tt.expected), len(got), got)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Mismatch at %d: expected %s, got %s", i, tt.expected[i], got[i])
				}
			}
		})
	}
}

func TestGetValidKingMoves(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e4", King, White)
	pos.SetPieceAtSquare("e5", Pawn, White)
	pos.SetPieceAtSquare("f5", Pawn, Black)
	king := pos.GetPieceAtSquare("e4")
	moves := GetValidKingMoves(pos, king).ToSquares()
	sort.Strings(moves)
	expected := []string{"d3", "e3", "f3", "d4", "f4", "d5", "f5"}
	sort.Strings(expected)
	if len(moves) != len(expected) {
		t.Fatalf("Expected %d moves, got %d: %v", len(expected), len(moves), moves)
	}
	for i := range moves {
		if moves[i] != expected[i] {
			t.Errorf("Mismatch at %d: expected %s, got %s", i, expected[i], moves[i])
		}
	}
}
