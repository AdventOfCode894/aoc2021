package main

import (
	"fmt"
	"io"
	"math"
	"unicode/utf8"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	rules := make(map[string]rune)
	pr.NextLine()
	polymerTemplate := pr.LineString()
	if err := pr.Err(); err != nil {
		return 0, err
	}
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		bigram, _ := tr.NextString(' ')
		if len(bigram) != 2 {
			return 0, fmt.Errorf("rule did not start with a bigram: %s", bigram)
		}
		tr.ConsumeString("-> ")
		insert, _ := tr.NextRune()
		tr.ConsumeEOL()
		rules[bigram] = insert
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	bigramInstances := make(map[string]uint)
	for i := 0; i < len(polymerTemplate)-1; i++ {
		bigramInstances[polymerTemplate[i:i+2]]++
	}

	for step := 0; step < 40; step++ {
		newInstances := make(map[string]uint)
		for bigram, count := range bigramInstances {
			insertion, hasRule := rules[bigram]
			if !hasRule {
				newInstances[bigram] = count
				continue
			}
			bigramRunes := []rune(bigram)
			first, second := bigramRunes[0], bigramRunes[1]
			newInstances[string(first)+string(insertion)] += count
			newInstances[string(insertion)+string(second)] += count
		}
		bigramInstances = newInstances
	}

	letterInstances := make(map[rune]uint)
	for bigram, count := range bigramInstances {
		l := []rune(bigram)[0]
		letterInstances[l] += count
	}
	lastLetter, _ := utf8.DecodeLastRuneInString(polymerTemplate)
	letterInstances[lastLetter]++

	minLetterCount := uint(math.MaxUint)
	maxLetterCount := uint(0)
	for _, count := range letterInstances {
		if count < minLetterCount {
			minLetterCount = count
		}
		if count > maxLetterCount {
			maxLetterCount = count
		}
	}

	answer := maxLetterCount - minLetterCount
	return int(answer), nil
}
