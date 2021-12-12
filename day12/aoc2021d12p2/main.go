package main

import (
	"errors"
	"io"
	"unicode"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

type cave struct {
	name  string
	big   bool
	links []*caveLink
}

type caveLink struct {
	dest *cave
}

type cavePath struct {
	caves []*cave
}

func solvePuzzle(r io.Reader) (int, error) {
	caves := make(map[string]*cave)
	pr := aocio.NewPuzzleReader(r)
	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		start, _ := tr.NextString('-')
		end, _ := tr.NextString(aocio.EOLDelim)
		tr.ConsumeEOL()

		if _, known := caves[start]; !known {
			caves[start] = &cave{name: start}
		}
		if _, known := caves[end]; !known {
			caves[end] = &cave{name: end}
		}
		caves[start].big = isCaveNameBig(start)
		caves[end].big = isCaveNameBig(end)
		caves[start].links = append(caves[start].links, &caveLink{dest: caves[end]})
		caves[end].links = append(caves[end].links, &caveLink{dest: caves[start]})
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}
	if _, haveStart := caves["start"]; !haveStart {
		return 0, errors.New("no starting location")
	}

	visitedSmall := make(map[*cave]struct{})
	visitedBig := make(map[*cave]struct{})
	paths := allPaths(caves["start"], visitedSmall, visitedBig, false, nil)

	return len(paths), nil
}

func isCaveNameBig(s string) bool {
	for _, r := range []rune(s) {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func allPaths(start *cave, visitedSmall map[*cave]struct{}, visitedBig map[*cave]struct{}, usedTwiceOption bool, partialPath []*cave) []*cavePath {
	if start.name == "end" {
		partialPath = append(partialPath, start)
		pathCopy := make([]*cave, len(partialPath))
		copy(pathCopy, partialPath)
		return []*cavePath{{caves: pathCopy}}
	}
	if start.big {
		if _, visited := visitedBig[start]; visited {
			return nil
		}
		visitedBig[start] = struct{}{}
		defer func() { delete(visitedBig, start) }()
	} else {
		_, visited := visitedSmall[start]
		if visited {
			if start.name == "start" {
				return nil
			}
			if usedTwiceOption {
				return nil
			}
			usedTwiceOption = true
		}
		visitedBig = make(map[*cave]struct{})
		visitedSmall[start] = struct{}{}
		if !visited {
			defer func() { delete(visitedSmall, start) }()
		}
	}
	var paths []*cavePath
	for _, next := range start.links {
		newPaths := allPaths(next.dest, visitedSmall, visitedBig, usedTwiceOption, append(partialPath, start))
		paths = append(paths, newPaths...)
	}
	return paths
}
