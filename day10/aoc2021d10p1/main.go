package main

import (
	"fmt"
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	score := 0
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		c, ok := parseNavigation(pr.LineRunes())
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
				return 0, fmt.Errorf("unknown first illegal rune: %c", c)
			}
		}
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}
	return score, nil
}

func parseNavigation(s []rune) (rune, bool) {
	var stack []rune
	for _, c := range s {
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
