package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if err := solvePuzzle(os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	var tree diagnosticTreeNode
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		diagnostic, err := strconv.ParseUint(line, 2, 64)
		if err != nil {
			return fmt.Errorf("failed to parse diagnostic \"%s\" as binary: %v", line, err)
		}
		tree.Insert(diagnostic, len(line))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read input line: %v", err)
	}

	oxygenGenRating := tree.FindRating(true)
	co2ScrubberRating := tree.FindRating(false)
	lifeSupportRating := oxygenGenRating * co2ScrubberRating
	_, _ = fmt.Fprintf(w, "Life support rating: %d\n  Oxygen generator rating: %d\n  CO2 scrubber rating: %d\n", lifeSupportRating, oxygenGenRating, co2ScrubberRating)

	return nil
}
