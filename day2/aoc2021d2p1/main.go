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
	pr := aocio.NewPuzzleReader(r)
	horizontal := uint(0)
	depth := uint(0)
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		command, _ := tr.NextString(' ')
		amount, _ := tr.NextUint(aocio.EOLDelim, 10)
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
	area := horizontal * depth
	return int(area), pr.Err()
}
