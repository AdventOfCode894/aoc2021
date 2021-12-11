package aoc2021d4

import (
	"github.com/AdventOfCode894/aoc2021/internal/aocio"
)

func ReadDraws(pr *aocio.PuzzleReader) ([]uint, error) {
	draws := pr.ReadUintArrayLine(',', 10)
	return draws, pr.Err()
}

func ReadBingoCard(pr *aocio.PuzzleReader) (*BingoCard, error) {
	if !pr.NextNonEmptyLine() {
		return nil, pr.Err()
	}

	cells, width, height := pr.Read2DUintArray(' ', 10)
	if err := pr.Err(); err != nil {
		return nil, err
	}
	unwrapped := make([]uint, 0, width*height)
	for _, row := range cells {
		unwrapped = append(unwrapped, row...)
	}

	return newBingoCard(unwrapped, width)
}
