package main

import "fmt"

type GeneratedMove struct {
	From       string
	To         string
	Notation   string
	IsCapture  bool
	Promotion  PieceKind
	Kind       PieceKind
	Color      Color
	IsCastle   bool
	CastleSide CastlingSide
}

func squareFromIndex(index uint64) string {
	file, rank := indexToFileRank(index)
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

func pieceSANLetter(kind PieceKind) string {
	switch kind {
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return ""
	}
}

// generatePossibleMoves enumerates pseudo-legal moves for the side to move (ignores checks)
// and returns them with simple SAN-like notation (captures marked with 'x').
func generatePossibleMoves(pos *Position) []GeneratedMove {
	var moves []GeneratedMove

	var enemyOccupancy Bitboard
	if pos.toMove == White {
		enemyOccupancy = pos.GetBlackOccupancy()
	} else {
		enemyOccupancy = pos.GetWhiteOccupancy()
	}

	for i := range pos.pieces {
		piece := &pos.pieces[i]
		if piece.Color != pos.toMove || piece.Location.IsEmpty() {
			continue
		}

		fromIndex := piece.Location.FirstSet()
		fromSquare := squareFromIndex(fromIndex)

		var destinations Bitboard
		switch piece.Kind {
		case Knight:
			destinations = GetPossibleKnightMoves(pos, piece)
		case King:
			destinations = GetPossibleKingMoves(pos, piece)
		case Rook:
			destinations = GetPossibleRayMoves(pos, piece).Orthogonal()
		case Bishop:
			destinations = GetPossibleRayMoves(pos, piece).Diagonal()
		case Queen:
			destinations = GetPossibleRayMoves(pos, piece).All()
		case Pawn:
			destinations = GetPossiblePawnMoves(pos, piece)
		default:
			destinations = EmptyBitboard()
		}

		for _, toIndex := range destinations.ToIndexes() {
			toSquare := squareFromIndex(toIndex)

			isCapture := false
			if piece.Kind == Pawn {
				// Pawn capture if destination occupied by enemy or equals en passant square
				isCapture = enemyOccupancy.IsSet(toIndex) || (!pos.GetEnpassant().IsEmpty() && pos.GetEnpassant().IsSet(toIndex))
			} else {
				isCapture = enemyOccupancy.IsSet(toIndex)
			}

			// Handle pawn promotions: when a pawn moves to last rank, emit 4 promotion variants
			if piece.Kind == Pawn {
				_, toRank := indexToFileRank(toIndex)
				isPromotionRank := (piece.Color == White && toRank == 7) || (piece.Color == Black && toRank == 0)
				if isPromotionRank {
					promotionKinds := []PieceKind{Rook, Bishop, Knight, Queen}
					for _, promo := range promotionKinds {
						var notation string
						if isCapture {
							notation = fmt.Sprintf("%cx%s=%s", fromSquare[0], toSquare, pieceSANLetter(promo))
						} else {
							notation = fmt.Sprintf("%s=%s", toSquare, pieceSANLetter(promo))
						}
						moves = append(moves, GeneratedMove{
							From:      fromSquare,
							To:        toSquare,
							Notation:  notation,
							IsCapture: isCapture,
							Promotion: promo,
							Kind:      piece.Kind,
							Color:     piece.Color,
						})
					}
					continue
				}
			}

			notation := ""
			if piece.Kind == Pawn {
				if isCapture {
					// Pawn capture notation: source file + 'x' + destination
					notation = fmt.Sprintf("%cx%s", fromSquare[0], toSquare)
				} else {
					notation = toSquare
				}
			} else {
				letter := pieceSANLetter(piece.Kind)
				if isCapture {
					notation = fmt.Sprintf("%sx%s", letter, toSquare)
				} else {
					notation = fmt.Sprintf("%s%s", letter, toSquare)
				}
			}

			moves = append(moves, GeneratedMove{
				From:      fromSquare,
				To:        toSquare,
				Notation:  notation,
				IsCapture: isCapture,
				Promotion: Empty,
				Kind:      piece.Kind,
				Color:     piece.Color,
			})
		}
	}

	// Castling moves (pseudo-legal; checks filtered later)
	addCastlingMoves(pos, &moves)

	return moves
}

// through squares for castling checks
var (
	whiteKingsideThrough  = []string{"e1", "f1", "g1"}
	whiteQueensideThrough = []string{"e1", "d1", "c1"}
	blackKingsideThrough  = []string{"e8", "f8", "g8"}
	blackQueensideThrough = []string{"e8", "d8", "c8"}
)

// addCastlingMoves appends pseudo-legal castling moves to dst if available and path squares are empty
func addCastlingMoves(pos *Position, dst *[]GeneratedMove) {
	add := func(color Color, side CastlingSide, fromSq, toSq string, emptySquares []string, right CastlingSide, notation string) {
		if pos.toMove != color || !pos.CanCastle(right) {
			return
		}
		fromFile, fromRank, ok := squareToFileRank(fromSq)
		if !ok {
			return
		}
		king := pos.GetPiece(fromFile, fromRank)
		if king == nil || king.Kind != King || king.Color != color {
			return
		}
		for _, sq := range emptySquares {
			f, r, ok2 := squareToFileRank(sq)
			if !ok2 || pos.GetPiece(f, r) != nil {
				return
			}
		}
		*dst = append(*dst, GeneratedMove{
			From:       fromSq,
			To:         toSq,
			Notation:   notation,
			IsCapture:  false,
			Promotion:  Empty,
			Kind:       King,
			Color:      color,
			IsCastle:   true,
			CastleSide: side,
		})
	}
	// White
	add(White, WhiteKingside, "e1", "g1", []string{"f1", "g1"}, WhiteKingside, "O-O")
	add(White, WhiteQueenside, "e1", "c1", []string{"b1", "c1", "d1"}, WhiteQueenside, "O-O-O")
	// Black
	add(Black, BlackKingside, "e8", "g8", []string{"f8", "g8"}, BlackKingside, "O-O")
	add(Black, BlackQueenside, "e8", "c8", []string{"b8", "c8", "d8"}, BlackQueenside, "O-O-O")
}

type AppliedMove struct {
	Move     GeneratedMove
	Position *Position
}

// generateLegalMoves enumerates legal moves by filtering out pseudo-legal moves that
// leave the moving side's king in check. Returns each move paired with its resulting position.
func generateLegalMoves(pos *Position) []AppliedMove {
	possible := generatePossibleMoves(pos)
	legal := make([]AppliedMove, 0, len(possible))
	for _, mv := range possible {
		after := pos.ApplyMove(mv)
		// Filter moves that leave king in check
		if after.IsKingInCheck(mv.Color) {
			continue
		}
		// For castling, ensure not castling through check
		if mv.IsCastle {
			if pos.IsCastlingThroughCheck(mv.Color, mv.CastleSide) {
				continue
			}
		}
		legal = append(legal, AppliedMove{Move: mv, Position: after})
	}
	return legal
}
