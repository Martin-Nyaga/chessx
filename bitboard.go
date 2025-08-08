package main

import (
	"fmt"
	"math/bits"
	"strings"
)

type Bitboard uint64

func NewBitboard(value uint64) Bitboard {
	return Bitboard(value)
}

func EmptyBitboard() Bitboard {
	return Bitboard(0)
}

func (b Bitboard) IsEmpty() bool {
	return b == 0
}

func (b Bitboard) IsSet(index uint64) bool {
	return (b & (1 << index)) != 0
}

func (b Bitboard) Set(index uint64) Bitboard {
	return b | (1 << index)
}

func (b Bitboard) Clear(index uint64) Bitboard {
	return b &^ (1 << index)
}

func (b Bitboard) Toggle(index uint64) Bitboard {
	return b ^ (1 << index)
}

func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}

func (b Bitboard) FirstSet() uint64 {
	if b == 0 {
		return ^uint64(0)
	}
	return uint64(bits.TrailingZeros64(uint64(b)))
}

func (b Bitboard) LastSet() uint64 {
	if b == 0 {
		return ^uint64(0)
	}
	return uint64(63 - bits.LeadingZeros64(uint64(b)))
}

func (b Bitboard) ToIndexes() []uint64 {
	var indexes []uint64
	for i := uint64(0); i < 64; i++ {
		if b.IsSet(i) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (b Bitboard) ToFileRanks() [][2]int {
	var fileRanks [][2]int
	for i := uint64(0); i < 64; i++ {
		if b.IsSet(i) {
			file, rank := indexToFileRank(i)
			fileRanks = append(fileRanks, [2]int{file, rank})
		}
	}
	return fileRanks
}

func (b Bitboard) ToSquares() []string {
	var squares []string
	for i := uint64(0); i < 64; i++ {
		if b.IsSet(i) {
			file, rank := indexToFileRank(i)
			square := fmt.Sprintf("%c%d", 'a'+file, rank+1)
			squares = append(squares, square)
		}
	}
	return squares
}

func FromIndex(index uint64) Bitboard {
	if index >= 64 {
		return EmptyBitboard()
	}
	return Bitboard(1 << index)
}

func FromFileRank(file, rank int) Bitboard {
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return EmptyBitboard()
	}
	index := fileRankToIndex(file, rank)
	return FromIndex(index)
}

func FromSquare(square string) Bitboard {
	file, rank, ok := squareToFileRank(square)
	if !ok {
		return EmptyBitboard()
	}
	return FromFileRank(file, rank)
}

func (b Bitboard) And(other Bitboard) Bitboard {
	return b & other
}

func (b Bitboard) Or(other Bitboard) Bitboard {
	return b | other
}

func (b Bitboard) Xor(other Bitboard) Bitboard {
	return b ^ other
}

func (b Bitboard) Not() Bitboard {
	return ^b
}

func (b Bitboard) ShiftLeft(amount uint64) Bitboard {
	return b << amount
}

func (b Bitboard) ShiftRight(amount uint64) Bitboard {
	return b >> amount
}

func (b Bitboard) String() string {
	var sb strings.Builder
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			index := fileRankToIndex(file, rank)
			if b.IsSet(index) {
				sb.WriteString("1 ")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b Bitboard) ToUint64() uint64 {
	return uint64(b)
}
