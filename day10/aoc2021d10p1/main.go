package main

import (
	"bufio"
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
	scanner := bufio.NewScanner(r)
	score := 0
	for scanner.Scan() {
		line := scanner.Text()
		c, ok := parseNavigation(line)
		if !ok {
			switch c {
			case ')':
				score += 3
			case ']':
				score += 57
			case '}':
				score += 1197
			case '>':
				score += 25137
			default:
				return fmt.Errorf("unknown first illegal rune: %c", c)
			}
		}
	}
	_, _ = fmt.Fprintln(w, score)
	return nil
}

func parseNavigation(s string) (rune, bool) {
	var stack []rune
	for _, c := range []rune(s) {
		var expected rune
		switch c {
		case '(', '[', '{', '<':
			stack = append(stack, c)
		case ')':
			expected = '('
		case ']':
			expected = '['
		case '}':
			expected = '{'
		case '>':
			expected = '<'
		default:
			return c, false
		}
		if expected > 0 {
			if len(stack) < 1 || stack[len(stack)-1] != expected {
				return c, false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return 0, true
}
