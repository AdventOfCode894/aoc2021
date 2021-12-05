package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := solvePuzzle(os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type coordinate struct {
	X int
	Y int
}

func (c *coordinate) Add(other coordinate) {
	c.X += other.X
	c.Y += other.Y
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	ventMap := make(map[coordinate]uint)
	for {
		var p1, p2 coordinate
		if _, err := fmt.Fscanf(r, "%d,%d -> %d,%d\n", &p1.X, &p1.Y, &p2.X, &p2.Y); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read input line: %v", err)
			}
			break
		}

		var step coordinate
		if p1.X < p2.X {
			step.X = 1
		} else if p1.X > p2.X {
			step.X = -1
		}
		if p1.Y < p2.Y {
			step.Y = 1
		} else if p1.Y > p2.Y {
			step.Y = -1
		}

		// For part 1, only consider horizontal or vertical lines
		if step.X != 0 && step.Y != 0 {
			continue
		}

		p2.Add(step) // Include p2
		for p := p1; p != p2; p.Add(step) {
			ventMap[p]++
		}
	}

	dangerousPoints := 0
	for _, overlap := range ventMap {
		if overlap <= 1 {
			continue
		}
		dangerousPoints++
	}

	_, _ = fmt.Fprintf(w, "Number of dangerous points: %d\n", dangerousPoints)
	return nil
}
