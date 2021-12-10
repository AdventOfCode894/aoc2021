package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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
	var scores []int
	for scanner.Scan() {
		line := scanner.Text()
		_, stack, ok := parseNavigation(line)
		if ok {
			score := 0
			for i := len(stack) - 1; i >= 0; i-- {
				score *= 5
				switch stack[i] {
				case '(':
					score += 1
				case '[':
					score += 2
				case '{':
					score += 3
				case '<':
					score += 4
				default:
					return fmt.Errorf("unknown completion stack rune: %c", stack[i])
				}
			}
			scores = append(scores, score)
		}
	}
	sort.Ints(scores)
	middleScore := scores[len(scores)/2]
	_, _ = fmt.Fprintln(w, middleScore)
	return nil
}

func parseNavigation(s string) (rune, []rune, bool) {
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
			return c, nil, false
		}
		if expected > 0 {
			if len(stack) < 1 || stack[len(stack)-1] != expected {
				return c, nil, false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return 0, stack, true
}
