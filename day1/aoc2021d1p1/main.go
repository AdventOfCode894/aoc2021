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
	lastDepth := ^uint(0)
	increases := 0
	for {
		var depth uint
		if _, err := fmt.Fscanf(r, "%d\n", &depth); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read line: %v", err)
			}
			break
		}
		if depth > lastDepth {
			increases++
		}
		lastDepth = depth
	}
	_, _ = fmt.Fprintf(w, "Depth increased %d times\n", increases)
	return nil
}
