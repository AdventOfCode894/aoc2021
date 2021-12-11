package main

import (
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

type coordinate struct {
	X int
	Y int
}

func (c *coordinate) Add(other coordinate) {
	c.X += other.X
	c.Y += other.Y
}

func solvePuzzle(r io.Reader) (int, error) {
	ventMap := make(map[coordinate]uint)
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		var p1, p2 coordinate
		p1.X, _ = tr.NextInt(',', 10)
		p1.Y, _ = tr.NextInt(' ', 10)
		tr.ConsumeString("->")
		tr.ConsumeSpaces()
		p2.X, _ = tr.NextInt(',', 10)
		p2.Y, _ = tr.NextInt(aocio.EOLDelim, 10)
		tr.ConsumeEOL()

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
	if pr.Err() != nil {
		return 0, pr.Err()
	}

	dangerousPoints := 0
	for _, overlap := range ventMap {
		if overlap <= 1 {
			continue
		}
		dangerousPoints++
	}

	return dangerousPoints, nil
}
