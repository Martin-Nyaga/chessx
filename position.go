package main

import (
	"fmt"
	"strings"
)

type Color int

const (
	White Color = iota
	Black
)

const (
	WhiteKingside  = 1 << 0
	WhiteQueenside = 1 << 1
	BlackKingside  = 1 << 2
	BlackQueenside = 1 << 3
)

type PieceKind int

const (
	Empty PieceKind = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

type Piece struct {
	Kind     PieceKind
	Color    Color
	Location Bitboard
}

type Position struct {
	pieces     []Piece
	board      [64]int
	toMove     Color
	moveNumber int
	castling   byte
	enpassant  Bitboard
	halfmoves  int
}

func NewPosition() *Position {
	pos := &Position{
		pieces:     make([]Piece, 0),
		board:      [64]int{},
		toMove:     White,
		moveNumber: 1,
		castling:   0,
		enpassant:  EmptyBitboard(),
		halfmoves:  0,
	}

	for i := 0; i < 64; i++ {
		pos.board[i] = -1
	}

	return pos
}

func fileRankToIndex(file, rank int) uint64 {
	return uint64(rank*8 + file)
}

func indexToFileRank(index uint64) (int, int) {
	file := int(index % 8)
	rank := int(index / 8)
	return file, rank
}

func (p *Position) SetPiece(file, rank int, kind PieceKind, color Color) {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return
	}

	index := fileRankToIndex(file, rank)

	if p.board[index] != -1 {
		pieceIndex := p.board[index]
		if pieceIndex < len(p.pieces) {
			p.pieces = append(p.pieces[:pieceIndex], p.pieces[pieceIndex+1:]...)
			for i := range p.board {
				if p.board[i] > pieceIndex {
					p.board[i]--
				}
			}
		}
		p.board[index] = -1
	}

	if kind != Empty {
		piece := Piece{
			Kind:     kind,
			Color:    color,
			Location: FromIndex(index),
		}
		p.pieces = append(p.pieces, piece)
		p.board[index] = len(p.pieces) - 1
	}
}

func (p *Position) GetPiece(file, rank int) *Piece {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return nil
	}

	index := fileRankToIndex(file, rank)
	pieceIndex := p.board[index]

	if pieceIndex == -1 || pieceIndex >= len(p.pieces) {
		return nil
	}

	return &p.pieces[pieceIndex]
}

func (p *Position) SetToMove(color Color) {
	p.toMove = color
}

func (p *Position) SetMoveNumber(moveNumber int) {
	p.moveNumber = moveNumber
}

func (p *Position) CanCastle(right byte) bool {
	return p.castling&right != 0
}

func (p *Position) SetCastling(right byte, available bool) {
	if available {
		p.castling |= right
	} else {
		p.castling &^= right
	}
}

func (p *Position) SetEnpassant(index uint64) {
	p.enpassant = FromIndex(index)
}

func (p *Position) GetEnpassant() Bitboard {
	return p.enpassant
}

func (p *Position) SetHalfmoves(halfmoves int) {
	p.halfmoves = halfmoves
}

func (p *Position) GetHalfmoves() int {
	return p.halfmoves
}

func (p *Position) String() string {
	var sb strings.Builder

	// Print the board
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			piece := p.GetPiece(file, rank)
			if piece == nil {
				sb.WriteString(". ")
			} else {
				sb.WriteString(pieceKindToFEN(piece.Kind, piece.Color))
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	// Print additional info
	sb.WriteString(fmt.Sprintf("To move: %s\n", colorToString(p.toMove)))
	sb.WriteString(fmt.Sprintf("Move: %d\n", p.moveNumber))
	sb.WriteString(fmt.Sprintf("Pieces: %d\n", len(p.pieces)))
	sb.WriteString(fmt.Sprintf("Castling: %08b\n", p.castling))
	squares := p.enpassant.ToSquares()
	if len(squares) == 0 {
		sb.WriteString("Enpassant: -\n")
	} else {
		sb.WriteString(fmt.Sprintf("Enpassant: %s\n", squares[0]))
	}
	sb.WriteString(fmt.Sprintf("Halfmoves: %d\n", p.halfmoves))

	return sb.String()
}

func pieceKindToFEN(kind PieceKind, color Color) string {
	var pieceChar string
	switch kind {
	case Pawn:
		pieceChar = "p"
	case Rook:
		pieceChar = "r"
	case Knight:
		pieceChar = "n"
	case Bishop:
		pieceChar = "b"
	case Queen:
		pieceChar = "q"
	case King:
		pieceChar = "k"
	default:
		return "."
	}

	if color == White {
		return strings.ToUpper(pieceChar)
	}
	return pieceChar
}

func colorToString(color Color) string {
	if color == White {
		return "White"
	}
	return "Black"
}

func ParseFEN(fen string) (*Position, error) {
	parts := strings.Fields(fen)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid FEN: insufficient parts")
	}

	pos := NewPosition()

	boardPart := parts[0]
	rank := 7
	file := 0

	for _, char := range boardPart {
		switch {
		case char == '/':
			rank--
			file = 0
		case char >= '1' && char <= '8':
			file += int(char - '0')
		case char == 'r':
			pos.SetPiece(file, rank, Rook, Black)
			file++
		case char == 'n':
			pos.SetPiece(file, rank, Knight, Black)
			file++
		case char == 'b':
			pos.SetPiece(file, rank, Bishop, Black)
			file++
		case char == 'q':
			pos.SetPiece(file, rank, Queen, Black)
			file++
		case char == 'k':
			pos.SetPiece(file, rank, King, Black)
			file++
		case char == 'p':
			pos.SetPiece(file, rank, Pawn, Black)
			file++
		case char == 'R':
			pos.SetPiece(file, rank, Rook, White)
			file++
		case char == 'N':
			pos.SetPiece(file, rank, Knight, White)
			file++
		case char == 'B':
			pos.SetPiece(file, rank, Bishop, White)
			file++
		case char == 'Q':
			pos.SetPiece(file, rank, Queen, White)
			file++
		case char == 'K':
			pos.SetPiece(file, rank, King, White)
			file++
		case char == 'P':
			pos.SetPiece(file, rank, Pawn, White)
			file++
		}
	}

	if len(parts) > 1 {
		if parts[1] == "w" {
			pos.SetToMove(White)
		} else if parts[1] == "b" {
			pos.SetToMove(Black)
		}
	}

	if len(parts) > 2 {
		castling := parts[2]
		pos.SetCastling(WhiteKingside, strings.ContainsRune(castling, 'K'))
		pos.SetCastling(WhiteQueenside, strings.ContainsRune(castling, 'Q'))
		pos.SetCastling(BlackKingside, strings.ContainsRune(castling, 'k'))
		pos.SetCastling(BlackQueenside, strings.ContainsRune(castling, 'q'))
	}

	if len(parts) > 3 {
		enpassant := parts[3]
		if enpassant != "-" {
			if len(enpassant) == 2 {
				file := int(enpassant[0] - 'a')
				rank := int(enpassant[1] - '1')
				if file >= 0 && file < 8 && rank >= 0 && rank < 8 {
					pos.SetEnpassant(fileRankToIndex(file, rank))
				}
			}
		}
	}

	if len(parts) > 4 {
		if _, err := fmt.Sscanf(parts[4], "%d", &pos.halfmoves); err != nil {
			pos.halfmoves = 0
		}
	}

	if len(parts) > 5 {
		if _, err := fmt.Sscanf(parts[5], "%d", &pos.moveNumber); err != nil {
			pos.moveNumber = 1
		}
	}

	return pos, nil
}
