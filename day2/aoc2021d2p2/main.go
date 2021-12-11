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
	aim := uint(0)
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		command, _ := tr.NextString(' ')
		amount, _ := tr.NextUint(aocio.EOLDelim, 10)
		switch command {
		case "forward":
			horizontal += amount
			depth += aim * amount
		case "down":
			aim += amount
		case "up":
			if aim < amount {
				aim = 0
				break
			}
			aim -= amount
		default:
			return 0, fmt.Errorf("unknown command: %s", command)
		}
	}
	area := horizontal * depth
	return int(area), pr.Err()
}
