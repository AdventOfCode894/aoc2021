package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"

	"github.com/AdventOfCode894/aoc2021/day4/aoc2021d4"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)

	draws, err := aoc2021d4.ReadDraws(pr)
	if err != nil {
		return 0, err
	}

	var firstWinningCard *aoc2021d4.BingoCard
	var firstWinningTurns int
	var firstWinningScore uint
	for {
		card, err := aoc2021d4.ReadBingoCard(pr)
		if err != nil {
			return 0, fmt.Errorf("failed to read bingo card: %v", err)
		}
		if card == nil {
			break
		}

		turns := 0
		var draw uint
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
	if pr.Err() != nil {
		return 0, pr.Err()
	}

	if firstWinningCard == nil {
		return 0, errors.New("no cards win")
	}

	return int(firstWinningScore), nil
}
