package aoc2021d4

import "fmt"

type BingoCard struct {
	cells  [][]uint // Row-major
	marked [][]bool
	won    bool

	width  int
	height int

	// Caches for performance enhancement
	rowMarks      []int
	colMarks      []int
	cellLocations map[uint]cardLocation
}

type cardLocation struct {
	row int
	col int
}

func newBingoCard(cells []uint, width int) (*BingoCard, error) {
	if len(cells)%width != 0 {
		return nil, fmt.Errorf("bingo card has impossible geometry: width = %d but cell count = %d", width, len(cells))
	}
	height := len(cells) / width

	bc := &BingoCard{
		cells:  make([][]uint, height),
		marked: make([][]bool, height),

		width:  width,
		height: height,

		rowMarks:      make([]int, height),
		colMarks:      make([]int, width),
		cellLocations: make(map[uint]cardLocation),
	}
	for row := 0; row < height; row++ {
		bc.cells[row] = make([]uint, width)
		bc.marked[row] = make([]bool, width)
		for col := 0; col < width; col++ {
			cell := cells[0]
			cells = cells[1:]
			bc.cells[row][col] = cell
			bc.cellLocations[cell] = cardLocation{row: row, col: col}
		}
	}

	return bc, nil
}

func (bc *BingoCard) Mark(draw uint) (winner bool) {
	loc, found := bc.cellLocations[draw]
	if found && !bc.marked[loc.row][loc.col] {
		bc.marked[loc.row][loc.col] = true
		bc.rowMarks[loc.row]++
		bc.colMarks[loc.col]++
		if bc.rowMarks[loc.row] >= bc.width || bc.colMarks[loc.col] >= bc.height {
			bc.won = true
		}
	}
	return bc.won
}

func (bc *BingoCard) Won() bool { return bc.won }

func (bc *BingoCard) UnmarkedSum() uint {
	sum := uint(0)
	for row := 0; row < bc.height; row++ {
		for col := 0; col < bc.width; col++ {
			if !bc.marked[row][col] {
				sum += bc.cells[row][col]
			}
		}
	}
	return sum
}
