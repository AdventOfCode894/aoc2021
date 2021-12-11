package main

import (
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	var tree diagnosticTreeNode
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		diagnostic := pr.ReadUintLine(2)
		tree.Insert(diagnostic, pr.LineLen())
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	oxygenGenRating := tree.FindRating(true)
	co2ScrubberRating := tree.FindRating(false)
	lifeSupportRating := oxygenGenRating * co2ScrubberRating
	return int(lifeSupportRating), nil
}
