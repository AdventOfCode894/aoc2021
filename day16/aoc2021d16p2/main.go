package main

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/day16/aoc2021d16"

	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	hr := hex.NewDecoder(r)

	value, _, err := aoc2021d16.EvaluateExpression(hr)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	return int(value), nil
}
