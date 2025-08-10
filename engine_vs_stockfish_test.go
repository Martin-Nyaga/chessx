package main

import (
	"bufio"
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
	"time"
)

func engineLegalMovesUCI(pos *Position) []string {
	applied := generateLegalMoves(pos)
	ucis := make([]string, 0, len(applied))
	for _, ap := range applied {
		ucis = append(ucis, ap.Move.UCINotation())
	}
	sort.Strings(ucis)
	return ucis
}

func stockfishLegalMovesForHistoryUCI(t *testing.T, history []string) ([]string, error) {
	if _, err := exec.LookPath("stockfish"); err != nil {
		t.Skip("stockfish not found in PATH; skipping Stockfish comparison test")
		return nil, nil
	}
	cmd := exec.Command("stockfish")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	// Send commands
	movesLine := "position startpos"
	if len(history) > 0 {
		movesLine += " moves " + strings.Join(history, " ")
	}
	ioBuf := bytes.NewBufferString("uci\nucinewgame\n" + movesLine + "\ngo perft 1\nquit\n")
	if _, err := ioBuf.WriteTo(stdin); err != nil {
		return nil, err
	}
	_ = stdin.Close()
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	// Parse output like: "e2e4: 20"
	var ucis []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ":"); idx != -1 {
			move := strings.TrimSpace(line[:idx])
			// basic validation of UCI shape
			if len(move) >= 4 && len(move) <= 5 {
				ucis = append(ucis, move)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.Strings(ucis)
	return ucis, nil
}

func stockfishPositionStringForHistory(t *testing.T, history []string) (string, error) {
	if _, err := exec.LookPath("stockfish"); err != nil {
		t.Skip("stockfish not found in PATH; skipping Stockfish comparison test")
		return "", nil
	}
	cmd := exec.Command("stockfish")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	movesLine := "position startpos"
	if len(history) > 0 {
		movesLine += " moves " + strings.Join(history, " ")
	}
	ioBuf := bytes.NewBufferString("uci\nucinewgame\n" + movesLine + "\nd\nquit\n")
	if _, err := ioBuf.WriteTo(stdin); err != nil {
		return "", err
	}
	_ = stdin.Close()
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func TestEngineMovesMatchStockfish_StartPosition(t *testing.T) {
	if os.Getenv("CHESSX_STOCKFISH") != "1" {
		t.Skip("CHESSX_STOCKFISH env not set; skipping Stockfish comparison test")
		return
	}

	rand.Seed(time.Now().UnixNano())

	// Start from initial position and empty history
	pos, err := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatalf("parse FEN: %v", err)
	}
	history := make([]string, 0, 40)

	for step := 0; step < 21; step++ { // starting position + 20 random plies
		engineMoves := engineLegalMovesUCI(pos)
		stockMoves, err := stockfishLegalMovesForHistoryUCI(t, history)
		if err != nil {
			t.Fatalf("stockfish error: %v", err)
		}
		if stockMoves == nil {
			return // skipped
		}

		// Optional progress logs
		if os.Getenv("CHESSX_VERBOSE") == "1" {
			t.Logf("step %d history=%s", step, strings.Join(history, " "))
			t.Logf("counts engine=%d stockfish=%d", len(engineMoves), len(stockMoves))
			t.Logf("engine moves: %s", strings.Join(engineMoves, ","))
			t.Logf("stock  moves: %s", strings.Join(stockMoves, ","))
			if os.Getenv("CHESSX_EXTRA_VERBOSE") == "1" {
				t.Logf("engine position:\n%s", pos.String())
				if s, err := stockfishPositionStringForHistory(t, history); err == nil && s != "" {
					t.Logf("stockfish position:\n%s", s)
				}
			}
		}

		// Compare sets
		if len(engineMoves) != len(stockMoves) || !equalSorted(engineMoves, stockMoves) {
			if os.Getenv("CHESSX_VERBOSE") == "1" {
				t.Logf("step %d history=%s", step, strings.Join(history, " "))
				t.Logf("engine: %s", strings.Join(engineMoves, ","))
				t.Logf("stock : %s", strings.Join(stockMoves, ","))
			}
			t.Fatalf("move sets differ at step %d", step)
		}

		if step == 20 {
			break
		}
		// Pick a random legal move from engine set
		pick := engineMoves[rand.Intn(len(engineMoves))]
		if os.Getenv("CHESSX_VERBOSE") == "1" {
			t.Logf("selected move: %s", pick)
		}
		history = append(history, pick)
		// Apply to engine position by matching the GeneratedMove with this UCI
		applied := generateLegalMoves(pos)
		var found bool
		for _, ap := range applied {
			if ap.Move.UCINotation() == pick {
				pos = ap.Position
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("selected move %s not found in engine move list", pick)
		}
	}
}

func equalSorted(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
