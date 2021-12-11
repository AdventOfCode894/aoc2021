package main

import (
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

const reproductionAge = 6
const newFishAge = 8
const simulationDays = 256

func solvePuzzle(r io.Reader) (int, error) {
	fishAges := make([]uint, newFishAge+1)
	pr := aocio.NewPuzzleReader(r)
	for _, age := range pr.ReadUintArrayLine(',', 10) {
		fishAges[age]++
	}
	if pr.Err() != nil {
		return 0, pr.Err()
	}

	for i := 0; i < simulationDays; i++ {
		reproducers := fishAges[0]
		copy(fishAges, fishAges[1:])
		fishAges[newFishAge] = reproducers
		fishAges[reproductionAge] += reproducers
	}

	totalFish := uint(0)
	for _, fish := range fishAges {
		totalFish += fish
	}
	return int(totalFish), nil
}
