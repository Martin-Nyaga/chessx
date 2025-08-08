package main

// GetValidPawnMoves returns legal pawn moves for the given pawn piece.
// Includes single/double pushes (if empty), captures, and en passant (if available in position).
// Assumes standard chess with white pawns moving towards increasing ranks (north).
func GetValidPawnMoves(pos *Position, piece *Piece) Bitboard {
	if piece == nil || piece.Location.IsEmpty() || piece.Kind != Pawn {
		return EmptyBitboard()
	}
	originIndex := piece.Location.FirstSet()
	if originIndex >= 64 {
		return EmptyBitboard()
	}

	file, rank := indexToFileRank(originIndex)

	var selfOccupancy, enemyOccupancy Bitboard
	var forwardDirection int
	var startRank int
	if piece.Color == White {
		selfOccupancy = pos.GetWhiteOccupancy()
		enemyOccupancy = pos.GetBlackOccupancy()
		forwardDirection = 1
		startRank = 1
	} else {
		selfOccupancy = pos.GetBlackOccupancy()
		enemyOccupancy = pos.GetWhiteOccupancy()
		forwardDirection = -1
		startRank = 6
	}

	moves := EmptyBitboard()

	// Single push
	nextRank := rank + forwardDirection
	if nextRank >= 0 && nextRank < 8 {
		nextIndex := fileRankToIndex(file, nextRank)
		if !pos.GetAllOccupancy().IsSet(nextIndex) {
			moves = moves.Set(nextIndex)

			// Double push from start rank
			if rank == startRank {
				doubleStepRank := rank + 2*forwardDirection
				if doubleStepRank >= 0 && doubleStepRank < 8 {
					doubleStepIndex := fileRankToIndex(file, doubleStepRank)
					if !pos.GetAllOccupancy().IsSet(doubleStepIndex) {
						moves = moves.Set(doubleStepIndex)
					}
				}
			}
		}
	}

	// Captures
	for _, deltaFile := range []int{-1, 1} {
		targetFile := file + deltaFile
		targetRank := rank + forwardDirection
		if targetFile >= 0 && targetFile < 8 && targetRank >= 0 && targetRank < 8 {
			targetIndex := fileRankToIndex(targetFile, targetRank)
			if enemyOccupancy.IsSet(targetIndex) {
				moves = moves.Set(targetIndex)
			}
		}
	}

	// En passant
	enPassant := pos.GetEnpassant()
	if !enPassant.IsEmpty() {
		enPassantIndex := enPassant.FirstSet()
		enPassantFile, enPassantRank := indexToFileRank(enPassantIndex)
		enPassantTargetRank := rank + forwardDirection
		for _, deltaFile := range []int{-1, 1} {
			targetFile := file + deltaFile
			if targetFile == enPassantFile && enPassantRank == enPassantTargetRank {
				moves = moves.Set(enPassantIndex)
			}
		}
	}

	// Remove own occupancy (safety, though pushes should be empty already)
	moves = moves.And(selfOccupancy.Not())
	return moves
}
