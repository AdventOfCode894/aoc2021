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

type point2D struct {
	x uint
	y uint
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	var points []point2D
	maxX := uint(0)
	maxY := uint(0)
	for pr.NextLine() && !pr.IsLineEmpty() {
		tr := pr.LineTokenReader()
		var p point2D
		p.x, _ = tr.NextUint(',', 10)
		p.y, _ = tr.NextUint(aocio.EOLDelim, 10)
		tr.ConsumeEOL()

		points = append(points, p)
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	paper := make([][]bool, maxY+1)
	for row := range paper {
		paper[row] = make([]bool, maxX+1)
	}

	for _, p := range points {
		paper[p.y][p.x] = true
	}

	for pr.NextNonEmptyLine() {
		tr := pr.LineTokenReader()
		tr.ConsumeString("fold along ")
		dir, _ := tr.NextRune()
		var vert bool
		switch dir {
		case 'y':
			vert = true
		case 'x':
			vert = false
		default:
			return 0, fmt.Errorf("unknown fold direction: %c", dir)
		}
		tr.ConsumeString("=")
		pos, _ := tr.NextInt(aocio.EOLDelim, 10)
		if vert {
			paper = foldPaperVert(paper, pos)
		} else {
			paper = foldPaperHoriz(paper, pos)
		}
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	for _, row := range paper {
		for _, b := range row {
			if b {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	return 0, nil
}

func foldPaperVert(paper [][]bool, foldY int) [][]bool {
	for from, to := foldY+1, foldY-1; from < len(paper) && to >= 0; from, to = from+1, to-1 {
		for x, b := range paper[from] {
			if !b {
				continue
			}
			paper[to][x] = true
		}
	}
	return paper[:foldY]
}

func foldPaperHoriz(paper [][]bool, foldX int) [][]bool {
	for y, row := range paper {
		for from, to := foldX+1, foldX-1; from < len(row) && to >= 0; from, to = from+1, to-1 {
			if !row[from] {
				continue
			}
			row[to] = true
		}
		paper[y] = paper[y][:foldX]
	}
	return paper
}
