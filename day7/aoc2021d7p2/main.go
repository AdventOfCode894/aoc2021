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
	var positions []uint
	maxPos := uint(0)
	minPos := ^uint(0)
	for {
		var position uint
		if _, err := fmt.Fscanf(r, "%d", &position); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read crab position: %v", err)
			}
			break
		}
		if maxPos < position {
			maxPos = position
		}
		if minPos > position {
			minPos = position
		}
		positions = append(positions, position)
	}

	minCost := ^uint(0)
	minCongregate := uint(0)
	for congregate := minPos; congregate <= maxPos; congregate++ {
		cost := uint(0)
		for _, position := range positions {
			distance := int(position) - int(congregate)
			if distance < 0 {
				distance = -distance
			}
			moveCost := (uint(distance) * (uint(distance) + 1)) / 2
			cost += moveCost
		}
		if cost < minCost {
			minCost = cost
			minCongregate = congregate
		}
	}

	_, _ = fmt.Fprintf(w, "Cheapest crab congregation point is %d with fuel cost %d\n", minCongregate, minCost)
	return nil
}
