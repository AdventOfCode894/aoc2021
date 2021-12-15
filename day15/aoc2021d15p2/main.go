package main

import (
	"io"
	"math"

	"github.com/AdventOfCode894/aoc2021/internal/aoclib/astar"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	tile, tileWidth, tileHeight := pr.Read2DUintArray(aocio.NoDelim, 10)
	if err := pr.Err(); err != nil {
		return 0, err
	}
	if err := pr.Err(); err != nil {
		return 0, err
	}

	width := tileWidth * 5
	height := tileHeight * 5

	indexToPoint := func(i int) (x int, y int) {
		x = i % width
		y = i / width
		return x, y
	}
	pointToIndex := func(x int, y int) int {
		return y*width + x
	}
	riskLevel := func(i int) uint {
		x, y := indexToPoint(i)
		tileX := x % tileWidth
		parallelUniverseX := x / tileWidth
		tileY := y % tileHeight
		parallelUniverseY := y / tileHeight

		tileRisk := tile[tileY][tileX]
		tileRisk += uint(parallelUniverseX)
		tileRisk += uint(parallelUniverseY)
		for tileRisk >= 10 {
			tileRisk -= 9
		}
		return tileRisk
	}

	shortestPath := astar.Search(0, func(from int) []int {
		x, y := indexToPoint(from)
		var neighbors []int
		if x > 0 {
			neighbors = append(neighbors, pointToIndex(x-1, y))
		}
		if y > 0 {
			neighbors = append(neighbors, pointToIndex(x, y-1))
		}
		if x < width-1 {
			neighbors = append(neighbors, pointToIndex(x+1, y))
		}
		if y < height-1 {
			neighbors = append(neighbors, pointToIndex(x, y+1))
		}
		return neighbors
	}, func(from, to int) float64 {
		return float64(riskLevel(to))
	}, func(from int) float64 {
		x, y := indexToPoint(from)
		destX, destY := width-1, height-1
		dx := float64(destX) - float64(x)
		dy := float64(destY) - float64(y)
		return math.Sqrt(dx*dx + dy*dy)
	}, func(at int) bool {
		return at == pointToIndex(width-1, height-1)
	})

	pathRisk := uint(0)
	for _, i := range shortestPath[1:] {
		pathRisk += riskLevel(i)
	}

	return int(pathRisk), nil
}
