package main

import (
	"fmt"
)

func main() {
	fmt.Println("ChessX Move Generation Demo")
	fmt.Println("===========================")

	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(startingFEN)
	if err != nil {
		fmt.Printf("Error parsing FEN: %v\n", err)
		return
	}

	fmt.Println("Starting Position:")
	fmt.Println(pos.String())

	fmt.Printf("Generating legal moves for %s...\n\n", colorToString(pos.toMove))

	legal := generateLegalMoves(pos)
	for _, item := range legal {
		m := item.Move
		capMark := ""
		if m.IsCapture {
			capMark = " (capture)"
		}
		fmt.Printf("%s -> %s  %s%s\n", m.From, m.To, m.Notation, capMark)
		fmt.Println(item.Position.String())
	}
}
