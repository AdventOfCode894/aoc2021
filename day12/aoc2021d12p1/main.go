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
	name    string
	big     bool
	tunnels []*cave
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
		caves[start].tunnels = append(caves[start].tunnels, caves[end])
		caves[end].tunnels = append(caves[end].tunnels, caves[start])
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}
	if _, haveStart := caves["start"]; !haveStart {
		return 0, errors.New("no starting location")
	}

	visited := make(map[*cave]struct{})
	paths := allPaths(caves["start"], visited)

	return paths, nil
}

func isCaveNameBig(s string) bool {
	for _, r := range []rune(s) {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func allPaths(c *cave, visited map[*cave]struct{}) int {
	if c.name == "end" {
		return 1
	}
	if !c.big {
		if _, v := visited[c]; v {
			return 0
		}
		visited[c] = struct{}{}
		defer func() { delete(visited, c) }()
	}
	var paths int
	for _, next := range c.tunnels {
		paths += allPaths(next, visited)
	}
	return paths
}
