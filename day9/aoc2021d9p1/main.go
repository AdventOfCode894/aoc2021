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
	heights, _, _ := pr.Read2DUintArray(aocio.NoDelim, 10)
	if err := pr.Err(); err != nil {
		return 0, err
	}

	totalRisk := uint(0)
	for row, rowHeights := range heights {
		for col, height := range rowHeights {
			if row > 0 {
				if heights[row-1][col] <= height {
					continue
				}
			}
			if col > 0 {
				if heights[row][col-1] <= height {
					continue
				}
			}
			if row < len(heights)-1 {
				if heights[row+1][col] <= height {
					continue
				}
			}
			if col < len(rowHeights)-1 {
				if heights[row][col+1] <= height {
					continue
				}
			}
			risk := height + 1
			totalRisk += risk
		}
	}

	return int(totalRisk), nil
}
