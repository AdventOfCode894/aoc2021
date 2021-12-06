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

const reproductionAge = 6
const newFishAge = 8
const simulationDays = 256

func solvePuzzle(r io.Reader, w io.Writer) error {
	fishAges := make([]uint, newFishAge+1)
	for {
		var age uint
		if _, err := fmt.Fscanf(r, "%d", &age); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read age: %v", err)
			}
			break
		}
		fishAges[age]++
	}

	for i := 0; i < simulationDays; i++ {
		reproducers := fishAges[0]
		copy(fishAges, fishAges[1:])
		fishAges[newFishAge] = reproducers
		fishAges[reproductionAge] += reproducers
	}

	totalFish := uint(0)
	for _, fish := range fishAges {
		totalFish += fish
	}
	_, _ = fmt.Fprintf(w, "Total fish after %d days: %d\n", simulationDays, totalFish)
	return nil
}
