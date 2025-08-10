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

// CastlingSide represents a specific castling right bit value stored in Position.castling
type CastlingSide byte

const (
	WhiteKingside  CastlingSide = 1 << 0
	WhiteQueenside CastlingSide = 1 << 1
	BlackKingside  CastlingSide = 1 << 2
	BlackQueenside CastlingSide = 1 << 3
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

type Position struct {
	pieces         []Piece
	board          [64]int
	toMove         Color
	moveNumber     int
	castling       byte
	enpassant      Bitboard
	halfmoves      int
	whiteOccupancy Bitboard
	blackOccupancy Bitboard
}

func NewPosition() *Position {
	pos := &Position{
		pieces:         make([]Piece, 0),
		board:          [64]int{},
		toMove:         White,
		moveNumber:     1,
		castling:       0,
		enpassant:      EmptyBitboard(),
		halfmoves:      0,
		whiteOccupancy: EmptyBitboard(),
		blackOccupancy: EmptyBitboard(),
	}

	for i := 0; i < 64; i++ {
		pos.board[i] = -1
	}

	return pos
}

// moved to board_utils.go: fileRankToIndex, indexToFileRank

// Clone returns a deep copy of the position.
func (p *Position) Clone() *Position {
	if p == nil {
		return nil
	}
	newPosition := &Position{
		pieces:         make([]Piece, len(p.pieces)),
		board:          p.board,
		toMove:         p.toMove,
		moveNumber:     p.moveNumber,
		castling:       p.castling,
		enpassant:      p.enpassant,
		halfmoves:      p.halfmoves,
		whiteOccupancy: p.whiteOccupancy,
		blackOccupancy: p.blackOccupancy,
	}
	copy(newPosition.pieces, p.pieces)
	return newPosition
}

// ApplyMove returns a new position with the given move applied.
// Supports captures, en passant, promotions, en passant availability, halfmove clock, move number, and castling rights updates.
func (p *Position) ApplyMove(move GeneratedMove) *Position {
	if p == nil {
		return nil
	}
	newPosition := p.Clone()

	fromFile, fromRank, okFrom := squareToFileRank(move.From)
	toFile, toRank, okTo := squareToFileRank(move.To)
	if !okFrom || !okTo {
		return newPosition
	}

	occupant := newPosition.GetPiece(fromFile, fromRank)
	if occupant == nil {
		return newPosition
	}

	// Reset en passant by default; set if double pawn push occurs below
	newPosition.enpassant = EmptyBitboard()

	// Determine if this is an en passant capture on the current position
	isEnPassantCapture := false
	if move.Kind == Pawn && !p.enpassant.IsEmpty() {
		epIndex := p.enpassant.FirstSet()
		if epIndex == fileRankToIndex(toFile, toRank) && newPosition.GetPiece(toFile, toRank) == nil {
			isEnPassantCapture = true
		}
	}

	// Handle captures
	captured := false
	if isEnPassantCapture {
		captured = true
		capRank := toRank - 1
		if move.Color == Black {
			capRank = toRank + 1
		}
		newPosition.SetPiece(toFile, capRank, Empty, White)
	} else if newPosition.GetPiece(toFile, toRank) != nil {
		captured = true
	}

	// Update castling rights for king/rook moves and rook captures
	// King moves: clear both rights for that color
	if move.Kind == King {
		if move.Color == White {
			newPosition.SetCastling(WhiteKingside, false)
			newPosition.SetCastling(WhiteQueenside, false)
		} else {
			newPosition.SetCastling(BlackKingside, false)
			newPosition.SetCastling(BlackQueenside, false)
		}
	}
	// Rook moves from original squares
	if move.Kind == Rook {
		if move.Color == White && fromRank == 0 {
			if fromFile == 0 {
				newPosition.SetCastling(WhiteQueenside, false)
			} else if fromFile == 7 {
				newPosition.SetCastling(WhiteKingside, false)
			}
		} else if move.Color == Black && fromRank == 7 {
			if fromFile == 0 {
				newPosition.SetCastling(BlackQueenside, false)
			} else if fromFile == 7 {
				newPosition.SetCastling(BlackKingside, false)
			}
		}
	}
	// Rook captured on original squares
	if captured && !isEnPassantCapture {
		// Check original board (p) to see captured piece color/kind
		if capPiece := p.GetPiece(toFile, toRank); capPiece != nil && capPiece.Kind == Rook {
			if capPiece.Color == White && toRank == 0 {
				if toFile == 0 {
					newPosition.SetCastling(WhiteQueenside, false)
				} else if toFile == 7 {
					newPosition.SetCastling(WhiteKingside, false)
				}
			} else if capPiece.Color == Black && toRank == 7 {
				if toFile == 0 {
					newPosition.SetCastling(BlackQueenside, false)
				} else if toFile == 7 {
					newPosition.SetCastling(BlackKingside, false)
				}
			}
		}
	}

	// Move the piece (handle castling specially)
	if move.IsCastle {
		// Move king
		newPosition.SetPiece(fromFile, fromRank, Empty, move.Color)
		newPosition.SetPiece(toFile, toRank, King, move.Color)
		// Move rook accordingly
		if move.Color == White {
			if move.CastleSide == WhiteKingside {
				// rook h1 -> f1
				newPosition.SetPiece(7, 0, Empty, White)
				newPosition.SetPiece(5, 0, Rook, White)
			} else {
				// rook a1 -> d1
				newPosition.SetPiece(0, 0, Empty, White)
				newPosition.SetPiece(3, 0, Rook, White)
			}
			newPosition.SetCastling(WhiteKingside, false)
			newPosition.SetCastling(WhiteQueenside, false)
		} else {
			if move.CastleSide == BlackKingside {
				// rook h8 -> f8
				newPosition.SetPiece(7, 7, Empty, Black)
				newPosition.SetPiece(5, 7, Rook, Black)
			} else {
				// rook a8 -> d8
				newPosition.SetPiece(0, 7, Empty, Black)
				newPosition.SetPiece(3, 7, Rook, Black)
			}
			newPosition.SetCastling(BlackKingside, false)
			newPosition.SetCastling(BlackQueenside, false)
		}
		// No captures, en passant, or promotion during castling
		captured = false
		// Halfmove clock increments by 1 (not a pawn move or capture)
		newPosition.halfmoves = p.halfmoves + 1
		// Side to move and move number
		if p.toMove == Black {
			newPosition.moveNumber = p.moveNumber + 1
		} else {
			newPosition.moveNumber = p.moveNumber
		}
		if p.toMove == White {
			newPosition.toMove = Black
		} else {
			newPosition.toMove = White
		}
		return newPosition
	}
	newPosition.SetPiece(fromFile, fromRank, Empty, move.Color)
	movedKind := move.Kind
	if movedKind == Pawn && move.Promotion != Empty {
		movedKind = move.Promotion
	}
	newPosition.SetPiece(toFile, toRank, movedKind, move.Color)

	// En passant availability after a double pawn push
	if move.Kind == Pawn {
		if dr := toRank - fromRank; dr == 2 || dr == -2 {
			midRank := (toRank + fromRank) / 2
			newPosition.SetEnpassant(fileRankToIndex(toFile, midRank))
		}
	}

	// Halfmove clock
	if move.Kind == Pawn || captured {
		newPosition.halfmoves = 0
	} else {
		newPosition.halfmoves = p.halfmoves + 1
	}

	// Side to move and move number
	if p.toMove == Black {
		newPosition.moveNumber = p.moveNumber + 1
	} else {
		newPosition.moveNumber = p.moveNumber
	}
	if p.toMove == White {
		newPosition.toMove = Black
	} else {
		newPosition.toMove = White
	}

	return newPosition
}

// IsKingInCheck returns true if the specified color's king is attacked in this position.
func (p *Position) IsKingInCheck(color Color) bool {
	if p == nil {
		return false
	}
	// Find king index
	var kingIndex uint64 = ^uint64(0)
	for i := range p.pieces {
		piece := &p.pieces[i]
		if piece.Color == color && piece.Kind == King && !piece.Location.IsEmpty() {
			kingIndex = piece.Location.FirstSet()
			break
		}
	}
	if kingIndex == ^uint64(0) {
		return true
	}
	enemy := White
	if color == White {
		enemy = Black
	}
	for i := range p.pieces {
		piece := &p.pieces[i]
		if piece.Color != enemy || piece.Location.IsEmpty() {
			continue
		}
		var attacks Bitboard
		switch piece.Kind {
		case Knight:
			attacks = GetPossibleKnightMoves(p, piece)
		case King:
			attacks = GetPossibleKingMoves(p, piece)
		case Rook:
			attacks = GetPossibleRayMoves(p, piece).Orthogonal()
		case Bishop:
			attacks = GetPossibleRayMoves(p, piece).Diagonal()
		case Queen:
			attacks = GetPossibleRayMoves(p, piece).All()
		case Pawn:
			attacks = GetPossiblePawnMoves(p, piece)
		default:
			attacks = EmptyBitboard()
		}
		if attacks.IsSet(kingIndex) {
			return true
		}
	}
	return false
}

// IsCastleThroughCheck returns true if any intermediate square the king passes through
// (including the destination) is attacked for the given castling move.
func (p *Position) IsCastlingThroughCheck(color Color, side CastlingSide) bool {
	// Determine squares to test
	var throughSquares []string
	if color == White {
		if side == WhiteKingside {
			throughSquares = whiteKingsideThrough
		} else {
			throughSquares = whiteQueensideThrough
		}
	} else {
		if side == BlackKingside {
			throughSquares = blackKingsideThrough
		} else {
			throughSquares = blackQueensideThrough
		}
	}
	for _, sq := range throughSquares {
		f, r, ok := squareToFileRank(sq)
		if !ok {
			return true
		}
		// Temporarily place king at square and check attacks
		// Clone to avoid mutating current position
		tmp := p.Clone()
		// Remove any king of this color
		for i := range tmp.pieces {
			piece := &tmp.pieces[i]
			if piece.Color == color && piece.Kind == King && !piece.Location.IsEmpty() {
				file, rank := indexToFileRank(piece.Location.FirstSet())
				tmp.SetPiece(file, rank, Empty, color)
				break
			}
		}
		tmp.SetPiece(f, r, King, color)
		if tmp.IsKingInCheck(color) {
			return true
		}
	}
	return false
}

func (p *Position) SetPiece(file, rank int, kind PieceKind, color Color) {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return
	}

	index := fileRankToIndex(file, rank)

	// Remove existing piece from occupancy bitboards
	if p.board[index] != -1 {
		pieceIndex := p.board[index]
		if pieceIndex < len(p.pieces) {
			existingPiece := p.pieces[pieceIndex]
			if existingPiece.Color == White {
				p.whiteOccupancy = p.whiteOccupancy.Clear(index)
			} else {
				p.blackOccupancy = p.blackOccupancy.Clear(index)
			}
			p.pieces = append(p.pieces[:pieceIndex], p.pieces[pieceIndex+1:]...)
			for i := range p.board {
				if p.board[i] > pieceIndex {
					p.board[i]--
				}
			}
		}
		p.board[index] = -1
	}

	// Add new piece to occupancy bitboards (only if not Empty)
	if kind != Empty {
		piece := Piece{
			Kind:     kind,
			Color:    color,
			Location: FromIndex(index),
		}
		p.pieces = append(p.pieces, piece)
		p.board[index] = len(p.pieces) - 1

		if color == White {
			p.whiteOccupancy = p.whiteOccupancy.Set(index)
		} else {
			p.blackOccupancy = p.blackOccupancy.Set(index)
		}
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

func (p *Position) CanCastle(right CastlingSide) bool {
	return p.castling&byte(right) != 0
}

func (p *Position) SetCastling(right CastlingSide, available bool) {
	if available {
		p.castling |= byte(right)
	} else {
		p.castling &^= byte(right)
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

func (p *Position) GetWhiteOccupancy() Bitboard {
	return p.whiteOccupancy
}

func (p *Position) GetBlackOccupancy() Bitboard {
	return p.blackOccupancy
}

func (p *Position) GetAllOccupancy() Bitboard {
	return p.whiteOccupancy.Or(p.blackOccupancy)
}

func (p *Position) GetPieceAtSquare(square string) *Piece {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return nil
	}
	return p.GetPiece(file, rank)
}

func (p *Position) SetPieceAtSquare(square string, kind PieceKind, color Color) {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return
	}
	p.SetPiece(file, rank, kind, color)
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
			if file, rank, ok := squareToFileRank(enpassant); ok {
				pos.SetEnpassant(fileRankToIndex(file, rank))
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
