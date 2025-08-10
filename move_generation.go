package main

import "fmt"

type GeneratedMove struct {
	From      string
	To        string
	Notation  string
	IsCapture bool
	Promotion PieceKind
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

// generateLegalMoves enumerates pseudo-legal moves for the side to move (ignores checks)
// and returns them with simple SAN-like notation (captures marked with 'x').
func generateLegalMoves(pos *Position) []GeneratedMove {
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
			})
		}
	}

	return moves
}
