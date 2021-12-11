package main

import (
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	horizontal := uint(0)
	depth := uint(0)
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		command, _ := tr.NextString(' ')
		amount, _ := tr.NextUint(aocio.EOLDelim, 10)
		tr.ConsumeEOL()
		switch command {
		case "forward":
			horizontal += amount
		case "down":
			depth += amount
		case "up":
			if depth < amount {
				depth = 0
				break
			}
			depth -= amount
		default:
			return 0, fmt.Errorf("unknown command: %s", command)
		}
	}
	if pr.Err() != nil {
		return 0, pr.Err()
	}
	area := horizontal * depth
	return int(area), nil
}
