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
	var zeroes, ones []uint
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		for i, c := range pr.LineRunes() {
			switch c {
			case '0', '1':
				if i >= len(zeroes) {
					zeroes = append(zeroes, 0)
					ones = append(ones, 0)
				}
				if c == '0' {
					zeroes[i]++
				} else {
					ones[i]++
				}
			default:
				return 0, fmt.Errorf("invalid character in input: U+%d", c)
			}
		}
	}
	if pr.Err() != nil {
		return 0, pr.Err()
	}

	gamma := uint(0)
	epsilon := uint(0)
	for i := range zeroes {
		gamma <<= 1
		epsilon <<= 1
		if zeroes[i] > ones[i] {
			epsilon |= 1
		} else {
			gamma |= 1
		}
	}
	power := gamma * epsilon
	return int(power), nil
}
