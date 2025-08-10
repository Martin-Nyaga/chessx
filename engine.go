package main

// Engine is a minimal interface for a move-selecting engine.
// It returns an applied legal move and the resulting position, or false when no moves exist.
type Engine interface {
	Name() string
	SelectMove(pos *Position) (AppliedMove, bool)
}
