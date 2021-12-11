package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/day4/aoc2021d4"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
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

	var lastWinningCard *aoc2021d4.BingoCard
	var lastWinningTurns int
	var lastWinningScore uint
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
			if lastWinningCard == nil || lastWinningTurns <= turns {
				sum := card.UnmarkedSum()
				score := sum * draw

				lastWinningCard = card
				lastWinningTurns = turns
				lastWinningScore = score
			}
		}
	}
	if pr.Err() != nil {
		return 0, pr.Err()
	}

	if lastWinningCard == nil {
		return 0, errors.New("no cards win")
	}

	return int(lastWinningScore), nil
}
