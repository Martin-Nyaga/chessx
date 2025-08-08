package main

import (
	"fmt"
	"os"
	"testing"
)

func TestBitboardBasics(t *testing.T) {
	empty := EmptyBitboard()
	if !empty.IsEmpty() {
		t.Errorf("Empty bitboard should be empty")
	}

	bb := NewBitboard(0x1234567890ABCDEF)
	if bb.IsEmpty() {
		t.Errorf("Non-zero bitboard should not be empty")
	}

	if bb.ToUint64() != 0x1234567890ABCDEF {
		t.Errorf("Expected 0x1234567890ABCDEF, got %x", bb.ToUint64())
	}
}

func TestBitboardOperations(t *testing.T) {
	bb := EmptyBitboard()

	bb = bb.Set(0)
	if !bb.IsSet(0) {
		t.Errorf("Bit 0 should be set")
	}

	bb = bb.Set(63)
	if !bb.IsSet(63) {
		t.Errorf("Bit 63 should be set")
	}

	bb = bb.Clear(0)
	if bb.IsSet(0) {
		t.Errorf("Bit 0 should be cleared")
	}

	bb = bb.Toggle(0)
	if !bb.IsSet(0) {
		t.Errorf("Bit 0 should be toggled on")
	}

	bb = bb.Toggle(0)
	if bb.IsSet(0) {
		t.Errorf("Bit 0 should be toggled off")
	}

	count := bb.Count()
	if count != 1 {
		t.Errorf("Expected 1 bit set, got %d", count)
	}
}

func TestBitboardConversions(t *testing.T) {
	bb := FromIndex(0)
	if !bb.IsSet(0) {
		t.Errorf("FromIndex(0) should set bit 0")
	}

	bb = FromFileRank(0, 0)
	if !bb.IsSet(0) {
		t.Errorf("FromFileRank(0,0) should set bit 0")
	}

	bb = FromSquare("a1")
	if !bb.IsSet(0) {
		t.Errorf("FromSquare(\"a1\") should set bit 0")
	}

	bb = FromSquare("h8")
	if !bb.IsSet(63) {
		t.Errorf("FromSquare(\"h8\") should set bit 63")
	}

	bb = FromSquare("e4")
	indexes := bb.ToIndexes()
	if len(indexes) != 1 || indexes[0] != 28 {
		t.Errorf("Expected [28], got %v", indexes)
	}

	fileRanks := bb.ToFileRanks()
	if len(fileRanks) != 1 || fileRanks[0][0] != 4 || fileRanks[0][1] != 3 {
		t.Errorf("Expected [[4 3]], got %v", fileRanks)
	}

	squares := bb.ToSquares()
	if len(squares) != 1 || squares[0] != "e4" {
		t.Errorf("Expected [\"e4\"], got %v", squares)
	}
}

func TestBitboardBitwise(t *testing.T) {
	bb1 := FromIndex(0)
	bb2 := FromIndex(1)

	and := bb1.And(bb2)
	if !and.IsEmpty() {
		t.Errorf("AND of different bits should be empty")
	}

	or := bb1.Or(bb2)
	if or.Count() != 2 {
		t.Errorf("OR should have 2 bits set")
	}

	xor := bb1.Xor(bb2)
	if xor.Count() != 2 {
		t.Errorf("XOR should have 2 bits set")
	}

	not := bb1.Not()
	if not.Count() != 63 {
		t.Errorf("NOT should have 63 bits set")
	}
}

func TestBitboardShifts(t *testing.T) {
	bb := FromIndex(0)

	shifted := bb.ShiftLeft(1)
	if !shifted.IsSet(1) {
		t.Errorf("Shift left should move bit to position 1")
	}

	shifted = bb.ShiftRight(1)
	if !shifted.IsEmpty() {
		t.Errorf("Shift right of bit 0 should be empty")
	}
}

func TestBitboardFirstLast(t *testing.T) {
	bb := FromIndex(5).Or(FromIndex(10)).Or(FromIndex(15))

	first := bb.FirstSet()
	if first != 5 {
		t.Errorf("First set bit should be 5, got %d", first)
	}

	last := bb.LastSet()
	if last != 15 {
		t.Errorf("Last set bit should be 15, got %d", last)
	}

	empty := EmptyBitboard()
	first = empty.FirstSet()
	if first != ^uint64(0) {
		t.Errorf("First set of empty should be max uint64, got %d", first)
	}

	last = empty.LastSet()
	if last != ^uint64(0) {
		t.Errorf("Last set of empty should be max uint64, got %d", last)
	}
}

func TestBitboardString(t *testing.T) {
	bb := FromSquare("e4")

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Bitboard for e4:")
		fmt.Println(bb.String())
	}

	expected := ". . . . . . . . \n. . . . . . . . \n. . . . . . . . \n. . . . . . . . \n. . . . 1 . . . \n. . . . . . . . \n. . . . . . . . \n. . . . . . . . \n"
	if bb.String() != expected {
		t.Errorf("String representation doesn't match expected")
	}
}

func TestBitboardMultipleSquares(t *testing.T) {
	bb := FromSquare("a1").Or(FromSquare("h8")).Or(FromSquare("e4"))

	squares := bb.ToSquares()
	expected := []string{"a1", "e4", "h8"}

	if len(squares) != len(expected) {
		t.Errorf("Expected %d squares, got %d", len(expected), len(squares))
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Bitboard for a1, e4, h8:")
		fmt.Println(bb.String())
		fmt.Printf("Squares: %v\n", squares)
	}
}
