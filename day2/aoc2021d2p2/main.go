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
	horizontal := uint(0)
	depth := uint(0)
	aim := uint(0)
	for {
		var command string
		var amount uint
		if _, err := fmt.Fscanf(r, "%s %d\n", &command, &amount); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read line: %v", err)
			}
			break
		}
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
			return fmt.Errorf("unknown command: %s", command)
		}
	}
	area := horizontal * depth
	_, _ = fmt.Fprintf(w, "Final position:\n  Horizontal: %d\n  Depth: %d\n  Multiplied: %d\n", horizontal, depth, area)
	return nil
}
