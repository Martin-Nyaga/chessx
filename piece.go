package main

import "fmt"

type Piece struct {
	Kind     PieceKind
	Color    Color
	Location Bitboard
}

func (p *Piece) FileRank() (int, int) {
	if p.Location.IsEmpty() {
		return -1, -1
	}
	index := p.Location.FirstSet()
	if index >= 64 {
		return -1, -1
	}
	return indexToFileRank(index)
}

func (p *Piece) Square() string {
	if p.Location.IsEmpty() {
		return ""
	}
	index := p.Location.FirstSet()
	if index >= 64 {
		return ""
	}
	file, rank := indexToFileRank(index)
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}
