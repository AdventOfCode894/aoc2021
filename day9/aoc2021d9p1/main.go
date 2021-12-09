package main

import (
	"errors"
	"fmt"
	"io"
	"os"
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

	_, _ = fmt.Fprintf(w, "Total risk: %d\n", totalRisk)
	return nil
}
