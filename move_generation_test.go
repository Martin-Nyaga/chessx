package main

import (
	"sort"
	"testing"
)

func TestGenerateLegalMoves_StartingPositionWhite(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(fen)
	if err != nil {
		t.Fatalf("failed to parse fen: %v", err)
	}
	got := generatePossibleMoves(pos)
	if len(got) != 20 {
		t.Fatalf("expected 20 moves, got %d", len(got))
	}
	notations := make([]string, 0, len(got))
	for _, m := range got {
		notations = append(notations, m.Notation)
	}
	// required subset
	required := []string{"e3", "e4", "Na3", "Nc3", "Nf3", "Nh3"}
	for _, r := range required {
		found := false
		for _, n := range notations {
			if n == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find move %s in generated list", r)
		}
	}
}

func TestGenerateLegalMoves_RookCaptureFiltering(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("a1", Rook, White)
	pos.SetPieceAtSquare("b1", Pawn, Black)
	pos.SetPieceAtSquare("a2", Pawn, White)
	pos.SetToMove(White)

	got := generatePossibleMoves(pos)
	// Only one rook move: capture on b1
	found := false
	for _, m := range got {
		if m.Notation == "Rxb1" {
			if !m.IsCapture || m.From != "a1" || m.To != "b1" {
				t.Fatalf("Rxb1 should be capture from a1 to b1, got %+v", m)
			}
			found = true
		}
		// ensure no illegal rook forward move past own pawn
		if m.From == "a1" && m.To == "a2" {
			t.Fatalf("rook should not move to a2 when own pawn blocks")
		}
	}
	if !found {
		t.Fatalf("expected Rxb1 in generated moves")
	}
}

func TestGenerateLegalMoves_PawnCaptureAndEnPassant(t *testing.T) {
	pos := NewPosition()
	pos.SetPieceAtSquare("e5", Pawn, White)
	pos.SetPieceAtSquare("d5", Pawn, Black)
	pos.SetEnpassant(fileRankToIndex(3, 5)) // d6
	pos.SetToMove(White)

	got := generatePossibleMoves(pos)
	notations := map[string]GeneratedMove{}
	for _, m := range got {
		notations[m.Notation] = m
	}

	if mv, ok := notations["exd6"]; !ok {
		t.Fatalf("expected en passant capture exd6 present")
	} else if !mv.IsCapture || mv.To != "d6" || mv.From != "e5" {
		t.Fatalf("exd6 should be a capture from e5 to d6, got %+v", mv)
	}

	if mv, ok := notations["e6"]; !ok {
		t.Fatalf("expected forward push e6 present")
	} else if mv.IsCapture {
		t.Fatalf("e6 should not be a capture")
	}
}

func TestGenerateLegalMoves_BlackStartingMoves(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"
	pos, err := ParseFEN(fen)
	if err != nil {
		t.Fatalf("failed to parse fen: %v", err)
	}
	applied := generateLegalMoves(pos)
	moves := make([]GeneratedMove, 0, len(applied))
	for _, ap := range applied {
		moves = append(moves, ap.Move)
	}
	if len(moves) != 20 {
		t.Fatalf("expected 20 moves for black, got %d", len(moves))
	}
	var notations []string
	for _, m := range moves {
		notations = append(notations, m.Notation)
	}
	sort.Strings(notations)
	// Knights must be available
	required := []string{"Na6", "Nc6", "Nf6", "Nh6"}
	for _, r := range required {
		found := false
		for _, n := range notations {
			if n == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find move %s in generated list", r)
		}
	}
}
