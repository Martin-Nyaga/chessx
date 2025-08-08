package main

import (
	"fmt"
	"os"
	"testing"
)

func TestPositionCreation(t *testing.T) {
	// Create a new position
	pos := NewPosition()

	// Set up a sample position (starting position)
	// White pieces on rank 1
	pos.SetPiece(0, 0, Rook, White)
	pos.SetPiece(1, 0, Knight, White)
	pos.SetPiece(2, 0, Bishop, White)
	pos.SetPiece(3, 0, Queen, White)
	pos.SetPiece(4, 0, King, White)
	pos.SetPiece(5, 0, Bishop, White)
	pos.SetPiece(6, 0, Knight, White)
	pos.SetPiece(7, 0, Rook, White)

	// White pawns on rank 2
	for file := 0; file < 8; file++ {
		pos.SetPiece(file, 1, Pawn, White)
	}

	// Black pieces on rank 8
	pos.SetPiece(0, 7, Rook, Black)
	pos.SetPiece(1, 7, Knight, Black)
	pos.SetPiece(2, 7, Bishop, Black)
	pos.SetPiece(3, 7, Queen, Black)
	pos.SetPiece(4, 7, King, Black)
	pos.SetPiece(5, 7, Bishop, Black)
	pos.SetPiece(6, 7, Knight, Black)
	pos.SetPiece(7, 7, Rook, Black)

	// Black pawns on rank 7
	for file := 0; file < 8; file++ {
		pos.SetPiece(file, 6, Pawn, Black)
	}

	// Set some additional pieces for demonstration
	pos.SetPiece(3, 3, Queen, White)  // White queen in center
	pos.SetPiece(4, 4, Bishop, Black) // Black bishop in center

	// Set the position to White's turn, move 5
	pos.SetToMove(White)
	pos.SetMoveNumber(5)

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Sample Chess Position:")
		fmt.Println(pos.String())
	}

	// Verify some pieces are in the correct positions
	if piece := pos.GetPiece(0, 0); piece == nil || piece.Kind != Rook || piece.Color != White {
		t.Errorf("Expected white rook at a1, got %v", piece)
	}

	if piece := pos.GetPiece(4, 7); piece == nil || piece.Kind != King || piece.Color != Black {
		t.Errorf("Expected black king at e8, got %v", piece)
	}

	if piece := pos.GetPiece(3, 3); piece == nil || piece.Kind != Queen || piece.Color != White {
		t.Errorf("Expected white queen at d4, got %v", piece)
	}
}

func TestEmptyPosition(t *testing.T) {
	pos := NewPosition()

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Empty Position:")
		fmt.Println(pos.String())
	}

	// Verify all squares are empty
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if piece := pos.GetPiece(file, rank); piece != nil {
				t.Errorf("Expected empty square at %d,%d, got %v", file, rank, piece)
			}
		}
	}
}

func TestCastling(t *testing.T) {
	pos := NewPosition()

	pos.SetCastling(WhiteKingside, true)
	pos.SetCastling(BlackQueenside, true)

	if !pos.CanCastle(WhiteKingside) {
		t.Errorf("Expected white kingside castling to be available")
	}

	if pos.CanCastle(WhiteQueenside) {
		t.Errorf("Expected white queenside castling to be unavailable")
	}

	if !pos.CanCastle(BlackQueenside) {
		t.Errorf("Expected black queenside castling to be available")
	}

	if pos.CanCastle(BlackKingside) {
		t.Errorf("Expected black kingside castling to be unavailable")
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Castling Test:")
		fmt.Println(pos.String())
	}
}

func TestEnpassant(t *testing.T) {
	pos := NewPosition()

	pos.SetEnpassant(20)

	if pos.GetEnpassant().ToUint64() != (1 << 20) {
		t.Errorf("Expected enpassant at 20, got %d", pos.GetEnpassant().ToUint64())
	}

	pos.SetEnpassant(^uint64(0))

	if !pos.GetEnpassant().IsEmpty() {
		t.Errorf("Expected no enpassant, got %s", pos.GetEnpassant().String())
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Enpassant Test:")
		fmt.Println(pos.String())
	}
}

func TestHalfmoves(t *testing.T) {
	pos := NewPosition()

	pos.SetHalfmoves(15)

	if pos.GetHalfmoves() != 15 {
		t.Errorf("Expected 15 halfmoves, got %d", pos.GetHalfmoves())
	}

	pos.SetHalfmoves(0)

	if pos.GetHalfmoves() != 0 {
		t.Errorf("Expected 0 halfmoves, got %d", pos.GetHalfmoves())
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Halfmoves Test:")
		fmt.Println(pos.String())
	}
}

func TestParseFEN(t *testing.T) {
	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	pos, err := ParseFEN(startingFEN)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	if pos.toMove != White {
		t.Errorf("Expected white to move, got %v", pos.toMove)
	}

	if !pos.CanCastle(WhiteKingside) {
		t.Errorf("Expected white kingside castling to be available")
	}

	if !pos.CanCastle(WhiteQueenside) {
		t.Errorf("Expected white queenside castling to be available")
	}

	if !pos.CanCastle(BlackKingside) {
		t.Errorf("Expected black kingside castling to be available")
	}

	if !pos.CanCastle(BlackQueenside) {
		t.Errorf("Expected black queenside castling to be available")
	}

	if !pos.GetEnpassant().IsEmpty() {
		t.Errorf("Expected no enpassant, got %s", pos.GetEnpassant().String())
	}

	if pos.GetHalfmoves() != 0 {
		t.Errorf("Expected 0 halfmoves, got %d", pos.GetHalfmoves())
	}

	if pos.moveNumber != 1 {
		t.Errorf("Expected move 1, got %d", pos.moveNumber)
	}

	if piece := pos.GetPiece(0, 0); piece == nil || piece.Kind != Rook || piece.Color != White {
		t.Errorf("Expected white rook at a1, got %v", piece)
	}

	if piece := pos.GetPiece(4, 7); piece == nil || piece.Kind != King || piece.Color != Black {
		t.Errorf("Expected black king at e8, got %v", piece)
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Starting Position from FEN:")
		fmt.Println(pos.String())
	}
}

func TestSpecificFEN(t *testing.T) {
	fen := "3r3k/pb2q1pp/3b1p2/2n5/2QRpP2/6B1/PP4PP/2R3K1 b - - 7 28"

	pos, err := ParseFEN(fen)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	if pos.toMove != Black {
		t.Errorf("Expected black to move, got %v", pos.toMove)
	}

	if pos.castling != 0 {
		t.Errorf("Expected no castling rights, got %08b", pos.castling)
	}

	if !pos.GetEnpassant().IsEmpty() {
		t.Errorf("Expected no enpassant, got %s", pos.GetEnpassant().String())
	}

	if pos.GetHalfmoves() != 7 {
		t.Errorf("Expected 7 halfmoves, got %d", pos.GetHalfmoves())
	}

	if pos.moveNumber != 28 {
		t.Errorf("Expected move 28, got %d", pos.moveNumber)
	}

	if piece := pos.GetPiece(3, 7); piece == nil || piece.Kind != Rook || piece.Color != Black {
		t.Errorf("Expected black rook at d8, got %v", piece)
	}

	if piece := pos.GetPiece(7, 7); piece == nil || piece.Kind != King || piece.Color != Black {
		t.Errorf("Expected black king at h8, got %v", piece)
	}

	if piece := pos.GetPiece(2, 0); piece == nil || piece.Kind != Rook || piece.Color != White {
		t.Errorf("Expected white rook at c1, got %v", piece)
	}

	if piece := pos.GetPiece(6, 0); piece == nil || piece.Kind != King || piece.Color != White {
		t.Errorf("Expected white king at g1, got %v", piece)
	}

	if os.Getenv("CHESSX_VERBOSE") == "1" {
		fmt.Println("Specific FEN Position:")
		fmt.Println(pos.String())
	}
}

func TestOccupancyBitboards(t *testing.T) {
	pos := NewPosition()

	// Test empty position
	if !pos.GetWhiteOccupancy().IsEmpty() {
		t.Errorf("Expected empty white occupancy for new position")
	}
	if !pos.GetBlackOccupancy().IsEmpty() {
		t.Errorf("Expected empty black occupancy for new position")
	}
	if !pos.GetAllOccupancy().IsEmpty() {
		t.Errorf("Expected empty total occupancy for new position")
	}

	// Test adding white pieces
	pos.SetPiece(0, 0, Rook, White)
	pos.SetPiece(1, 0, Knight, White)
	pos.SetPiece(2, 0, Bishop, White)

	whiteOccupancy := pos.GetWhiteOccupancy()
	if whiteOccupancy.Count() != 3 {
		t.Errorf("Expected 3 white pieces, got %d", whiteOccupancy.Count())
	}
	if !whiteOccupancy.IsSet(0) {
		t.Errorf("Expected white piece at a1 (index 0)")
	}
	if !whiteOccupancy.IsSet(1) {
		t.Errorf("Expected white piece at b1 (index 1)")
	}
	if !whiteOccupancy.IsSet(2) {
		t.Errorf("Expected white piece at c1 (index 2)")
	}

	// Test adding black pieces
	pos.SetPiece(0, 7, Rook, Black)
	pos.SetPiece(1, 7, Knight, Black)
	pos.SetPiece(2, 7, Bishop, Black)

	blackOccupancy := pos.GetBlackOccupancy()
	if blackOccupancy.Count() != 3 {
		t.Errorf("Expected 3 black pieces, got %d", blackOccupancy.Count())
	}
	if !blackOccupancy.IsSet(56) {
		t.Errorf("Expected black piece at a8 (index 56)")
	}
	if !blackOccupancy.IsSet(57) {
		t.Errorf("Expected black piece at b8 (index 57)")
	}
	if !blackOccupancy.IsSet(58) {
		t.Errorf("Expected black piece at c8 (index 58)")
	}

	// Test total occupancy
	allOccupancy := pos.GetAllOccupancy()
	if allOccupancy.Count() != 6 {
		t.Errorf("Expected 6 total pieces, got %d", allOccupancy.Count())
	}

	// Test replacing pieces
	pos.SetPiece(0, 0, Queen, White) // Replace white rook with queen
	whiteOccupancy = pos.GetWhiteOccupancy()
	if whiteOccupancy.Count() != 3 {
		t.Errorf("Expected still 3 white pieces after replacement, got %d", whiteOccupancy.Count())
	}
	if !whiteOccupancy.IsSet(0) {
		t.Errorf("Expected white piece still at a1 after replacement")
	}

	// Test removing pieces
	pos.SetPiece(0, 0, Empty, White) // Remove piece
	whiteOccupancy = pos.GetWhiteOccupancy()
	if whiteOccupancy.Count() != 2 {
		t.Errorf("Expected 2 white pieces after removal, got %d", whiteOccupancy.Count())
	}
	if whiteOccupancy.IsSet(0) {
		t.Errorf("Expected no white piece at a1 after removal")
	}

	// Test replacing with different color
	pos.SetPiece(1, 0, Bishop, Black) // Replace white knight with black bishop
	whiteOccupancy = pos.GetWhiteOccupancy()
	blackOccupancy = pos.GetBlackOccupancy()
	if whiteOccupancy.Count() != 1 {
		t.Errorf("Expected 1 white piece after color change, got %d", whiteOccupancy.Count())
	}
	if blackOccupancy.Count() != 4 {
		t.Errorf("Expected 4 black pieces after color change, got %d", blackOccupancy.Count())
	}
	if whiteOccupancy.IsSet(1) {
		t.Errorf("Expected no white piece at b1 after color change")
	}
	if !blackOccupancy.IsSet(1) {
		t.Errorf("Expected black piece at b1 after color change")
	}
}

func TestOccupancyBitboardsWithFEN(t *testing.T) {
	// Test starting position
	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(startingFEN)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	whiteOccupancy := pos.GetWhiteOccupancy()
	blackOccupancy := pos.GetBlackOccupancy()
	allOccupancy := pos.GetAllOccupancy()

	if whiteOccupancy.Count() != 16 {
		t.Errorf("Expected 16 white pieces in starting position, got %d", whiteOccupancy.Count())
	}
	if blackOccupancy.Count() != 16 {
		t.Errorf("Expected 16 black pieces in starting position, got %d", blackOccupancy.Count())
	}
	if allOccupancy.Count() != 32 {
		t.Errorf("Expected 32 total pieces in starting position, got %d", allOccupancy.Count())
	}

	// Test specific squares in starting position
	if !whiteOccupancy.IsSet(0) { // a1
		t.Errorf("Expected white piece at a1")
	}
	if !whiteOccupancy.IsSet(4) { // e1
		t.Errorf("Expected white piece at e1")
	}
	if !blackOccupancy.IsSet(56) { // a8
		t.Errorf("Expected black piece at a8")
	}
	if !blackOccupancy.IsSet(60) { // e8
		t.Errorf("Expected black piece at e8")
	}

	// Test a more complex position
	complexFEN := "3r3k/pb2q1pp/3b1p2/2n5/2QRpP2/6B1/PP4PP/2R3K1 b - - 7 28"
	pos2, err := ParseFEN(complexFEN)
	if err != nil {
		t.Fatalf("Failed to parse complex FEN: %v", err)
	}

	whiteOccupancy2 := pos2.GetWhiteOccupancy()
	blackOccupancy2 := pos2.GetBlackOccupancy()
	allOccupancy2 := pos2.GetAllOccupancy()

	// Count pieces manually from FEN
	whitePieces := 0
	blackPieces := 0
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			piece := pos2.GetPiece(file, rank)
			if piece != nil {
				if piece.Color == White {
					whitePieces++
				} else {
					blackPieces++
				}
			}
		}
	}

	if whiteOccupancy2.Count() != whitePieces {
		t.Errorf("Expected %d white pieces, got %d", whitePieces, whiteOccupancy2.Count())
	}
	if blackOccupancy2.Count() != blackPieces {
		t.Errorf("Expected %d black pieces, got %d", blackPieces, blackOccupancy2.Count())
	}
	if allOccupancy2.Count() != whitePieces+blackPieces {
		t.Errorf("Expected %d total pieces, got %d", whitePieces+blackPieces, allOccupancy2.Count())
	}
}
