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
	pr := aocio.NewPuzzleReader(r)
	lastDepth := ^uint(0)
	increases := 0
	for pr.NextNonEmptyLine() {
		depth := pr.ReadUintLine(10)
		if depth > lastDepth {
			increases++
		}
		lastDepth = depth
	}
	return increases, pr.Err()
}
