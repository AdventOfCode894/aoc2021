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
	// This is really bad dynamic programming. Better solution: day 6 with fish named after digraphs!

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

	subQuantities := make(map[uint]map[rune]map[rune]map[rune]uint)

	const maxSteps = 40
	quantities := make(map[rune]uint)
	for i := 0; i < len(polymerTemplate)-1; i++ {
		first := polymerTemplate[i]
		second := polymerTemplate[i+1]
		recordQuantitiesFor(first, second, maxSteps, rules, subQuantities)
		for r, c := range subQuantities[maxSteps][first][second] {
			quantities[r] += c
		}
		quantities[second]--
	}
	quantities[polymerTemplate[len(polymerTemplate)-1]]++

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

	return int(answer), nil
}

func recordQuantitiesFor(first rune, second rune, steps uint, rules map[rune]map[rune]rune, quantities map[uint]map[rune]map[rune]map[rune]uint) {
	if _, ok := quantities[steps]; !ok {
		quantities[steps] = make(map[rune]map[rune]map[rune]uint)
	}
	if _, ok := quantities[steps][first]; !ok {
		quantities[steps][first] = make(map[rune]map[rune]uint)
	}
	if _, alreadyDone := quantities[steps][first][second]; alreadyDone {
		return
	}

	quantities[steps][first][second] = make(map[rune]uint)
	if steps < 1 {
		quantities[steps][first][second][first]++
		quantities[steps][first][second][second]++
		return
	}
	if _, ok := rules[first]; !ok {
		quantities[steps][first][second][first]++
		quantities[steps][first][second][second]++
		return
	}
	insert, ok := rules[first][second]
	if !ok {
		quantities[steps][first][second][first]++
		quantities[steps][first][second][second]++
		return
	}

	recordQuantitiesFor(first, insert, steps-1, rules, quantities)
	recordQuantitiesFor(insert, second, steps-1, rules, quantities)
	quantities[steps][first][second] = make(map[rune]uint)
	for r, c := range quantities[steps-1][first][insert] {
		quantities[steps][first][second][r] += c
	}
	for r, c := range quantities[steps-1][insert][second] {
		quantities[steps][first][second][r] += c
	}
	quantities[steps][first][second][insert]--
}
