package main

import (
	"os"
	"sort"
	"strings"
	"testing"
)

func TestKnightMoves_Basic(t *testing.T) {
	tests := []struct {
		square   string
		expected []string
	}{
		{square: "a1", expected: []string{"b3", "c2"}},
		{square: "e4", expected: []string{"c3", "c5", "d2", "d6", "f2", "f6", "g3", "g5"}},
		{square: "h8", expected: []string{"f7", "g6"}},
	}

	for _, tt := range tests {
		t.Run(tt.square, func(t *testing.T) {
			bb := GetKnightMovesFromSquare(tt.square)
			squares := bb.ToSquares()
			sort.Strings(squares)
			sort.Strings(tt.expected)
			if len(squares) != len(tt.expected) {
				t.Fatalf("Expected %d squares, got %d: %v", len(tt.expected), len(squares), squares)
			}
			for i := range squares {
				if squares[i] != tt.expected[i] {
					t.Errorf("Mismatch at %d: expected %s, got %s", i, tt.expected[i], squares[i])
				}
			}
		})
	}
}

func TestGetPossibleKnightMoves(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e4", Knight, White)
	pos.SetPieceAtSquare("c3", Pawn, White) // own blocker
	pos.SetPieceAtSquare("f6", Pawn, White) // own blocker
	pos.SetPieceAtSquare("d6", Pawn, Black) // capturable
	pos.SetPieceAtSquare("g5", Pawn, Black) // capturable

	knight := pos.GetPieceAtSquare("e4")
    moves := GetPossibleKnightMoves(pos, knight)
	got := moves.ToSquares()
	sort.Strings(got)
	expected := []string{"c5", "d2", "d6", "f2", "g3", "g5"}
	sort.Strings(expected)

	if len(got) != len(expected) {
		if os.Getenv("CHESSX_VERBOSE") == "1" {
			t.Logf("Knight moves diff. got=%s expected=%s", strings.Join(got, ","), strings.Join(expected, ","))
		}
		t.Fatalf("Expected %d moves, got %d", len(expected), len(got))
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("Mismatch at %d: expected %s, got %s", i, expected[i], got[i])
		}
	}
}

func TestKnightMoves_InvalidInput(t *testing.T) {
	invalid := []string{"", "a", "z9", "i1"}
	for _, s := range invalid {
		if !GetKnightMovesFromSquare(s).IsEmpty() {
			t.Errorf("Expected empty moves for invalid square %q", s)
		}
	}
}
