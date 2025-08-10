package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func clearScreen() {
	// Try ANSI clear; fallback to calling clear
	fmt.Print("\033[2J\033[H")
	_ = exec.Command("clear").Run()
}

func readUserMove() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter move (UCI or SAN-like e2e4, Nf3, exd5, e8=Q), or 'q' to quit: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// matchInputToMove finds a legal move matching user's text. Accepts UCI fully; basic SAN-like fallback.
func matchInputToMove(pos *Position, input string) (AppliedMove, bool) {
	legal := generateLegalMoves(pos)
	// First try UCI exact
	for _, ap := range legal {
		if ap.Move.UCINotation() == input {
			return ap, true
		}
	}
	// Very basic SAN-like parsing: handle forms like e2e4, exd5, e8=Q already covered by UCI
	// Handle piece-letter destination like Nf3, Bb5: pick the first legal that has matching SAN letter and destination
	if len(input) >= 2 {
		dest := input[len(input)-2:]
		var letter string
		switch input[0] {
		case 'N':
			letter = "N"
		case 'B':
			letter = "B"
		case 'R':
			letter = "R"
		case 'Q':
			letter = "Q"
		case 'K':
			letter = "K"
		}
		for _, ap := range legal {
			if strings.HasSuffix(ap.Move.Notation, dest) {
				if letter == "" {
					// likely pawn move
					if ap.Move.Kind == Pawn {
						return ap, true
					}
				} else {
					if pieceSANLetter(ap.Move.Kind) == letter {
						return ap, true
					}
				}
			}
		}
	}
	return AppliedMove{}, false
}

func main() {
	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	pos, err := ParseFEN(startingFEN)
	if err != nil {
		fmt.Printf("Error parsing FEN: %v\n", err)
		return
	}

	var engine Engine = Dumbfish{}

	for {
		clearScreen()
		fmt.Printf("Side to move: %s\n\n", colorToString(pos.toMove))
		fmt.Println(pos.String())

		input, err := readUserMove()
		if err != nil {
			fmt.Printf("input error: %v\n", err)
			return
		}
		if input == "q" || input == "quit" || input == "exit" {
			fmt.Println("Goodbye!")
			return
		}

		userMove, ok := matchInputToMove(pos, input)
		if !ok {
			fmt.Println("Illegal or unrecognized move. Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			continue
		}
		pos = userMove.Position

		// Engine reply (Dumbfish)
		if reply, ok := engine.SelectMove(pos); ok {
			pos = reply.Position
		} else {
			fmt.Println("No legal moves for dumbfish. Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		}
	}
}
