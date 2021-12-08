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
	// Sorry for this mess! Yikes.
	uniqueValues := 0
	for {
		var patterns []string
		for i := 0; i < 10; i++ {
			var pattern string
			if _, err := fmt.Fscanf(r, "%s", &pattern); err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("failed to read pattern: %v", err)
				}
				break
			}
			patterns = append(patterns, pattern)
		}
		var delim rune
		if _, err := fmt.Fscanf(r, "%c", &delim); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read delimeter: %v", err)
			}
			break
		}
		if delim != '|' {
			return fmt.Errorf("wrong delimeter: %c", delim)
		}
		var outputs []string
		for i := 0; i < 4; i++ {
			var output string
			if _, err := fmt.Fscanf(r, "%s", &output); err != nil {
				if !errors.Is(err, io.EOF) {
					return fmt.Errorf("failed to read output: %v", err)
				}
				break
			}
			outputs = append(outputs, output)
		}

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

	_, _ = fmt.Fprintf(w, "Unique values in output %d\n", uniqueValues)
	return nil
}
