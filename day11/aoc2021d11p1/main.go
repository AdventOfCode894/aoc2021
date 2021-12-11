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
	energies, _, _ := pr.Read2DUintArray(aocio.NoDelim, 10)
	if err := pr.Err(); err != nil {
		return 0, err
	}

	flashes := make([][]bool, len(energies))
	for i := range flashes {
		flashes[i] = make([]bool, len(energies[i]))
	}

	totalFlashes := 0
	for steps := 0; steps < 100; steps++ {
		totalFlashes += simulateEnergy(energies, flashes)
	}

	return totalFlashes, nil
}

func simulateEnergy(energies [][]uint, flashes [][]bool) int {
	for i, row := range energies {
		for j := range row {
			energies[i][j]++
			flashes[i][j] = false
		}
	}
	totalFlashes := 0
	for {
		newFlashes := applyFlashes(energies, flashes)
		if newFlashes < 1 {
			break
		}
		totalFlashes += newFlashes
	}
	for i, row := range energies {
		for j, energy := range row {
			if energy <= 9 {
				continue
			}
			energies[i][j] = 0
		}
	}
	return totalFlashes
}

func applyFlashes(energies [][]uint, flashes [][]bool) int {
	newFlashes := 0
	for i, row := range energies {
		for j, energy := range row {
			if flashes[i][j] {
				continue
			}
			if energy > 9 {
				for k := -1; k <= 1; k++ {
					if i+k < 0 || i+k >= len(energies) {
						continue
					}
					for l := -1; l <= 1; l++ {
						if k == 0 && l == 0 {
							continue
						}
						if j+l < 0 || j+l >= len(row) {
							continue
						}
						energies[i+k][j+l]++
					}
				}
				flashes[i][j] = true
				newFlashes++
			}
		}
	}
	return newFlashes
}
