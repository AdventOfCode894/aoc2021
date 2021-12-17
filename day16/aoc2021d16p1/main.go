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

	var parser aoc2021d16.ExpressionParser
	_, versionSum, err := parser.Evaluate(hr)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	return int(versionSum), nil
}
