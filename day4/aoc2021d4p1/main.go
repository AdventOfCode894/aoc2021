package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/AdventOfCode894/aoc2021/day4/aoc2021d4"
)

func main() {
	if err := solvePuzzle(os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)

	draws, err := aoc2021d4.ReadDraws(s)
	if err != nil {
		return fmt.Errorf("failed to read draw sequence: %v", err)
	}

	var firstWinningCard *aoc2021d4.BingoCard
	var firstWinningTurns int
	var firstWinningScore uint64
	for {
		card, err := aoc2021d4.ReadBingoCard(s)
		if err != nil {
			return fmt.Errorf("failed to read bingo card: %v", err)
		}
		if card == nil {
			break
		}

		turns := 0
		var draw uint64
		for _, draw = range draws {
			turns++
			if card.Mark(draw) {
				break
			}
		}
		if card.Won() {
			if firstWinningCard == nil || firstWinningTurns > turns {
				sum := card.UnmarkedSum()
				score := sum * draw

				firstWinningCard = card
				firstWinningTurns = turns
				firstWinningScore = score
			}
		}
	}

	if firstWinningCard == nil {
		_, _ = fmt.Fprintln(w, "No cards win")
		return nil
	}

	_, _ = fmt.Fprintf(w, "First winning card has score %d\n", firstWinningScore)
	return nil
}
