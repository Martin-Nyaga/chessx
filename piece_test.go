package main

import (
	"testing"
)

func TestPieceCoordinates(t *testing.T) {
	// Test piece at a1 (file 0, rank 0)
	piece := &Piece{
		Kind:     Rook,
		Color:    White,
		Location: FromSquare("a1"),
	}

	file, rank := piece.FileRank()
	if file != 0 {
		t.Errorf("Expected file 0 for a1, got %d", file)
	}
	if rank != 0 {
		t.Errorf("Expected rank 0 for a1, got %d", rank)
	}
	if piece.Square() != "a1" {
		t.Errorf("Expected square 'a1', got '%s'", piece.Square())
	}

	// Test piece at e4 (file 4, rank 3)
	piece2 := &Piece{
		Kind:     Bishop,
		Color:    Black,
		Location: FromSquare("e4"),
	}

	file2, rank2 := piece2.FileRank()
	if file2 != 4 {
		t.Errorf("Expected file 4 for e4, got %d", file2)
	}
	if rank2 != 3 {
		t.Errorf("Expected rank 3 for e4, got %d", rank2)
	}
	if piece2.Square() != "e4" {
		t.Errorf("Expected square 'e4', got '%s'", piece2.Square())
	}

	// Test piece at h8 (file 7, rank 7)
	piece3 := &Piece{
		Kind:     Queen,
		Color:    White,
		Location: FromSquare("h8"),
	}

	file3, rank3 := piece3.FileRank()
	if file3 != 7 {
		t.Errorf("Expected file 7 for h8, got %d", file3)
	}
	if rank3 != 7 {
		t.Errorf("Expected rank 7 for h8, got %d", rank3)
	}
	if piece3.Square() != "h8" {
		t.Errorf("Expected square 'h8', got '%s'", piece3.Square())
	}
}

func TestPieceEmptyLocation(t *testing.T) {
	piece := &Piece{
		Kind:     Rook,
		Color:    White,
		Location: EmptyBitboard(),
	}

	file, rank := piece.FileRank()
	if file != -1 {
		t.Errorf("Expected file -1 for empty location, got %d", file)
	}
	if rank != -1 {
		t.Errorf("Expected rank -1 for empty location, got %d", rank)
	}
	if piece.Square() != "" {
		t.Errorf("Expected empty square for empty location, got '%s'", piece.Square())
	}
}
