package main

import (
	"fmt"
	"io"
	"sort"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	var scores []int
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		stack, ok := parseNavigation(pr.LineRunes())
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
					return 0, fmt.Errorf("unknown completion stack rune: %c", stack[i])
				}
			}
			scores = append(scores, score)
		}
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}
	sort.Ints(scores)
	middleScore := scores[len(scores)/2]
	return middleScore, nil
}

func parseNavigation(s []rune) ([]rune, bool) {
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
			return nil, false
		}
		if expected > 0 {
			if len(stack) < 1 || stack[len(stack)-1] != expected {
				return nil, false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return stack, true
}
