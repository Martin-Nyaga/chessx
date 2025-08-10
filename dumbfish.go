package main

type Dumbfish struct{}

func (d Dumbfish) Name() string { return "Dumbfish" }

// SelectMove returns the first legal move and resulting position, if any.
func (d Dumbfish) SelectMove(pos *Position) (AppliedMove, bool) {
	legal := generateLegalMoves(pos)
	if len(legal) == 0 {
		return AppliedMove{}, false
	}
	return legal[0], true
}
