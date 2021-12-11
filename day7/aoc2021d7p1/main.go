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
	maxPos := uint(0)
	minPos := ^uint(0)
	pr := aocio.NewPuzzleReader(r)
	positions := pr.ReadUintArrayLine(',', 10)
	if pr.Err() != nil {
		return 0, pr.Err()
	}
	for _, position := range positions {
		if maxPos < position {
			maxPos = position
		}
		if minPos > position {
			minPos = position
		}
	}

	minCost := ^uint(0)
	for congregate := minPos; congregate <= maxPos; congregate++ {
		cost := uint(0)
		for _, position := range positions {
			distance := int(position) - int(congregate)
			if distance < 0 {
				distance = -distance
			}
			cost += uint(distance)
		}
		if cost < minCost {
			minCost = cost
		}
	}

	return int(minCost), nil
}
