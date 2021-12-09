package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	var in io.Reader = os.Stdin
	if len(os.Args) == 2 {
		var err error
		if in, err = os.Open(os.Args[1]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error attempting to open input file \"%s\": %v", os.Args[0], err)
		}
	}
	if err := solvePuzzle(in, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	var heights [][]uint
	var newRow []uint
	row := 0
	for {
		var c rune
		if _, err := fmt.Fscanf(r, "%c", &c); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read height: %v", err)
			}
			break
		}
		if c == '\n' {
			heights = append(heights, newRow)
			newRow = nil
			row++
			continue
		}
		height := uint(c - '0')
		newRow = append(newRow, height)
	}
	if len(newRow) > 0 {
		heights = append(heights, newRow)
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

	_, _ = fmt.Fprintf(w, "Largest three basins multiplied: %d\n", largestBasinMul)
	return nil
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
