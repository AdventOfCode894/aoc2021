package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := solvePuzzle(os.Stdin, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func solvePuzzle(r io.Reader, w io.Writer) error {
	ring := newMeasurementRing(3)
	lastSum := ^uint(0)
	increases := 0
	for {
		var depth uint
		if _, err := fmt.Fscanf(r, "%d\n", &depth); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read line: %v", err)
			}
			break
		}
		ring.Record(depth)
		if ring.IsFull() {
			sum := ring.Sum()
			if sum > lastSum {
				increases++
			}
			lastSum = sum
		}
	}
	_, _ = fmt.Fprintf(w, "Sliding depth increased %d times\n", increases)
	return nil
}
