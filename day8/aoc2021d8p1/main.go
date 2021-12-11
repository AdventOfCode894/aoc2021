package main

import (
	"errors"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	// Sorry for this mess! Yikes.
	uniqueValues := 0
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		tokens := pr.ReadStringArrayLine(' ')
		if err := pr.Err(); err != nil {
			return 0, err
		}
		if len(tokens) != 10+1+4 || tokens[10] != "|" {
			return 0, errors.New("invalid input data")
		}
		outputs := tokens[11:]

		for _, o := range outputs {
			switch len(o) {
			case 2: //1
				uniqueValues++
			case 4: //4
				uniqueValues++
			case 3: //7
				uniqueValues++
			case 7: //8
				uniqueValues++
			}
		}
	}

	return uniqueValues, nil
}
