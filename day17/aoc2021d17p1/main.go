package main

import (
	"io"

	"github.com/AdventOfCode894/aoc2021/internal/aocio"
	"github.com/AdventOfCode894/aoc2021/internal/aocmain"
)

func main() {
	aocmain.HandlePuzzle(solvePuzzle)
}

func solvePuzzle(r io.Reader) (int, error) {
	pr := aocio.NewPuzzleReader(r)
	pr.NextNonEmptyLine()
	tr := pr.LineTokenReader()
	tr.ConsumeString("target area: x=")
	x1, _ := tr.NextInt('.', 10)
	tr.ConsumeString(".")
	x2, _ := tr.NextInt(',', 10)
	tr.ConsumeString(" y=")
	y1, _ := tr.NextInt('.', 10)
	tr.ConsumeString(".")
	y2, _ := tr.NextInt(aocio.EOLDelim, 10)
	tr.ConsumeEOL()
	if err := pr.Err(); err != nil {
		return 0, err
	}

	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	minX := 0
	maxX := x2
	minY := y1
	maxY := x2

	maxHighestY := 0
	for x := minX; x <= maxX; x++ {
		for y := minY; y < maxY; y++ {
			highestY, ok := simulate(x, y, x1, x2, y1, y2)
			if !ok {
				continue
			}
			if highestY > maxHighestY {
				maxHighestY = highestY
			}
		}
	}

	return maxHighestY, nil
}

func simulate(initialX int, initialY int, x1 int, x2 int, y1 int, y2 int) (int, bool) {
	velX := initialX
	velY := initialY
	posX := 0
	posY := 0
	highestY := 0
	steps := 0
	hitTarget := false
	for {
		posX += velX
		if velX > 0 {
			velX--
		}
		posY += velY
		velY--
		if posX > x2 {
			break
		}
		if posY < y1 {
			break
		}
		if posY > highestY {
			highestY = posY
		}
		steps++
		if posX >= x1 && posY <= y2 {
			hitTarget = true
		}
	}
	if !hitTarget {
		return 0, false
	}
	return highestY, true
}
