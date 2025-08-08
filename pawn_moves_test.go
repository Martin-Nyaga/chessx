package main

import (
	"sort"
	"testing"
)

func TestGetValidPawnMoves_WhiteBasic(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e2", Pawn, White)
	moves := GetValidPawnMoves(pos, pos.GetPieceAtSquare("e2")).ToSquares()
	sort.Strings(moves)
	expected := []string{"e3", "e4"}
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

func TestGetValidPawnMoves_BlackBasic(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("d7", Pawn, Black)
	moves := GetValidPawnMoves(pos, pos.GetPieceAtSquare("d7")).ToSquares()
	sort.Strings(moves)
	expected := []string{"d5", "d6"}
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

func TestGetValidPawnMoves_CapturesAndBlocks(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e4", Pawn, White)
	pos.SetPieceAtSquare("e5", Pawn, White) // blocks push
	pos.SetPieceAtSquare("d5", Pawn, Black) // capturable
	pos.SetPieceAtSquare("f5", Pawn, Black) // capturable
	moves := GetValidPawnMoves(pos, pos.GetPieceAtSquare("e4")).ToSquares()
	sort.Strings(moves)
	expected := []string{"d5", "f5"}
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

func TestGetValidPawnMoves_EnPassant(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e5", Pawn, White)
	pos.SetPieceAtSquare("d5", Pawn, Black)
	// en passant target square d6
	pos.SetEnpassant(fileRankToIndex(3, 5))
	moves := GetValidPawnMoves(pos, pos.GetPieceAtSquare("e5")).ToSquares()
	sort.Strings(moves)
	expected := []string{"d6", "e6"}
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
