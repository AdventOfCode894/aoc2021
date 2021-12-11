package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	var energies [][]uint
	var newRow []uint
	row := 0
	for {
		var c rune
		if _, err := fmt.Fscanf(r, "%c", &c); err != nil {
			if !errors.Is(err, io.EOF) {
				return 0, fmt.Errorf("failed to read energy: %v", err)
			}
			break
		}
		if c == '\n' {
			energies = append(energies, newRow)
			newRow = nil
			row++
			continue
		}
		energy := uint(c - '0')
		newRow = append(newRow, energy)
	}
	if len(newRow) > 0 {
		energies = append(energies, newRow)
	}

	flashes := make([][]bool, len(energies))
	for i := range flashes {
		flashes[i] = make([]bool, len(energies[i]))
	}

	steps := 0
	for {
		simulateEnergy(energies, flashes)
		steps++
		if areAllFlashing(energies) {
			break
		}
	}

	return steps, nil
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

func areAllFlashing(energies [][]uint) bool {
	for _, row := range energies {
		for _, energy := range row {
			if energy > 0 {
				return false
			}
		}
	}
	return true
}
