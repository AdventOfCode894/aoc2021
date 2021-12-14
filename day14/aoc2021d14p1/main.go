package main

import (
	"io"
	"sort"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	rules := make(map[rune]map[rune]rune)
	pr.NextLine()
	polymerTemplate := []rune(pr.LineString())
	if err := pr.Err(); err != nil {
		return 0, err
	}
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		first, _ := tr.NextRune()
		second, _ := tr.NextRune()
		tr.ConsumeString(" -> ")
		insert, _ := tr.NextRune()
		tr.ConsumeEOL()
		if _, seenFirst := rules[first]; !seenFirst {
			rules[first] = make(map[rune]rune)
		}
		rules[first][second] = insert
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	for step := 0; step < 10; step++ {
		i := 0
		for i < len(polymerTemplate)-1 {
			first := polymerTemplate[i]
			second := polymerTemplate[i+1]
			i++
			firstRule, ok := rules[first]
			if !ok {
				continue
			}
			insert, ok := firstRule[second]
			if !ok {
				continue
			}
			polymerTemplate = append(polymerTemplate[:i], append([]rune{insert}, polymerTemplate[i:]...)...)
			i++
		}
	}

	quantities := make(map[rune]int)
	for _, r := range polymerTemplate {
		quantities[r]++
	}
	var polymers []rune
	for r := range quantities {
		polymers = append(polymers, r)
	}
	sort.Slice(polymers, func(i, j int) bool {
		return quantities[polymers[i]] < quantities[polymers[j]]
	})

	leastCommon := quantities[polymers[0]]
	mostCommon := quantities[polymers[len(polymers)-1]]
	answer := mostCommon - leastCommon

	return answer, nil
}
