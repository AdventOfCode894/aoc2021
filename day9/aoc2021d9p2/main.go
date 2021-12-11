package main

import (
	"io"
	"sort"

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

	var basinSizes []uint
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
			basinSize := findBasinSize(heights, row, col)
			basinSizes = append(basinSizes, basinSize)
		}
	}

	sort.Slice(basinSizes, func(i, j int) bool {
		return basinSizes[i] > basinSizes[j]
	})

	largestBasinMul := basinSizes[0] * basinSizes[1] * basinSizes[2]
	return int(largestBasinMul), nil
}

type coordinate struct {
	row int
	col int
}

func findBasinSize(heights [][]uint, row int, col int) uint {
	visited := make([][]bool, len(heights))
	for i := range heights {
		visited[i] = make([]bool, len(heights[i]))
	}
	queue := []coordinate{{row, col}}
	basinSize := uint(0)
	for len(queue) > 0 {
		p := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if p.row < 0 {
			continue
		}
		if p.col < 0 {
			continue
		}
		if p.row >= len(visited) {
			continue
		}
		if p.col >= len(visited[p.row]) {
			continue
		}
		if visited[p.row][p.col] {
			continue
		}
		visited[p.row][p.col] = true
		if heights[p.row][p.col] >= 9 {
			continue
		}
		basinSize++
		queue = append(queue,
			coordinate{p.row - 1, p.col},
			coordinate{p.row + 1, p.col},
			coordinate{p.row, p.col - 1},
			coordinate{p.row, p.col + 1})
	}
	return basinSize
}
